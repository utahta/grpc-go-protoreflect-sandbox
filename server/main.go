package main

import (
	"context"
	"log"
	"net"

	"github.com/utahta/grpc-go-protoreflect-example/gen/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50050"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}

func (s *server) SayHello2(ctx context.Context, in *helloworld.Hello2Request) (*helloworld.Hello2Reply, error) {
	return &helloworld.Hello2Reply{Message: "Hello2 " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	helloworld.RegisterGreeterServer(s, &server{})
	helloworld.RegisterGreeter2Server(s, &server{})

	reflection.Register(s) // for protoreflect

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
