package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	pbProject "todo/proto/project"
	pbUser "todo/proto/user"
	db "todo/shared/db"
	"todo/shared/interceptor"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"google.golang.org/grpc"
)

const port = ":50051"

func main() {
	dbConn := db.ConnectDB()
	projectConn, err := grpc.Dial(os.Getenv("PROJECT_SERVICE_ADDR"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial to project: %s", err)
	}

	srv := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		interceptor.XTraceID(),
		interceptor.XUserID(),
		interceptor.Logging(),
	)))
	pbUser.RegisterUserServiceServer(srv, &userService{
		db:            dbConn,
		projectClient: pbProject.NewProjectServiceClient(projectConn),
	})

	go func() {
		listener, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to create listener: %s", err)
		}
		log.Println("start server on port", port)

		if err := srv.Serve(listener); err != nil {
			log.Println("failed to exit serve: ", err)
		}
	}()

	// グレースフルストップ
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM)
	// sigint チャネルに値が送信されるまで(=受信するまで)待機
	<-sigint
	// シグナルを受け取ったら graceful shutdown 開始
	log.Println("received a signal of graceful shutdown")
	stopped := make(chan struct{})
	go func() {
		// GracefulStop は、新しいgRPC接続とRPCの受け付けをやめて、処理中のすべてのRPCが終了するまでブロック
		srv.GracefulStop()
		close(stopped)
	}()
	// タイムアウト1min
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	select {
	case <-ctx.Done():
		// タイムアウト（1min経過）の場合Stopする
		srv.Stop()
	case <-stopped:
		cancel()
	}
	log.Println("completed graceful shutdown")
}
