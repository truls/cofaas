package producer_component

import (
	context "context"
	"fmt"
	pb "github.com/truls/chained-service-example/helloworld_stub"
	"github.com/truls/chained-service-example/producer/impl"

	gen "github.com/truls/chained-service-example/producer/component/gen"

)

type HelloWorldImpl struct{}

func init() {
	a := HelloWorldImpl{}
	gen.SetHelloWorld(a)
}

func (HelloWorldImpl) InitComponent() {
	impl.Main()
}

func (HelloWorldImpl) SayHello(arg gen.HelloWorldHelloRequest) gen.Result[gen.HelloWorldHelloReply, int32] {
	var param = pb.HelloRequest{Name: arg.Name}
	res, err := pb.Implementation.SayHello(context.TODO(), &param)
	if err != nil {
		fmt.Print("Test failed")
	}

	fmt.Print(res)
	ret := gen.Result[gen.HelloWorldHelloReply, int32]{Kind: gen.Ok, Val: gen.HelloWorldHelloReply{Message: res.Message}, Err: 0}
	return ret

}
