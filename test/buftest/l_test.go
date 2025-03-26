package main

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/air-iot/api-client-go/v4/api"
	pb "github.com/air-iot/api-client-go/v4/test"
	"google.golang.org/grpc"
)

// TestProvider 为你的 gRPC 服务实现
type TestProvider struct {
	pb.UnimplementedFlowTaskServiceServer
}

func (a *TestProvider) Get(ctx context.Context, req *api.GetOrDeleteRequest) (*api.Response, error) {
	// 这里可以添加更复杂的逻辑，当前仅模拟返回数据
	return &api.Response{
		Status: true,
		Code:   200,
		Info:   "ok",
		Result: []byte("===================================="),
	}, nil
}

// setupLocalServer 用于基准测试时初始化本地 gRPC 服务及客户端
func setupLocalServer() (pb.FlowTaskServiceClient, context.Context, func()) {
	ls := NewLocalServer()
	// 注册服务实现
	pb.RegisterFlowTaskServiceServer(ls.Server, &TestProvider{})
	ls.Start()

	// 使用较长的超时，避免基准测试过程中因超时失败
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(ls.Dialer), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial bufnet: %v", err)
	}

	client := pb.NewFlowTaskServiceClient(conn)

	cleanup := func() {
		cancel()
		conn.Close()
		ls.Stop()
	}

	return client, ctx, cleanup
}

// BenchmarkGet 对 Get 方法进行性能基准测试
func BenchmarkGet(b *testing.B) {
	client, ctx, cleanup := setupLocalServer()
	defer cleanup()

	req := &api.GetOrDeleteRequest{Id: "test"}

	// 重置计时器，开始正式测量
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.Get(ctx, req)
		if err != nil {
			b.Fatalf("client.Get error: %v", err)
		}
	}
}
