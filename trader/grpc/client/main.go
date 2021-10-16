package main

import (
	"context"
	"gonpy/trader/grpc/proto"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err!=nil{
		panic(err)
	}

	defer conn.Close()

	c := proto.NewGreeterClient(conn)
	r, err := c.SayHello(context.Background(), &proto.HelloRequest{Name: "fangyang"})
	if err!=nil{
		panic(err)
	}
	log.Println(r.Message)
}