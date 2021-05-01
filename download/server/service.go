package main

import (
	"io"
	"os"
	"path/filepath"

	pb "github.com/nomuyoshi/grpc-samples/download/proto"
)

// download/proto/file_grpc.pb.go のFileServiceServerインターフェースを満たすように実装する
type fileService struct {
	// 未実装のRPCにはUnimplementedFileServiceServerがエラーコードUnimplementedを返してくれる
	pb.UnimplementedFileServiceServer
}

func (s *fileService) Download(req *pb.FileRequest, stream pb.FileService_DownloadServer) error {
	fp := filepath.Join("./download/resource", req.GetName())
	fs, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer fs.Close()

	buf := make([]byte, 1000*1024)
	for {
		n, err := fs.Read(buf)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		// 読み取った内容を含めた FileResponse を渡してレスポンスを返す
		err = stream.Send(&pb.FileResponse{Data: buf[:n]})
		if err != nil {
			return err
		}
	}

	return nil
}
