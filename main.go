package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpcClient/proto/article"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:9095", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// 初始化客户端
	c := article.NewArticleClient(conn)

	// 调用方法
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*1)
	res, err := c.GetArticleById(ctx, &article.ArticleRequest{Id: 1})

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(res)
}
