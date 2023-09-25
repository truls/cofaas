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
	"flag"
	"fmt"
	"io"
	"net"
	"os"

	ctrdlog "github.com/containerd/containerd/log"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "cofaas_orig/protos/prodcon"
)

const (
	INLINE      = "INLINE"
	S3          = "S3"
	ELASTICACHE = "ELASTICACHE"
)

var verbose = flag.Bool("v", false, "Be verbose")

type consumerServer struct {
	transferType   string
	pb.UnimplementedProducerConsumerServer
}

func (s *consumerServer) ConsumeByte(ctx context.Context, str *pb.ConsumeByteRequest) (*pb.ConsumeByteReply, error) {
	if *verbose {
		log.Printf("[consumer] Consumed %d bytes\n", len(str.Value))
	}
	return &pb.ConsumeByteReply{Value: true}, nil
}

func (s *consumerServer) ConsumeStream(stream pb.ProducerConsumer_ConsumeStreamServer) error {
	for {
		str, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.ConsumeByteReply{Value: true})
		}
		if err != nil {
			return err
		}
		log.Printf("[consumer] Consumed string of length %d\n", len(str.Value))
	}
}

func main() {
	//flagAddress := flag.String("addr", "consumer.default.192.168.1.240.sslip.io", "Server IP address")
	port := flag.Int("ps", 80, "Port")
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: ctrdlog.RFC3339NanoFixed,
		FullTimestamp:   true,
	})
	log.SetOutput(os.Stdout)

		log.Println("consumer has tracing DISABLED")

	transferType, ok := os.LookupEnv("TRANSFER_TYPE")
	if !ok {
		log.Infof("TRANSFER_TYPE not found, using INLINE transfer")
		transferType = "INLINE"
	}

	//set up server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("[consumer] failed to listen: %v", err)
	}

	var grpcServer *grpc.Server
	grpcServer = grpc.NewServer()
	cs := consumerServer{transferType: transferType}
	pb.RegisterProducerConsumerServer(grpcServer, &cs)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[consumer] failed to serve: %s", err)
	}
}
