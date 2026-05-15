package consul

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/calmlax/aevons-framework/config"
)

const DefaultHealthPath = "/health"

// Instance 描述一个 Consul 服务实例。
type Instance struct {
	ID      string
	Name    string
	Address string
	Port    int
	// GRPCPort 来自 Consul Service.Meta.grpc_port。
	// 当服务同时暴露 HTTP 与 gRPC 端口时，调用方应优先使用这个端口建立 gRPC 连接。
	GRPCPort int
}

// Registry 提供 Consul 注册、注销和发现能力。
type Registry struct {
	baseURL    string
	httpClient *http.Client
}

// Managed 负责将健康等待、注册和注销封装为统一生命周期。
type Managed struct {
	registry   *Registry
	serverCfg  config.ServerConfig
	healthPath string
	serviceID  string
}

// New 根据配置创建 Consul 注册中心客户端。
func New(cfg config.ConsulConfig) (*Registry, error) {
	baseURL, err := normalizeAddress(cfg.Address)
	if err != nil {
		return nil, err
	}
	return &Registry{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}, nil
}

// NewManaged 创建一个带注册生命周期管理的 Consul 帮助对象。
func NewManaged(consulCfg config.ConsulConfig, serverCfg config.ServerConfig, healthPath string) (*Managed, error) {
	registry, err := New(consulCfg)
	if err != nil {
		return nil, err
	}
	if healthPath == "" {
		healthPath = DefaultHealthPath
	}
	return &Managed{
		registry:   registry,
		serverCfg:  serverCfg,
		healthPath: healthPath,
	}, nil
}

// Register 等待健康接口就绪后再执行 Consul 注册。
func (m *Managed) Register(timeout time.Duration) error {
	if m == nil || m.registry == nil {
		return nil
	}

	if err := waitForHealth(m.serverCfg.Host, m.serverCfg.Port, m.healthPath, timeout); err != nil {
		return err
	}

	serviceID, err := m.registry.Register(m.serverCfg, m.healthPath)
	if err != nil {
		return err
	}
	m.serviceID = serviceID
	return nil
}

// Deregister 注销当前已注册的服务实例。
func (m *Managed) Deregister() error {
	if m == nil || m.registry == nil {
		return nil
	}
	return m.registry.Deregister(m.serviceID)
}

// Discover 查询当前服务名对应的健康实例。
func (m *Managed) Discover() ([]Instance, error) {
	if m == nil || m.registry == nil {
		return nil, nil
	}
	return m.registry.Discover(m.serverCfg.Name)
}

// Register 将当前服务实例注册到 Consul。
// 除了注册 HTTP 健康检查地址外，还会把 gRPC 端口写入 Service.Meta，
// 供其他服务通过 Consul 发现后建立 gRPC 连接。
func (r *Registry) Register(serverCfg config.ServerConfig, healthPath string) (string, error) {
	if r == nil {
		return "", nil
	}

	serviceName := strings.TrimSpace(serverCfg.Name)
	if serviceName == "" {
		serviceName = "service"
	}

	address := registrationAddress(serverCfg.Host)
	if serverCfg.Port <= 0 {
		return "", fmt.Errorf("invalid service port: %d", serverCfg.Port)
	}

	if healthPath == "" {
		healthPath = DefaultHealthPath
	}
	if !strings.HasPrefix(healthPath, "/") {
		healthPath = "/" + healthPath
	}

	serviceID := fmt.Sprintf("%s-%s-%d", serviceName, address, serverCfg.Port)
	registration := map[string]any{
		"ID":      serviceID,
		"Name":    serviceName,
		"Address": address,
		"Port":    serverCfg.Port,
		"Meta": map[string]string{
			"grpc_port": fmt.Sprintf("%d", serverCfg.GRPCPort),
		},
		"Check": map[string]any{
			"HTTP":                           fmt.Sprintf("http://%s:%d%s", address, serverCfg.Port, healthPath),
			"Method":                         "GET",
			"Interval":                       "10s",
			"Timeout":                        "3s",
			"DeregisterCriticalServiceAfter": "1m",
		},
	}

	if err := r.put("/v1/agent/service/register", registration); err != nil {
		return "", err
	}

	return serviceID, nil
}

// Deregister 将服务实例从 Consul 注销。
func (r *Registry) Deregister(serviceID string) error {
	if r == nil || serviceID == "" {
		return nil
	}
	req, err := http.NewRequest(http.MethodPut, r.baseURL+"/v1/agent/service/deregister/"+serviceID, nil)
	if err != nil {
		return err
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("consul deregister failed: %s", resp.Status)
	}
	return nil
}

// Discover 查询指定服务当前健康实例。
// 如果实例在注册时带了 grpc_port 元信息，这里会一并解析出来。
func (r *Registry) Discover(serviceName string) ([]Instance, error) {
	if r == nil {
		return nil, nil
	}

	resp, err := r.httpClient.Get(r.baseURL + "/v1/health/service/" + serviceName + "?passing=true")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("consul discover failed: %s", resp.Status)
	}

	var entries []struct {
		Service struct {
			ID      string            `json:"ID"`
			Service string            `json:"Service"`
			Address string            `json:"Address"`
			Port    int               `json:"Port"`
			Meta    map[string]string `json:"Meta"`
		} `json:"Service"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, err
	}

	instances := make([]Instance, 0, len(entries))
	for _, entry := range entries {
		grpcPort := 0
		if rawPort := entry.Service.Meta["grpc_port"]; rawPort != "" {
			if parsedPort, convErr := strconv.Atoi(rawPort); convErr == nil {
				grpcPort = parsedPort
			}
		}
		instances = append(instances, Instance{
			ID:       entry.Service.ID,
			Name:     entry.Service.Service,
			Address:  entry.Service.Address,
			Port:     entry.Service.Port,
			GRPCPort: grpcPort,
		})
	}
	return instances, nil
}

func (r *Registry) put(path string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, r.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("consul request failed: %s", resp.Status)
	}
	return nil
}

// normalizeAddress 将配置中的 Consul 地址标准化为可直接发 HTTP 请求的 baseURL。
func normalizeAddress(rawAddress string) (string, error) {
	address := strings.TrimSpace(rawAddress)
	if address == "" {
		return "http://127.0.0.1:8500", nil
	}

	if strings.Contains(address, "://") {
		u, err := url.Parse(address)
		if err != nil {
			return "", fmt.Errorf("parse consul address %q: %w", rawAddress, err)
		}
		return strings.TrimRight(u.String(), "/"), nil
	}

	return "http://" + address, nil
}

// registrationAddress 生成注册到 Consul 的服务地址。
// 对 0.0.0.0、:: 这类监听地址，会回退为 127.0.0.1，避免注册出不可访问的通配地址。
func registrationAddress(host string) string {
	switch strings.TrimSpace(host) {
	case "", "0.0.0.0", "::", "[::]":
		return "127.0.0.1"
	default:
		return strings.TrimSpace(host)
	}
}

// waitForHealth 在注册到 Consul 前等待本地健康检查接口 ready，
// 避免服务刚启动时就因为探测过早而被标记为不健康。
func waitForHealth(host string, port int, path string, timeout time.Duration) error {
	address := host
	switch address {
	case "", "0.0.0.0", "::", "[::]":
		address = "127.0.0.1"
	}

	deadline := time.Now().Add(timeout)
	checkURL := fmt.Sprintf("http://%s:%d%s", address, port, path)
	client := &http.Client{Timeout: 1 * time.Second}

	for time.Now().Before(deadline) {
		resp, err := client.Get(checkURL)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
				return nil
			}
		}
		time.Sleep(200 * time.Millisecond)
	}

	return errors.New("timeout waiting for health endpoint")
}
