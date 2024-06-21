package api

import "net/http"

type HttpServer interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func WithHttpHandler(addr string, handler HttpServer) ApiOption {
	return newFuncApiOption(func(a *Api) {
		a.servers = append(a.servers, newHttpServer(addr, handler))
	})
}

func newHttpServer(addr string, handler HttpServer) *httpServer {
	return &httpServer{addr, handler}
}

type httpServer struct {
	addr    string
	handler HttpServer
}

func (server *httpServer) Serve() error {
	return http.ListenAndServe(server.addr, server.handler)
}

func (server *httpServer) Addr() string {
	return server.addr
}
