module github.com/truls/chained-service-example/producer/component

go 1.20

replace (
	github.com/truls/chained-service-example/helloworld_stub => ../../helloworld_stub
	github.com/truls/chained-service-example/prodcon_stub => ../../prodcon_stub
	github.com/truls/chained-service-example/producer/impl => ../impl
	github.com/truls/chained-service-example/stubs/grpc => ../../stubs/grpc
	github.com/truls/chained-service-example/stubs/net => ../../stubs/net
)

require (
	github.com/truls/chained-service-example/helloworld_stub v0.0.0-00010101000000-000000000000
	github.com/truls/chained-service-example/producer/impl v0.0.0-00010101000000-000000000000
)

require (
	github.com/truls/chained-service-example/prodcon_stub v0.0.0-00010101000000-000000000000 // indirect
	github.com/truls/chained-service-example/stubs/grpc v1.0.0 // indirect
	github.com/truls/chained-service-example/stubs/net v0.0.0-00010101000000-000000000000 // indirect
)
