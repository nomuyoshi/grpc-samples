package main

import (
	"io"
	"log"
	"sync"

	pb "github.com/nomuyoshi/grpc-samples/chat/proto"
)

type chatService struct {
	pb.UnimplementedChatServiceServer
}

var streams sync.Map

func (s *chatService) Connect(stream pb.ChatService_ConnectServer) error {
	log.Println("connect", &stream)
	streams.Store(stream, struct{}{})
	defer func() {
		log.Println("disconnect", &stream)
		streams.Delete(stream)
	}()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		streams.Range(func(key, value interface{}) bool {
			stream := key.(pb.ChatService_ConnectServer)
			stream.Send(&pb.Post{
				Name:    req.GetName(),
				Message: req.GetMessage(),
			})
			return true
		})
	}
}
