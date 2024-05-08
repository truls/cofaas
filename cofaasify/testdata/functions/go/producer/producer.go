// MIT License
//
// Copyright (c) 2021 Michal Baczun and EASE lab
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"context"
	"time"
	//"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc/credentials/insecure"

	ctrdlog "github.com/containerd/containerd/log"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/reflection"

	pb_client "cofaas_orig/protos/prodcon"

	pb "cofaas_orig/protos/helloworld"

	"google.golang.org/grpc"
)

type producerServer struct {
	consumerAddr string
	consumerPort int
	payloadData  []byte
	transferType string
	randomStr    string
	pb.UnimplementedGreeterServer
}

const (
	INLINE      = "INLINE"
	XDT         = "XDT"
	S3          = "S3"
	ELASTICACHE = "ELASTICACHE"
)

// var verbose = flag.Bool("v", false, "Be verbose")
// var repeats = flag.Int("r", 1, "Repeat message");
var v = false
var verbose = &v

var repetitions = 1

var measure_time = false
var timings = []float64{}

//var

func getGRPCclient(addr string) (pb_client.ProducerConsumerClient, *grpc.ClientConn) {
	// establish a connection
	var conn *grpc.ClientConn
	var err error
	conn, err = grpc.Dial(addr, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("[producer] fail to dial: %s", err)
	}
	return pb_client.NewProducerConsumerClient(conn), conn
}

func (ps *producerServer) SayHello(ctx context.Context, req *pb.HelloRequest) (_ *pb.HelloReply, err error) {
	addr := fmt.Sprintf("%v:%v", ps.consumerAddr, ps.consumerPort)
	client, conn := getGRPCclient(addr)
	defer conn.Close()
	payloadToSend := ps.payloadData
	if measure_time {
		start := time.Now()
		for i := 1; i <= repetitions; i++ {
			ack, err := client.ConsumeByte(ctx, &pb_client.ConsumeByteRequest{Value: payloadToSend})
			if err != nil {
				log.Fatalf("[producer] client error in string consumption: %s", err)
			}
			if *verbose {
				log.Printf("[producer] (single) Ack: %v\n", ack.Value)
			}
		}
		duration := time.Since(start)
		latency := duration.Microseconds() / int64(repetitions)

		if *verbose {
			log.Printf("[producer] Returing latency %d", latency)
		}

		return &pb.HelloReply{Message: fmt.Sprintf("%d", latency)}, err
	} else {
		for i := 1; i <= repetitions; i++ {
			ack, err := client.ConsumeByte(ctx, &pb_client.ConsumeByteRequest{Value: payloadToSend})
			if err != nil {
				log.Fatalf("[producer] client error in string consumption: %s", err)
			}
			if *verbose {
				log.Printf("[producer] (single) Ack: %v\n", ack.Value)
			}
		}
		return &pb.HelloReply{Message: "Success"}, err
	}
}

func main() {
	// flagAddress := flag.String("addr", "consumer.default.192.168.1.240.sslip.io", "Server IP address")
	// flagClientPort := flag.Int("pc", 80, "Client Port")
	// flagServerPort := flag.Int("ps", 80, "Server Port")
	// flag.Parse()
	flagAddressV := "consumer"
	flagAddress := &flagAddressV
	flagClientPortV := 3030
	flagServerPortV := 3031
	flagClientPort := &flagClientPortV
	flagServerPort := &flagServerPortV

	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: ctrdlog.RFC3339NanoFixed,
		FullTimestamp:   true,
	})
	log.SetOutput(os.Stdout)

	log.Println("producer has tracing DISABLED")

	var grpcServer *grpc.Server
	grpcServer = grpc.NewServer()

	//client setup
	log.Printf("[producer] Client using address: %v:%d\n", *flagAddress, *flagClientPort)

	ps := producerServer{consumerAddr: *flagAddress, consumerPort: *flagClientPort}
	transferType, ok := os.LookupEnv("TRANSFER_TYPE")
	if !ok {
		log.Infof("TRANSFER_TYPE not found, using INLINE transfer")
		transferType = INLINE
	}
	log.Infof("[producer] transfering via %s", transferType)
	ps.transferType = transferType

	transferSizeKB := 1 //4095
	if value, ok := os.LookupEnv("TRANSFER_SIZE_KB"); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			transferSizeKB = intValue
		} else {
			log.Infof("invalid TRANSFER_SIZE_KB: %s, using default %d", value, transferSizeKB)
		}
	}

	if value, ok := os.LookupEnv("REPEATS"); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			repetitions = intValue
		} else {
			log.Infof("invalid REPEATS: %s, using default %d", value, repetitions)
		}
	}

	if value, ok := os.LookupEnv("MEASURE_LAT"); ok {
		if value == "true" {
			log.Info("RUnning in latency measurement mode")
			measure_time = true
			timings = make([]float64, repetitions)
		}
	} else {
		log.Infof("invalid MEASURE_LAT: %s, using default %d", value, false)
	}

	if value, ok := os.LookupEnv("VERBOSE"); ok {
		v = value == "true"
	} else {
		log.Infof("invalid VERBOSE: %s, using default %b", value, false)
	}

	// 4194304 bytes is the limit by gRPC
	payloadData := make([]byte, transferSizeKB*1024)
	if _, err := rand.Read(payloadData); err != nil {
		log.Fatal(err)
	}
	ps.randomStr = os.Getenv("HOSTNAME")

	log.Infof("sending %d bytes to consumer", len(payloadData))
	log.Infof("repeating message %d times", repetitions)
	ps.payloadData = payloadData
	pb.RegisterGreeterServer(grpcServer, &ps)
	reflection.Register(grpcServer)

	//server setup
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *flagServerPort))
	if err != nil {
		log.Fatalf("[producer] failed to listen: %v", err)
	}

	log.Println("[producer] Server Started")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[producer] failed to serve: %s", err)
	}

}
