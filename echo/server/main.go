package main

import (
	"log"
	"net"

	pb "github.com/nomuyoshi/grpc-samples/echo/proto"
	"google.golang.org/grpc"
)

func init() {
	log.SetFlags(log.Ldate)
	log.SetPrefix("[echo] ")
}

func main() {
	port := ":50051"
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	srv := grpc.NewServer()
	pb.RegisterEchoServiceServer(srv, &echoService{})
	log.Printf("start server on port%s\n", port)
	if err := srv.Serve(listen); err != nil {
		log.Printf("failed to serve: %v\n", err)
	}
}
