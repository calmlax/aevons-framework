package grpcx

import "google.golang.org/grpc"

// DefaultServerOptions 返回内部 gRPC 服务端默认启动参数。
func DefaultServerOptions() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ForceServerCodec(JSONCodec{}),
	}
}

// NewServer 使用 framework 约定的默认参数创建 gRPC 服务端。
func NewServer(opts ...grpc.ServerOption) *grpc.Server {
	serverOptions := append(DefaultServerOptions(), opts...)
	return grpc.NewServer(serverOptions...)
}
