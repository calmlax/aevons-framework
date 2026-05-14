package grpcx

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// DefaultDialOptions 返回内部 gRPC 客户端默认拨号参数。
func DefaultDialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(JSONCodec{})),
	}
}

// NewClientConn 使用 framework 约定的默认参数创建 gRPC 客户端连接。
func NewClientConn(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	dialOptions := append(DefaultDialOptions(), opts...)
	return grpc.NewClient(target, dialOptions...)
}
