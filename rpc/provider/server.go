package provider

import (
	"log"
	"reflect"
)

type Server interface {
	Register(string, interface{})
	Run()
	Close()
}

type RPCServer struct {
	listener Listener
}
func NewRPCServer(ip string, port int) *RPCServer {
	return &RPCServer{
		listener: NewRPCListener(ip, port),
	}
}
func (svr *RPCServer) Run() {
	go svr.listener.Run()
}
func (svr *RPCServer) Close() {
	if svr.listener != nil {
		svr.listener.Close()
	}
}

func (svr *RPCServer) RegisterName(name string, class interface{}) {
	handler := &RPCServerHandler{class: reflect.ValueOf(class)}
	svr.listener.SetHandler(name, handler)
	log.Printf("%s registered success!\n", name)
}
