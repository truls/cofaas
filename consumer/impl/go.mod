module consumer

go 1.20

require (
	github.com/truls/chained-service-example/stubs/grpc v0.0.0-00010101000000-000000000000
	github.com/truls/chained-service-example/stubs/net v0.0.0-00010101000000-000000000000
	github.com/truls/chained-service-example/prodcon_stub v0.0.0-00010101000000-000000000000
)

replace (
	github.com/truls/chained-service-example/stubs/grpc => ../../stubs/grpc
	//github.com/truls/chained-service-example/helloworld => ../helloworld
	github.com/truls/chained-service-example/helloworld_stub => ../../helloworld_stub
	github.com/truls/chained-service-example/stubs/net => ../../stubs/net
	//github.com/truls/chained-service-example/prodcon => ../prodcon
	github.com/truls/chained-service-example/prodcon_stub => ../../prodcon_stub
//google.golang.com/grpc => ../stubs/grpc
)
