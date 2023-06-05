package grpc

import "context"

// CallOption configures a Call before it starts or extracts information from
// a Call after it completes.
type CallOption interface {}

type ClientConn struct {}

func (x *ClientConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...interface{}) error {
	return nil
}

func (x *ClientConn) NewStream(ctx context.Context, desc *interface{}, method string, opts ...interface{}) (interface{}, error) {
	return nil, nil
}

func (x *ClientConn) Close() {
}

func WithBlock() interface{} {
	return nil
}

func Dial(target string, opts ...interface{}) (*ClientConn, error) {
	return &ClientConn{}, nil
}

func WithTransportCredentials(creds interface{}) interface{} {
	return nil//grpc.WithTransportCredentials(creds)
}

type Server struct {}

func NewServer() (*Server) {
	return &Server{}
}

func (*Server) Serve(lis interface{}) error {
	return nil
}
