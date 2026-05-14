package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/calmlax/aevons-framework/config"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

// XSSMiddleware 防止 XSS 攻击的全局拦截中间件
func XSSMiddleware(cfg *config.Config) gin.HandlerFunc {
	// 使用严格策略，剥离所有 HTML 标签，仅保留纯文本
	policy := bluemonday.StrictPolicy()

	return func(c *gin.Context) {
		// 1. 如果未开启，直接放行
		if !cfg.XSS.Enabled {
			c.Next()
			return
		}

		// 2. 检查是否在排除名单中（如富文本接口）
		if xssIsExcluded(c.Request.URL.Path, cfg.XSS.Excludes) {
			c.Next()
			return
		}

		// 3. 处理请求体 JSON 数据
		contentType := c.Request.Header.Get("Content-Type")
		if c.Request.Body != nil && strings.Contains(contentType, "application/json") {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil && len(bodyBytes) > 0 {
				var data any
				// 尝试把 JSON 解码为弱类型的 map 结构
				if err := json.Unmarshal(bodyBytes, &data); err == nil {
					// 递归清理所有文本字段中的恶意注入
					sanitizedData := sanitizeData(policy, data)
					newBodyBytes, _ := json.Marshal(sanitizedData)
					// 回写覆盖原始 Body，给后续绑定使用
					c.Request.Body = io.NopCloser(bytes.NewBuffer(newBodyBytes))
					c.Request.ContentLength = int64(len(newBodyBytes))
				} else {
					// 格式不标准导致解析失败时，原样复原，后续框架的 Binding 也会自行抛错拦截
					c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				}
			}
		}

		// 4. 处理 Query/Form URL 参数
		query := c.Request.URL.Query()
		modified := false
		for key, values := range query {
			for i, val := range values {
				sanitized := policy.Sanitize(val)
				if sanitized != val {
					query[key][i] = sanitized
					modified = true
				}
			}
		}
		if modified {
			// 如果有被转义清理的值，要重写 URL 以便 c.Query() 获取清理后的值
			c.Request.URL.RawQuery = query.Encode()
		}

		c.Next()
	}
}

// sanitizeData 递归清理数据中的字符串
func sanitizeData(policy *bluemonday.Policy, data any) any {
	switch v := data.(type) {
	case string:
		// 最底层核心机制，过滤 HTML
		return policy.Sanitize(v)
	case map[string]any:
		for key, val := range v {
			v[key] = sanitizeData(policy, val)
		}
		return v
	case []any:
		for i, val := range v {
			v[i] = sanitizeData(policy, val)
		}
		return v
	default:
		// 数字、布尔等其它类型直接放行
		return data
	}
}

// xssIsExcluded 判断请求路径是否命中免过滤规则列表。
func xssIsExcluded(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if xssMatchPath(path, pattern) {
			return true
		}
	}
	return false
}

// xssMatchPath 通用的路径规则单点匹配器
func xssMatchPath(path, pattern string) bool {
	// 后缀通配符：/api/v1/content/*
	if strings.HasSuffix(pattern, "/*") {
		prefix := pattern[:len(pattern)-1]
		return strings.HasPrefix(path, prefix)
	}

	// 无通配符，精确比较
	if !strings.Contains(pattern, "*") {
		return path == pattern
	}

	// 段内通配符匹配
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")
	if len(patternParts) != len(pathParts) {
		return false
	}
	for i, p := range patternParts {
		if p == "*" {
			continue
		}
		if p != pathParts[i] {
			return false
		}
	}
	return true
}
