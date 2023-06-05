package net_stub

import "net"

type ListenerImpl struct {
}

type Addr struct {
}


func (*ListenerImpl) Acccept() (net.Conn, error) {
	return nil, nil
}


func (*ListenerImpl) Close() (error) {
	return nil
}

func (*ListenerImpl) Addr() (Addr) {
	return Addr{}
}

func Listen (a string, b string) (ListenerImpl, error) {
	return ListenerImpl{}, nil
}
