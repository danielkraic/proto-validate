package main

// based on https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_server/main.go

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/danielkraic/proto-validate-plugin/example/person"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()

	addr := fmt.Sprintf("0.0.0.0:%d", *port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	person.RegisterGreeterServer(s, &Server{})

	log.Printf("server listening at %v", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type Server struct {
	person.UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, req *person.HelloRequest) (*person.HelloReply, error) {
	return &person.HelloReply{Message: "Hello " + req.Person.Name}, nil
}
