all: prodcon.pb.go prodcon_grpc.pb.go

prodcon.pb.go:prodcon.proto
	protoc --go_out=.  --go_opt=paths=source_relative $<

prodcon_grpc.pb.go:prodcon.proto
	protoc --go-grpc_out=.  --go-grpc_opt=paths=source_relative $<

clean:
	rm -r prodcon.pb.go prodcon_grpc.pb.go

.PHONY: clean
