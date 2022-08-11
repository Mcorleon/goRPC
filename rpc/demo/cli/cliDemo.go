package main

import (
	"context"
	"encoding/gob"
	api "grpcClient/rpc/api/user"
	"grpcClient/rpc/client"
	"log"
	"time"
)

func main() {
	gob.Register(api.User{})
	gob.Register([]interface{}{})
	cli := client.NewClientProxy(client.DefaultOption)
	ctx, _ := context.WithTimeout(context.Background(),time.Second)
	var GetUserById func(id int) (api.User, error)
	cli.Call(ctx, "UserService.User.GetUserById", &GetUserById)
	u, err := GetUserById(541)
	log.Println("result:", u, err)

}
