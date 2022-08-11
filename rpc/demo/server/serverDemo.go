package main

import (
	"encoding/gob"
	api "grpcClient/rpc/api/user"
	"grpcClient/rpc/provider"
	"os"
	"os/signal"
	"syscall"
)




func main() {
	srv := provider.NewRPCServer("127.0.0.1", 8811)
	srv.RegisterName("User", &api.UserHandler{})
	gob.Register(api.User{})
	gob.Register([]interface{}{})
	go srv.Run()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-quit
	srv.Close()
}
