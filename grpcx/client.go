package grpcx

import (
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ErrEmptyTarget 表示 gRPC 目标地址为空。
var ErrEmptyTarget = errors.New("grpc target is empty")

// DefaultDialOptions 返回内部 gRPC 客户端默认拨号参数。
func DefaultDialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(JSONCodec{})),
	}
}

// NewClientConn 使用 framework 约定的默认参数创建 gRPC 客户端连接。
func NewClientConn(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	if target == "" {
		return nil, ErrEmptyTarget
	}
	dialOptions := append(DefaultDialOptions(), opts...)
	return grpc.NewClient(target, dialOptions...)
}

// CloseClientConn 安全关闭 gRPC 客户端连接。
func CloseClientConn(conn *grpc.ClientConn) error {
	if conn == nil {
		return nil
	}
	return conn.Close()
}
