module github.com/truls/chained-service-example/prodcon_stub

go 1.20

replace github.com/truls/chained-service-example/grpc_stub => ../grpc_stub

require google.golang.org/grpc v1.53.0

require (
	github.com/golang/protobuf v1.5.3 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)
