module github.com/truls/chained-service-example/producer

go 1.20

require (
	github.com/containerd/containerd v1.7.1
	github.com/sirupsen/logrus v1.9.2
	//github.com/truls/chained-service-example/helloworld v1.0.0
	github.com/truls/chained-service-example/helloworld_stub v0.0.0-00010101000000-000000000000
	github.com/truls/chained-service-example/stubs/grpc v1.0.0
)

require (
	github.com/truls/chained-service-example/prodcon_stub v0.0.0-00010101000000-000000000000
	github.com/truls/chained-service-example/stubs/net v0.0.0-00010101000000-000000000000
)

require golang.org/x/sys v0.7.0 // indirect

replace (
	//github.com/truls/chained-service-example/helloworld => ../helloworld
	github.com/truls/chained-service-example/helloworld_stub => ../../helloworld_stub
	github.com/truls/chained-service-example/prodcon_stub => ../../prodcon_stub
	github.com/truls/chained-service-example/stubs/grpc => ../../stubs/grpc
	github.com/truls/chained-service-example/stubs/net => ../../stubs/net

)
