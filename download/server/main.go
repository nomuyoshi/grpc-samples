package main

import (
	"log"
	"net"

	pb "github.com/nomuyoshi/grpc-samples/download/proto"
	"google.golang.org/grpc"
)

const port = ":50051"

func init() {
	log.SetFlags(log.Ltime)
	log.SetPrefix("[download] ")
}

// go run ./download/server
func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to Listen: %v\n", err)
	}

	srv := grpc.NewServer()
	pb.RegisterFileServiceServer(srv, &fileService{})
	log.Printf("start server on port%s\n", port)
	if err := srv.Serve(listen); err != nil {
		log.Printf("failed to serve: %v\n", err)
	}
}
