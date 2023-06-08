package main

import (
	context "context"
	"fmt"
	pb "github.com/truls/chained-service-example/helloworld_stub"
	pb_prodcon "github.com/truls/chained-service-example/prodcon_stub"
	"github.com/truls/chained-service-example/producer/impl"

	gen "github.com/truls/chained-service-example/producer/component/gen"

)

type HelloWorldImpl struct{}

func init() {
	a := HelloWorldImpl{}
	gen.SetExportsChainedServiceApiHelloWorld(a)

	c := ProducerConsumerClientImpl{}
	pb_prodcon.SetProducerConsumerClientImplementation(c)
}

func (HelloWorldImpl) InitComponent() {
	impl.Main()
	gen.ChainedServiceApiProdconInitComponent()
}

type ProducerConsumerClientImpl struct {}

func (ProducerConsumerClientImpl) ConsumeByte(ctx context.Context, in *pb_prodcon.ConsumeByteRequest, opts ...interface{}) (*pb_prodcon.ConsumeByteReply, error) {
	var request_param = gen.ChainedServiceApiProdconConsumeByteRequest{Value: in.Value}
	// TODO: Handle error
	res := gen.ChainedServiceApiProdconConsumeByte(request_param).Unwrap()
	return &pb_prodcon.ConsumeByteReply{Value: res.Value, Length: res.Length}, nil
}

func (ProducerConsumerClientImpl) ConsumeStream(ctx context.Context, opts ...interface{}) (pb_prodcon.ProducerConsumer_ConsumeStreamClient, error) {
	return nil, nil
}


func (HelloWorldImpl) SayHello (arg gen.ChainedServiceApiHelloWorldHelloRequest) gen.Result[gen.ChainedServiceApiHelloWorldHelloReply, int32] {
	var param = pb.HelloRequest{Name: arg.Name}
	res, err := pb.Implementation.SayHello(context.TODO(), &param)
	if err != nil {
		fmt.Print("Test failed")
	}


	fmt.Print(res)
	ret := gen.Result[gen.ChainedServiceApiHelloWorldHelloReply, int32]{Kind: gen.Ok, Val: gen.ChainedServiceApiHelloWorldHelloReply{Message: res.Message}, Err: 0}
	return ret

}

//go:generate wit-bindgen tiny-go ../../wit --world producer-interface --out-dir=gen
func main() {}
