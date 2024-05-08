module cofaas/tests/consumer

go 1.20

require (
	cofaas_orig/protos/prodcon v0.0.0-00010101000000-000000000000
	github.com/containerd/containerd v1.7.3
	github.com/sirupsen/logrus v1.9.4-0.20230606125235-dd1b4c2e81af
	google.golang.org/grpc v1.57.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230525234030-28d5490b6b19 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

replace cofaas_orig/protos/prodcon => ../protos/prodcon
