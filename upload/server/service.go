package main

import (
	"io"
	"io/ioutil"
	"log"
	"path/filepath"

	pb "github.com/nomuyoshi/grpc-samples/upload/proto"
)

// upload/proto/file_grpc.pb.go のFileServiceServerインターフェースを満たすように実装する
type fileService struct {
	// 未実装のRPCにはUnimplementedFileServiceServerがエラーコードUnimplementedを返してくれる
	pb.UnimplementedFileServiceServer
}

func (s *fileService) Upload(stream pb.FileService_UploadServer) error {
	var blob []byte
	var name string
	for {
		c, err := stream.Recv()
		// stream.Recv() はクライアントからの最後のリクエストのとき、io.EOFを返す
		if err == io.EOF {
			log.Printf("done %d bytes\n", len(blob))
			break
		}
		if err != nil {
			panic(err)
		}
		name = c.GetName()
		blob = append(blob, c.GetData()...)
	}
	fp := filepath.Join("./upload/resource", name)
	ioutil.WriteFile(fp, blob, 0644)
	// クライアントにレスポンスを返す
	stream.SendAndClose(&pb.FileResponse{Size: int64(len(blob))})

	return nil
}
