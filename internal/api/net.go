package api

import "net"

type NetListenerServer interface {
	Serve(net.Listener) error
}

func WithNetListenerServer(addr string, handler NetListenerServer) ApiOption {
	return newFuncApiOption(func(a *Api) {
		a.servers = append(a.servers, newNetListenerServer(addr, handler))
	})
}

func newNetListenerServer(addr string, s NetListenerServer) *netListenerServer {
	return &netListenerServer{addr, s}
}

type netListenerServer struct {
	addr string
	s    NetListenerServer
}

func (server *netListenerServer) Serve() error {
	lis, err := net.Listen("tcp", server.addr)
	if err != nil {
		return err
	}

	return server.s.Serve(lis)
}

func (server *netListenerServer) Addr() string {
	return server.addr
}
