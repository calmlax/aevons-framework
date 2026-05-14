package utils

import (
	"bytes"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
)

func ToSet(items []string) map[string]bool {
	m := make(map[string]bool, len(items))
	for _, item := range items {
		m[item] = true
	}
	return m
}

// removePrefixByPatterns 从 exclude_tables 的纯前缀模式（不含 %）中去除表名前缀
func RemovePrefixByPatterns(name string, patterns []string) string {
	for _, p := range patterns {
		// 只处理纯前缀（不含通配符），如 "easy_"、"xc_"
		if !strings.Contains(p, "%") && strings.HasPrefix(name, p) {
			return name[len(p):]
		}
	}
	return name
}

// MysqlOrPgTypeToDataType
// 自动支持 MySQL + PostgreSQL
// 无需区分库，自动识别 unsigned，自动适配 PG 无符号特性
func MysqlOrPgTypeToDataType(dbType string) string {
	t := strings.ToLower(strings.TrimSpace(dbType))
	isUnsigned := strings.Contains(t, "unsigned")

	// 布尔类型
	if strings.HasPrefix(t, "bool") ||
		strings.HasPrefix(t, "boolean") {
		return "bool"
	}

	// 整数类型
	switch {
	case strings.HasPrefix(t, "smallint"), strings.HasPrefix(t, "tinyint"):
		if isUnsigned {
			return "uint16"
		}
		return "int16"

	case strings.HasPrefix(t, "int"), strings.HasPrefix(t, "integer"):
		if isUnsigned {
			return "uint"
		}
		return "int"

	case strings.HasPrefix(t, "bigint"):
		if isUnsigned {
			return "uint64"
		}
		return "int64"
	}

	// 浮点
	if strings.HasPrefix(t, "float") && !strings.HasPrefix(t, "float8") {
		return "float32"
	}
	if strings.HasPrefix(t, "double") || strings.HasPrefix(t, "float8") {
		return "float64"
	}

	// 高精度小数
	if strings.HasPrefix(t, "decimal") || strings.HasPrefix(t, "numeric") {
		return "decimal.Decimal"
	}

	// 字符串
	if strings.Contains(t, "char") ||
		strings.Contains(t, "text") ||
		strings.HasPrefix(t, "enum") {
		return "string"
	}

	// 时间
	if strings.Contains(t, "time") ||
		strings.HasPrefix(t, "date") ||
		strings.HasPrefix(t, "timestamp") {
		return "time.Time"
	}

	// JSON
	if strings.HasPrefix(t, "json") {
		return "json.RawMessage"
	}

	// 字节
	if strings.HasPrefix(t, "byte") || strings.HasPrefix(t, "blob") {
		return "[]byte"
	}

	return "string"
}

// ParseTemplateFile 从 templates/ 目录读取模板文件，渲染后返回字符串
func ParseTemplateFile(tplFilename string, data any) (string, error) {
	tplPath := filepath.Join("templates", tplFilename)

	tpl, err := template.New(filepath.Base(tplPath)).Funcs(template.FuncMap{
		"toLowerCamel": func(s string) string {
			return ToLowerCamel(s)
		},
		"toUpperCamel": func(s string) string {
			return ToUpperCamel(s)
		},
		"contains": func(s []string, v string) bool {
			return slices.Contains(s, v)
		},
	}).ParseFiles(tplPath)

	if err != nil {
		return "", err
	}

	tpl.Option("missingkey=error")

	// 3. 渲染模板
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// 渲染：传入模板字符串 + 数据 → 返回渲染后的字符串
func ParseTemplateString(tpl string, data any) (string, error) {
	t, err := template.New("template").Parse(tpl)
	if err != nil {
		return "", err
	}

	// 渲染到 buffer
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
