package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// LogoutPayload Back-channel Logout 通知的请求体。
type LogoutPayload struct {
	UId       int64  `json:"uid"`
	Action    string `json:"action"` // "logout_all"
	Timestamp int64  `json:"timestamp"`
}

// SLONotifier 定义 Back-channel Logout 通知接口。
type SLONotifier interface {
	// Notify 向指定 URL 发送 Back-channel Logout 通知。
	// 发送失败时记录日志但不返回错误，不阻断登出流程。
	Notify(ctx context.Context, webhookURL string, payload *LogoutPayload)
}

// HttpSLONotifier 基于 HTTP POST 实现 SLONotifier。
type HttpSLONotifier struct {
	client *http.Client
}

// NewHttpSLONotifier 创建 HttpSLONotifier，超时 5 秒。
func NewHttpSLONotifier() *HttpSLONotifier {
	return &HttpSLONotifier{
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

// Notify 向 webhookURL 发送 POST 请求，失败时仅记录日志。
func (n *HttpSLONotifier) Notify(ctx context.Context, webhookURL string, payload *LogoutPayload) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("slo_notifier: 序列化 LogoutPayload 失败: %v", err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewReader(data))
	if err != nil {
		log.Printf("slo_notifier: 创建请求失败 url=%s err=%v", webhookURL, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		log.Printf("slo_notifier: 发送通知失败 url=%s err=%v", webhookURL, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Printf("slo_notifier: 通知返回异常状态码 url=%s status=%d", webhookURL, resp.StatusCode)
		return
	}

	log.Printf("slo_notifier: 通知发送成功 url=%s", webhookURL)
}

// NoopSLONotifier 空实现，用于不需要 SLO 通知的场景（测试或单机部署）。
type NoopSLONotifier struct{}

func (n *NoopSLONotifier) Notify(_ context.Context, webhookURL string, _ *LogoutPayload) {
	log.Printf("slo_notifier: noop，跳过通知 url=%s", webhookURL)
}
