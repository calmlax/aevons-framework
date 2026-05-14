package middleware

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/mileusna/useragent"
)

// ResponseBodyWriter 包装gin.ResponseWriter，用于捕获响应体大小
type ResponseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *ResponseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger 增强版日志中间件，输出丰富的请求/响应信息
// 包含：请求Id、IP、方法、路径、状态码、耗时、客户端信息、请求/响应体大小等
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 初始化基础信息
		start := time.Now()
		req := c.Request
		// 获取真实客户端IP（处理反向代理/CDN场景）
		clientIP := getClientIP(c)
		// 获取客户端User-Agent（浏览器/系统信息）
		userAgent := req.UserAgent()
		// 解析User-Agent为易读格式
		clientInfo := parseUserAgent(userAgent)
		// 获取请求Id（需配合gin-contrib/requestid中间件）
		reqId := requestid.Get(c)
		if reqId == "" {
			reqId = fmt.Sprintf("req-%d", time.Now().UnixNano()) // 兜底生成
		}

		// 2. 捕获响应体大小
		w := &ResponseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = w

		// 3. 处理请求
		c.Next()

		// 4. 计算耗时（毫秒级，更易读）
		latency := time.Since(start)
		latencyMs := float64(latency.Nanoseconds()) / 1e6 // 转毫秒

		// 5. 提取请求/响应信息
		method := req.Method
		path := req.URL.Path
		statusCode := c.Writer.Status()
		reqSize := req.ContentLength // 请求体大小（字节）
		respSize := w.body.Len()     // 响应体大小（字节）

		// 6. 结构化日志输出（易解析，支持ELK等日志系统）
		log.Printf(
			"[GIN] [reqId:%s] [IP:%s] [Client:%s] [Method:%s] [Path:%s] [Status:%d] [Latency:%.2fms] [ReqSize:%dB] [RespSize:%dB]",
			reqId,
			clientIP,
			clientInfo,
			method,
			path,
			statusCode,
			latencyMs,
			reqSize,
			respSize,
		)
	}
}

// getClientIP 获取真实客户端IP（处理X-Forwarded-For等反向代理头）
func getClientIP(c *gin.Context) string {
	// 优先从X-Forwarded-For头获取（反向代理/CDN场景）
	xForwardedFor := c.Request.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// 格式：X-Forwarded-For: clientIP, proxy1, proxy2
		ips := strings.Split(xForwardedFor, ",")
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			if ip != "" && !isInternalIP(ip) {
				return ip
			}
		}
	}

	// 其次从X-Real-IP头获取
	xRealIP := c.Request.Header.Get("X-Real-IP")
	if xRealIP != "" && !isInternalIP(xRealIP) {
		return xRealIP
	}

	// 最后从TCP连接获取
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return ip
}

// isInternalIP 判断是否为内网IP（过滤代理IP）
func isInternalIP(ip string) bool {
	ipv4 := net.ParseIP(ip)
	if ipv4 == nil {
		return false
	}
	// 内网IP段：10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16, 127.0.0.1/32
	return ipv4.IsPrivate() || ipv4.Equal(net.IPv4(127, 0, 0, 1))
}

func parseUserAgent(ua string) string {
	if ua == "" {
		return "Unknown/Unknown"
	}
	client := useragent.Parse(ua)
	os := client.OS
	if os == "" {
		os = "UnknownOS"
	}
	browser := client.Name
	if browser == "" {
		browser = "UnknownBrowser"
	}
	return fmt.Sprintf("%s/%s", os, browser)
}
