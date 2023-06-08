package main

import (
	"testing"
	gen "github.com/truls/chained-service-example/consumer/component/gen" )


func TestSayHello_Wrapper(t *testing.T) {
	a := HelloWorldImpl{}
	a.InitComponent()
	ret := a.SayHello(gen.HelloWorldHelloRequest{Name: "foo"})
	expected := gen.Result[gen.HelloWorldHelloReply, int32]{Kind: gen.Ok, Val: gen.HelloWorldHelloReply{Message: "Success"}, Err: 0}
	if ret.IsErr() {
		t.Fatalf("Call failed %s\n", ret.UnwrapErr())
	} else {
		if ret != expected {
			t.Fatal("Wrong result\n")
		}
	}
}
