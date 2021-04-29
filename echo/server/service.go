package main

import (
	"context"

	pb "github.com/nomuyoshi/grpc-samples/echo/proto"
)

// proto/echo_grpc.pb.go の EchoServiceServer Interface を満たすように実装
type echoService struct {
	// 未実装のRPCにはUnimplementedEchoServiceServerがエラーコードUnimplementedを返してくれる
	pb.UnimplementedEchoServiceServer
}

func (s *echoService) Echo(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{
		Message: req.GetMessage(),
	}, nil
}
