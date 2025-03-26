package main

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

// LocalServer 封装了一个内存 gRPC 服务及其 Listener
type LocalServer struct {
	lis    *bufconn.Listener
	Server *grpc.Server
}

// NewLocalServer 创建一个内存 gRPC 服务，传入可选的 grpc.ServerOption
func NewLocalServer(opts ...grpc.ServerOption) *LocalServer {
	lis := bufconn.Listen(bufSize)
	srv := grpc.NewServer(opts...)
	return &LocalServer{
		lis:    lis,
		Server: srv,
	}
}

// Start 在一个独立的 goroutine 中启动 gRPC 服务
func (ls *LocalServer) Start() error {
	go func() {
		if err := ls.Server.Serve(ls.lis); err != nil {
			// 正常情况下这里可以记录日志或通知错误
		}
	}()
	return nil
}

// Stop 停止 gRPC 服务
func (ls *LocalServer) Stop() {
	ls.Server.Stop()
}

// Dialer 返回一个自定义的 Dialer，可供 grpc.DialContext 使用，连接到内存中的 Listener
func (ls *LocalServer) Dialer(ctx context.Context, s string) (net.Conn, error) {
	return ls.lis.Dial()
}
