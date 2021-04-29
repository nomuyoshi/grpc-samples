package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/nomuyoshi/grpc-samples/echo/proto"
	"google.golang.org/grpc"
)

func init() {
	log.SetFlags(log.Ldate)
	log.SetPrefix("[echo] ")
}

func main() {
	target := "localhost:50051"
	// 本来はTLSが必要だが、省略
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	client := pb.NewEchoServiceClient(conn)
	msg := os.Args[1]
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.Echo(ctx, &pb.EchoRequest{Message: msg})
	if err != nil {
		log.Println(err)
	}
	log.Println(r.GetMessage())
}
