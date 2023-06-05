all: helloworld.pb.go helloworld_grpc.pb.go helloworld_vproto.pb.go

helloworld.pb.go:helloworld.proto
	protoc --go_out=. --go_opt=paths=source_relative helloworld.proto

helloworld_grpc.pb.go:helloworld.proto
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative helloworld.proto

helloworld_vproto.pb.go:helloworld.proto
	protoc --go-vtproto_out=. --plugin protoc-gen-go-vtproto="$$(which protoc-gen-go-vtproto)" --go-vtproto_opt=features=marshal+unmarshal+size,paths=source_relative helloworld.proto

clean:
	rm -r helloworld.pb.go helloworld_grpc.pb.go helloworld_vtproto.pb.go

.PHONY: clean
