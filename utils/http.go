package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/calmlax/aevons-framework/response"

	"github.com/gin-gonic/gin"
)

// ResponseBodyWriter 包装 gin.ResponseWriter，用于捕获响应体内容。
type ResponseBodyWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w *ResponseBodyWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

// BadRequest 返回 400 响应
func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, response.Response{
		Code:    400,
		Message: msg,
	})
}

// Unauthorized 返回 401 响应
func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, response.Response{
		Code:    401,
		Message: msg,
	})
}

// ReadReusableBody 读取请求体并回写，便于后续处理链继续消费 Body。
func ReadReusableBody(req *http.Request) ([]byte, error) {
	if req == nil || req.Body == nil {
		return nil, nil
	}

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes, nil
}

// BuildRequestURL 组装请求路径，包含查询参数。
func BuildRequestURL(req *http.Request) string {
	if req == nil || req.URL == nil {
		return ""
	}
	if req.URL.RawQuery == "" {
		return req.URL.Path
	}
	return req.URL.Path + "?" + req.URL.RawQuery
}

// ClientIPFromRequest 获取真实客户端 IP，优先处理反向代理场景。
func ClientIPFromRequest(req *http.Request) string {
	if req == nil {
		return ""
	}

	xForwardedFor := req.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			if ip != "" && !IsInternalIP(ip) {
				return ip
			}
		}
	}

	xRealIP := req.Header.Get("X-Real-IP")
	if xRealIP != "" && !IsInternalIP(xRealIP) {
		return xRealIP
	}

	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return req.RemoteAddr
	}
	return ip
}

// IsInternalIP 判断 IP 是否为本机或内网地址。
func IsInternalIP(ip string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	return parsed.IsPrivate() || parsed.Equal(net.IPv4(127, 0, 0, 1))
}

// TruncateText 按最大长度截断文本，并在必要时追加省略号。
func TruncateText(text string, max int) string {
	if max <= 0 || len(text) <= max {
		return text
	}
	if max <= 3 {
		return text[:max]
	}
	return text[:max-3] + "..."
}

// CompactText 去除换行与制表符，便于日志存储。
func CompactText(text string) string {
	return strings.ReplaceAll(strings.ReplaceAll(text, "\n", ""), "\t", "")
}

// CaptureRequestPayload 读取并压缩请求体文本，适合写入日志。
func CaptureRequestPayload(req *http.Request, max int) string {
	bodyBytes, err := ReadReusableBody(req)
	if err != nil || len(bodyBytes) == 0 {
		return ""
	}
	return TruncateText(CompactText(string(bodyBytes)), max)
}

// ExtractResponseError 从响应体和 gin 错误栈中提取用于日志记录的错误信息。
func ExtractResponseError(statusCode int, c *gin.Context, body []byte, max int) string {
	if statusCode < http.StatusBadRequest {
		return ""
	}

	var respMap map[string]any
	if err := json.Unmarshal(body, &respMap); err == nil {
		for _, key := range []string{"msg", "message", "error"} {
			if msg, ok := respMap[key].(string); ok && msg != "" {
				return TruncateText(msg, max)
			}
		}
	}

	if errText := strings.TrimSpace(c.Errors.String()); errText != "" {
		return TruncateText(errText, max)
	}

	return TruncateText(http.StatusText(statusCode), max)
}
