module github.com/truls/chained-service-example/producer_component

go 1.20

replace (
	github.com/truls/chained-service-example/helloworld_stub => ../../helloworld_stub
	github.com/truls/chained-service-example/prodcon_stub => ../../prodcon_stub
	github.com/truls/chained-service-example/producer/impl => ../impl
	github.com/truls/chained-service-example/stubs/grpc => ../../stubs/grpc
	github.com/truls/chained-service-example/stubs/net => ../../stubs/net
)

require github.com/truls/chained-service-example/producer/impl v0.0.0-00010101000000-000000000000

require (
	github.com/containerd/containerd v1.7.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/sirupsen/logrus v1.9.2 // indirect
	github.com/truls/chained-service-example/helloworld_stub v0.0.0-00010101000000-000000000000 // indirect
	github.com/truls/chained-service-example/prodcon_stub v0.0.0-00010101000000-000000000000 // indirect
	github.com/truls/chained-service-example/stubs/grpc v1.0.0 // indirect
	github.com/truls/chained-service-example/stubs/net v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/sys v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
	google.golang.org/grpc v1.55.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)
