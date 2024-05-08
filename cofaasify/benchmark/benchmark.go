package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path"
	"time"

	"cofaas/benchmark/invoker"
	"cofaas/benchmark/invoker/endpoint"

	ctrdlog "github.com/containerd/containerd/log"
	"github.com/go-errors/errors"
	log "github.com/sirupsen/logrus"
)

type benchCfg struct {
	name           string
	arg            string
	startFun       startFunc
	envVars        []string
	sizes          []int
	recordInnerLat bool
}

var sizes_go_wasm = []int{1, 2, 4, 8, 16} //, 128}
var sizes_go_nogc = []int{1, 2, 4, 8, 16, 32, 64}
var sizes_rust = []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512}

// var sizes = []int{8, 16} //, 32, 64, 128}
// var repeats = []int{10,20}
var repeats = []int{1, 10, 20}
var endp = []*endpoint.Endpoint{{
	Hostname: "localhost:3031",
}}

type benchmark interface {
	getName() string
	stop() error
}

// type startFunc
type stopFunc func(*exec.Cmd) error
type startFunc func(string, string, []string) (benchmark, error)

type runner struct {
	benchmark
	ctx      context.Context
	sout     io.ReadCloser
	serr     io.ReadCloser
	cmd      *exec.Cmd
	stopFun  stopFunc
	canceler context.CancelFunc
	name     string
}

// TODO: add context with a timeout to prevent benchmarking process
// from freezing

func (r *runner) stop() error {
	r.canceler()
	return r.stopFun(r.cmd)
}

func (r *runner) getName() string {
	return r.name
}

func getMyPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}

	return path.Dir(exe), nil
}

// Adapted from https://stackoverflow.com/a/56336811
func try_connect(host string, port string) (bool, error) {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		// fmt.Printf("%v", err)
		return false, nil
	}
	if conn != nil {
		defer conn.Close()
		return true, nil
	}
	return false, errors.Errorf("inconclusive state in connection test")
}

func wait_for_port(host string, port string) error {
	tick := time.Tick(time.Second)
	for i := 1; i <= 10; i++ {
		<-tick
		res, err := try_connect(host, port)
		if err != nil {
			return errors.Wrap(err, 0)
		}
		if res {
			return nil
		}
	}
	return errors.Errorf("timeout whiile waiting for port")
}

func startWasm(arg string, name string, environ []string) (benchmark, error) {
	myDir, err := getMyPath()
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}
	wrapperPath := path.Join(myDir, "../../wrapper/target/release/helloworld-server")
	componentPath := path.Join(myDir, arg)
	stopFun := func(cmd *exec.Cmd) error {
		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			return err
		}

		if err := cmd.Wait(); err != nil {
			return err
		}
		return nil
	}

	return startCommand(environ, name, stopFun, wrapperPath, componentPath)
}

func startDocker(arg string, name string, environ []string) (benchmark, error) {
	myDir, err := getMyPath()
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	composePath := path.Join(myDir, arg)

	stopFun := func(cmd *exec.Cmd) error {
		waitErr := cmd.Wait()
		c := exec.Command("podman-compose", "-f", composePath, "down")
		if err := c.Run(); err != nil {
			return errors.Errorf("failed to stop podman compose: %v", err)
		}

		if waitErr != nil {
			return errors.Wrap(err, 0)
		}
		return nil
	}

	return startCommand(environ, name, stopFun, "podman-compose", "-f", composePath, "up")
}

func startCommand(environ []string, name string, stopFun stopFunc, exe string, args ...string) (benchmark, error) {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, exe, args...)
	cmd.Env = append(os.Environ(), environ...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	// sout, err := cmd.StdoutPipe()
	// if err != nil {
	// 	cancel()
	// 	return nil, errors.Wrap(err, 0)
	// }
	// serr, err := cmd.StderrPipe()
	// if err != nil {
	// 	cancel()
	// 	return nil, errors.Wrap(err, 0)
	// }
	log.Infof("Starting benchmark server %s", cmd.String())
	if err := cmd.Start(); err != nil {
		cancel()
		fmt.Printf("Staring failed")
		return nil, errors.Wrap(err, 0)
	}

	if err := wait_for_port("::1", "3031"); err != nil {
		fmt.Printf("Staring failed2")

		// var buf []byte
		// if _, err := sout.Read(buf); err != nil {
		// 	return nil, errors.Wrap(err, 0)
		// }
		// if _, err := serr.Read(buf); err != nil {
		// 	return nil, errors.Wrap(err, 0)
		// }
		// sout.Close()
		// serr.Close()
		cancel()
		cmd.Wait()

		//res, _ := cmd.CombinedOutput()
		// return nil, errors.Errorf("%s", string(buf))
		return nil, errors.Errorf("%s", "")
	}

	return &runner{
		cmd:      cmd,
		name:     name,
		ctx:      ctx,
		canceler: cancel,
		stopFun:  stopFun,
	}, nil
}

func doBenchmark(config benchCfg, size int, repeats int, invocations int) {
	outputFile := fmt.Sprintf("result_%s_%d_%d.csv", config.name, size, repeats)

	if res, err := invoker.ResultPathExists(outputFile); err != nil {
		log.Fatal("Result path exist check failed")
	} else {
		if res {
			log.Infof("Skipping benchmark since output file %s already exists", invoker.GetResultPath(outputFile))
			return
		}
	}

	envVars := append(config.envVars,
		fmt.Sprintf("TRANSFER_SIZE_KB=%d", size),
		fmt.Sprintf("REPEATS=%d", repeats))
	b, err := config.startFun(config.arg, config.name, envVars)
	if err != nil {
		log.Fatalf("Failed to start benchmark %v", err)
	}
	defer b.stop()

	log.Infof("Invoking benchmark of %s for size %d repeating %d times", b.getName(), size, repeats)

	// if _, err := os.Stat(outputFile); os.IsNotExist(err error)

	// invoker.RunExperiment(endp, 60, 200, b.getName(), size, repeats, outputFile)
	invoker.RunExperimentSync(endp, invocations, b.getName(), size, repeats, outputFile, config.recordInnerLat)
}

func main() {
	debug := flag.Bool("dbg", false, "Enable debug logging")
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: ctrdlog.RFC3339NanoFixed,
		FullTimestamp:   true,
	})
	log.SetOutput(os.Stdout)
	if *debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug logging is enabled")
	} else {
		log.SetLevel(log.InfoLevel)
	}

	invoker.InitExperiment(endp)

	benchmarks := []benchCfg{
		{
			name:     "native-go-nogc",
			arg:      "../testdata/functions/go/docker-compose.yml",
			startFun: startDocker,
			envVars:  []string{"GOGC=off"},
			sizes:    sizes_go_nogc,
		},
		{
			name:     "native-go-gccmp",
			arg:      "../testdata/functions/go/docker-compose.yml",
			startFun: startDocker,
			sizes:    sizes_go_nogc,
		},
		{
			name:     "native-go",
			arg:      "../testdata/functions/go/docker-compose.yml",
			startFun: startDocker,
			sizes:    sizes_rust,
		},
		{
			name:     "native-rust",
			arg:      "../testdata/functions/rust/docker-compose.yml",
			startFun: startDocker,
			sizes:    sizes_rust,
		},
		{
			name:     "cofaas-go-nogc",
			arg:      "../testdata/component/go/composed.wasm",
			startFun: startWasm,
			sizes:    sizes_go_wasm,
		},
		{
			name:     "cofaas-rust",
			arg:      "../testdata/component/rust/composed.wasm",
			startFun: startWasm,
			sizes:    sizes_rust,
		},
	}

	lat_benchmarks := []benchCfg{
		{
			name:           "latency-native-go-nogc",
			arg:            "../testdata/functions/go/docker-compose.yml",
			startFun:       startDocker,
			envVars:        []string{"GOGC=off", "MEASURE_LAT=true"},
			sizes:          sizes_go_wasm,
			recordInnerLat: true,
		},
		{
			name:           "latency-native-rust",
			arg:            "../testdata/functions/rust/docker-compose.yml",
			startFun:       startDocker,
			sizes:          sizes_rust,
			envVars:        []string{"MEASURE_LAT=true", "VERBOSE=1"},
			recordInnerLat: true,
		},
		{
			name:           "latency-cofaas-go-nogc",
			arg:            "../testdata/component/go/composed.wasm",
			startFun:       startWasm,
			sizes:          sizes_go_wasm,
			envVars:        []string{"MEASURE_LAT=true"},
			recordInnerLat: true,
		},
		{
			name:           "latency-cofaas-rust",
			arg:            "../testdata/component/rust/composed.wasm",
			startFun:       startWasm,
			sizes:          sizes_rust,
			envVars:        []string{"MEASURE_LAT=true"},
			recordInnerLat: true,
		},
	}

	for _, b := range benchmarks {
		for _, r := range repeats {
			for _, s := range b.sizes {
				doBenchmark(b, s, r, 6000)
			}
		}
	}

	for _, b := range lat_benchmarks {
		for _, s := range b.sizes {
			doBenchmark(b, s, 100, 100)
		}
	}
}
