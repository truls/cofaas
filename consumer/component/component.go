package main

import (
	context "context"
	"fmt"
	pb "github.com/truls/chained-service-example/prodcon_stub"
	"github.com/truls/chained-service-example/consumer/impl"

	gen "github.com/truls/chained-service-example/producer/component/gen"

)

type ProdconImpl struct{}

func init() {
	a := ProdconImpl{}
	gen.SetExportsChainedServiceApiProdcon(a)
}

func (ProdconImpl) InitComponent() {
	impl.Main()
}

func (ProdconImpl) ConsumeByte (arg gen.ChainedServiceApiProdconConsumeByteRequest) gen.Result[gen.ChainedServiceApiProdconConsumeByteReply, int32] {
	var param = pb.ConsumeByteRequest{Value: arg.Value}
	res, err := pb.Implementation.ConsumeByte(context.TODO(), &param)
	if err != nil {
		fmt.Print("Test failed")
	}

	ret := gen.Result[gen.ChainedServiceApiProdconConsumeByteReply, int32]{
		Kind: gen.Ok,
		Val: gen.ChainedServiceApiProdconConsumeByteReply{Value: res.Value, Length: res.Length},
		Err: 0,
	}
	return ret

}

//go:generate wit-bindgen tiny-go ../../wit --world consumer-interface --out-dir=gen
func main() {}
