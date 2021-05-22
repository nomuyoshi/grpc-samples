package interceptor

import (
	"context"
	"log"
	"time"
	"todo/shared/md"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func XTraceID() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		traceID := md.GetTraceIDFromContext(ctx)
		// ContextにTraceIDを詰めてhandlerに渡す
		ctx = md.AddTraceIDToContext(ctx, traceID)
		// handler()関数がサーバーが実行するRPC
		// RPC実行前にはさみたい処理はhandler()を呼び出す前に書く
		return handler(ctx, req)
	}
}

func XUserID() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		userID, err := md.SafeGetUserIDFromContext(ctx)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		ctx = md.AddUserIDToContext(ctx, userID)
		return handler(ctx, req)
	}
}

const loggingFmt = "TraceID:%s\tFullMethod:%s\tElapsedTime:%s\tStatusCode:%s\tError:%s\n"

func Logging() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// RPC実行前（handler呼び出し前）に時刻を取得
		start := time.Now()
		h, err := handler(ctx, req)
		// RPC実行後（handler呼び出し後）にエラーがあればメッセージを取り出してログ出力
		var errMsg string
		if err != nil {
			errMsg = err.Error()
		}
		log.Printf(loggingFmt,
			md.GetTraceIDFromContext(ctx),
			info.FullMethod,
			time.Since(start),
			status.Code(err), errMsg)
		return h, err
	}
}
