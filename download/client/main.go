package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	pb "github.com/nomuyoshi/grpc-samples/download/proto"
	"google.golang.org/grpc"
)

func init() {
	log.SetFlags(log.Ltime)
	log.SetPrefix("[download] ")
}

// go run ./download/client ファイル名
func main() {
	target := "localhost:50051"
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}
	defer conn.Close()
	c := pb.NewFileServiceClient(conn)
	name := os.Args[1]
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stream, err := c.Download(ctx, &pb.FileRequest{Name: name})
	if err != nil {
		log.Fatalf("could not download: %s\n", err)
	}

	var blob []byte
	for {
		c, err := stream.Recv()
		// Recv() はストリーム完了（サーバーからの最後のレスポンス）時には、io.EOFを返す
		if err == io.EOF {
			log.Printf("done %d bytes\n", len(blob))
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		blob = append(blob, c.GetData()...)
	}
	// レスポンスで受け取ったデータをファイルに書き込む
	ioutil.WriteFile(name, blob, 0644)
}
