module github.com/truls/chained-service-example/producer

go 1.20

require (
	github.com/containerd/containerd v1.7.1
	github.com/sirupsen/logrus v1.9.2
	github.com/truls/chained-service-example/grpc_stub v1.0.0
	//github.com/truls/chained-service-example/helloworld v1.0.0
	github.com/truls/chained-service-example/helloworld_stub v0.0.0-00010101000000-000000000000
)

require (
	github.com/truls/chained-service-example/net_stub v0.0.0-00010101000000-000000000000
	github.com/truls/chained-service-example/prodcon_stub v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.55.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/sys v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace (
	github.com/truls/chained-service-example/grpc_stub => ../grpc_stub
	//github.com/truls/chained-service-example/helloworld => ../helloworld
	github.com/truls/chained-service-example/helloworld_stub => ../helloworld_stub
	github.com/truls/chained-service-example/net_stub => ../net_stub
	github.com/truls/chained-service-example/prodcon_stub => ../prodcon_stub

)
