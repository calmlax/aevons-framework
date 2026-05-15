package grpcx

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/calmlax/aevons-framework/config"
	"github.com/calmlax/aevons-framework/core/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ErrEmptyTarget 表示 gRPC 目标地址为空。
var ErrEmptyTarget = errors.New("grpc target is empty")

// ErrNoHealthyInstance 表示 Consul 中没有发现可用实例。
var ErrNoHealthyInstance = errors.New("no healthy service instance found")

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

// ResolveTargetFromConsul 从 Consul 中发现服务实例，并返回一个可用于 gRPC 拨号的目标地址。
// 当前策略是：
// 1. 仅选择 Consul 中 passing=true 的健康实例
// 2. 从健康实例中随机选一个，作为最小可用的客户端负载均衡
// 3. 优先使用实例元信息里的 grpc_port；若未提供，则退回实例主端口
func ResolveTargetFromConsul(consulCfg config.ConsulConfig, serviceName string) (string, error) {
	registry, err := consul.New(consulCfg)
	if err != nil {
		return "", err
	}

	instances, err := registry.Discover(serviceName)
	if err != nil {
		return "", err
	}
	if len(instances) == 0 {
		return "", fmt.Errorf("%w: %s", ErrNoHealthyInstance, serviceName)
	}

	chosen := instances[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(instances))]
	port := chosen.GRPCPort
	if port <= 0 {
		port = chosen.Port
	}
	if chosen.Address == "" || port <= 0 {
		return "", fmt.Errorf("invalid consul instance for %s", serviceName)
	}

	return fmt.Sprintf("%s:%d", chosen.Address, port), nil
}

// NewClientConnFromConsul 通过 Consul 服务发现创建 gRPC 客户端连接。
// 业务侧只需要提供 Consul 配置和目标服务名，不必再写死实例 IP:端口。
func NewClientConnFromConsul(consulCfg config.ConsulConfig, serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	target, err := ResolveTargetFromConsul(consulCfg, serviceName)
	if err != nil {
		return nil, err
	}
	return NewClientConn(target, opts...)
}

// CloseClientConn 安全关闭 gRPC 客户端连接。
func CloseClientConn(conn *grpc.ClientConn) error {
	if conn == nil {
		return nil
	}
	return conn.Close()
}
