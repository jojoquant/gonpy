package main

import (
	"context"
	"gonpy/trader/grpc/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct{
	proto.UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{Message: "[gonpy] hello " + req.Name}, nil
}

func main(){
	s := grpc.NewServer()
	proto.RegisterGreeterServer(s, &Server{})

	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}