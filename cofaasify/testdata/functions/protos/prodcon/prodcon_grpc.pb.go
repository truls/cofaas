// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.6
// source: prodcon.proto

package prodcon

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ProducerConsumerClient is the client API for ProducerConsumer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProducerConsumerClient interface {
	ConsumeByte(ctx context.Context, in *ConsumeByteRequest, opts ...grpc.CallOption) (*ConsumeByteReply, error)
}

type producerConsumerClient struct {
	cc grpc.ClientConnInterface
}

func NewProducerConsumerClient(cc grpc.ClientConnInterface) ProducerConsumerClient {
	return &producerConsumerClient{cc}
}

func (c *producerConsumerClient) ConsumeByte(ctx context.Context, in *ConsumeByteRequest, opts ...grpc.CallOption) (*ConsumeByteReply, error) {
	out := new(ConsumeByteReply)
	err := c.cc.Invoke(ctx, "/prodcon.ProducerConsumer/ConsumeByte", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProducerConsumerServer is the server API for ProducerConsumer service.
// All implementations must embed UnimplementedProducerConsumerServer
// for forward compatibility
type ProducerConsumerServer interface {
	ConsumeByte(context.Context, *ConsumeByteRequest) (*ConsumeByteReply, error)
	mustEmbedUnimplementedProducerConsumerServer()
}

// UnimplementedProducerConsumerServer must be embedded to have forward compatible implementations.
type UnimplementedProducerConsumerServer struct {
}

func (UnimplementedProducerConsumerServer) ConsumeByte(context.Context, *ConsumeByteRequest) (*ConsumeByteReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConsumeByte not implemented")
}
func (UnimplementedProducerConsumerServer) mustEmbedUnimplementedProducerConsumerServer() {}

// UnsafeProducerConsumerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProducerConsumerServer will
// result in compilation errors.
type UnsafeProducerConsumerServer interface {
	mustEmbedUnimplementedProducerConsumerServer()
}

func RegisterProducerConsumerServer(s grpc.ServiceRegistrar, srv ProducerConsumerServer) {
	s.RegisterService(&ProducerConsumer_ServiceDesc, srv)
}

func _ProducerConsumer_ConsumeByte_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConsumeByteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProducerConsumerServer).ConsumeByte(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/prodcon.ProducerConsumer/ConsumeByte",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProducerConsumerServer).ConsumeByte(ctx, req.(*ConsumeByteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ProducerConsumer_ServiceDesc is the grpc.ServiceDesc for ProducerConsumer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProducerConsumer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "prodcon.ProducerConsumer",
	HandlerType: (*ProducerConsumerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ConsumeByte",
			Handler:    _ProducerConsumer_ConsumeByte_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "prodcon.proto",
}
