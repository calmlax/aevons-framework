package consul

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
}

// Registry 提供 Consul 注册、注销和发现能力。
type Registry struct {
	baseURL    string
	httpClient *http.Client
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

// Register 将当前服务实例注册到 Consul。
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
			ID      string `json:"ID"`
			Service string `json:"Service"`
			Address string `json:"Address"`
			Port    int    `json:"Port"`
		} `json:"Service"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, err
	}

	instances := make([]Instance, 0, len(entries))
	for _, entry := range entries {
		instances = append(instances, Instance{
			ID:      entry.Service.ID,
			Name:    entry.Service.Service,
			Address: entry.Service.Address,
			Port:    entry.Service.Port,
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

func registrationAddress(host string) string {
	switch strings.TrimSpace(host) {
	case "", "0.0.0.0", "::", "[::]":
		return "127.0.0.1"
	default:
		return strings.TrimSpace(host)
	}
}
