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

package impl

import (
	"context"
	"flag"
	"fmt"
	//storage "github.com/vhive-serverless/vSwarm/utils/storage/go"
	"math/rand"
	//"net"
	net "github.com/truls/chained-service-example/stubs/net"
	"os"
	"strconv"

	//sdk "github.com/ease-lab/vhive-xdt/sdk/golang"
	///"github.com/ease-lab/vhive-xdt/utils"
	//"google.golang.org/grpc/credentials/insecure"

	//log "github.com/sirupsen/logrus"
	//"google.golang.org/grpc/reflection"

	pb_client "github.com/truls/chained-service-example/prodcon_stub"

	pb "github.com/truls/chained-service-example/helloworld_stub"

	//tracing "github.com/vhive-serverless/vSwarm/utils/tracing/go"
	grpc "github.com/truls/chained-service-example/stubs/grpc"
)

type producerServer struct {
	consumerAddr   string
	consumerPort   int
	payloadData    []byte
	transferType   string
	randomStr      string
	pb.UnimplementedGreeterServer
}

const (
	INLINE      = "INLINE"
	S3          = "S3"
	ELASTICACHE = "ELASTICACHE"
)


func getGRPCclient(addr string) (pb_client.ProducerConsumerClient, *grpc.ClientConn) {
	// establish a connection
	var conn *grpc.ClientConn
	var err error
	// if tracing.IsTracingEnabled() {
	// 	conn, err = tracing.DialGRPCWithUnaryInterceptor(addr, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	// } else {
	//conn, err = grpc.Dial(addr, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		conn, err = grpc.Dial(addr, grpc.WithBlock(), grpc.WithTransportCredentials(nil))
	//}
	if err != nil {
		fmt.Printf("[producer] fail to dial: %s", err)
		os.Exit(1)
	}
	return pb_client.NewProducerConsumerClient(conn), conn
}

func (ps *producerServer) SayHello(ctx context.Context, req *pb.HelloRequest) (_ *pb.HelloReply, err error) {
	addr := fmt.Sprintf("%v:%v", ps.consumerAddr, ps.consumerPort)
	client, conn := getGRPCclient(addr)
	defer conn.Close()
	payloadToSend := ps.payloadData
	ack, err := client.ConsumeByte(ctx, &pb_client.ConsumeByteRequest{Value: payloadToSend})
	if err != nil {
		fmt.Printf("[producer] client error in string consumption: %s", err)
		os.Exit(1)
	}
	fmt.Printf("[producer] (single) Ack: %v\n", ack.Value)
	return &pb.HelloReply{Message: "Success"}, err
}

func Main() {
	flagAddress := flag.String("addr", "consumer.default.192.168.1.240.sslip.io", "Server IP address")
	flagClientPort := flag.Int("pc", 80, "Client Port")
	flagServerPort := flag.Int("ps", 80, "Server Port")
	//url := flag.String("zipkin", "http://zipkin.istio-system.svc.cluster.local:9411/api/v2/spans", "zipkin url")
	//dockerCompose := flag.Bool("dockerCompose", false, "Env docker Compose?")
	flag.Parse()


	// if tracing.IsTracingEnabled() {
	// 	log.Println("producer has tracing enabled")
	// 	shutdown, err := tracing.InitBasicTracer(*url, "producer")
	// 	if err != nil {
	// 		log.Warn(err)
	// 	}
	// 	defer shutdown()
	// } else {
		fmt.Println("producer has tracing DISABLED")
	// }

	var grpcServer *grpc.Server
	// if tracing.IsTracingEnabled() {
	// 	grpcServer = tracing.GetGRPCServerWithUnaryInterceptor()
	// } else {
	grpcServer = grpc.NewServer()
	// }

	//client setup
	fmt.Printf("[producer] Client using address: %v:%d\n", *flagAddress, *flagClientPort)

	ps := producerServer{consumerAddr: *flagAddress, consumerPort: *flagClientPort}

	transferType, ok := os.LookupEnv("TRANSFER_TYPE")
	if !ok {
		fmt.Printf("TRANSFER_TYPE not found, using INLINE transfer")
		transferType = INLINE
	}
	fmt.Printf("[producer] transfering via %s", transferType)
	ps.transferType = transferType

	transferSizeKB := 4095
	if value, ok := os.LookupEnv("TRANSFER_SIZE_KB"); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			transferSizeKB = intValue
		} else {
			fmt.Printf("invalid TRANSFER_SIZE_KB: %s, using default %d", value, transferSizeKB)
		}
	}

	// 4194304 bytes is the limit by gRPC
	payloadData := make([]byte, transferSizeKB*1024)
	if _, err := rand.Read(payloadData); err != nil {
		fmt.Print(err, "\n")
		os.Exit(1)
	}
	ps.randomStr = os.Getenv("HOSTNAME")

	fmt.Printf("sending %d bytes to consumer", len(payloadData))
	ps.payloadData = payloadData
	pb.RegisterGreeterServer(grpcServer, &ps)
	//reflection.Register(grpcServer)

	//server setup
	// TODO: Handle this
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *flagServerPort))
	if err != nil {
		fmt.Printf("[producer] failed to listen: %v", err)
		os.Exit(1)
	}

	fmt.Println("[producer] Server Started")

	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("[producer] failed to serve: %s", err)
		os.Exit(1)
	}

}
