// Package utils 通用字符串工具类
package utils

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

// IsEmpty 判断字符串是否为空（不忽略空白）
func IsEmpty(v any) bool {
	if v == nil {
		return true
	}

	val := reflect.ValueOf(v)

	switch val.Kind() {

	case reflect.String:
		return val.Len() == 0

	case reflect.Array, reflect.Slice, reflect.Map:
		return val.Len() == 0

	case reflect.Ptr, reflect.Interface:
		return val.IsNil()

	case reflect.Struct:
		// struct 默认不认为空（更安全）
		return false
	}

	return false
}

// IsNotEmpty 判断字符串是否非空（不忽略空白）
func IsNotEmpty(v any) bool {
	return !IsEmpty(v)
}

// DefaultIfEmpty 若字符串为空则返回默认值
func DefaultIfEmpty(s, def string) string {
	if IsEmpty(s) {
		return def
	}
	return s
}

// ContainsIgnoreCase 判断是否包含子串（忽略大小写）
func ContainsIgnoreCase(s, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
}

// Truncate 截断字符串到指定长度，超出部分用 ellipsis 替代，支持 Unicode
func Truncate(s string, maxLen int, ellipsis string) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen <= 0 {
		return ellipsis
	}
	return string(runes[:maxLen]) + ellipsis
}

// RemoveSpaces 移除字符串中所有空白字符
func RemoveSpaces(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}

// ReverseStr 反转字符串，支持 Unicode
func ReverseStr(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ToSnake 驼峰转下划线
func ToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ToUpperCamel 下划线转大驼峰
func ToUpperCamel(s string) string {
	if s == "" {
		return s
	}

	// 统一分隔符（可选增强）
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, " ", "_")

	// 如果包含下划线 → 下划线转驼峰
	if strings.Contains(s, "_") {
		s = strings.ToLower(s)
		parts := strings.Split(s, "_")

		for i := 0; i < len(parts); i++ {
			if len(parts[i]) > 0 {
				parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
			}
		}
		return strings.Join(parts, "")
	}

	// 已经是驼峰或普通字符串 → 只处理首字母
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// ToLowerCamel 下划线转小驼峰
func ToLowerCamel(s string) string {
	if s == "" {
		return s
	}

	// 如果包含下划线 → 走下划线逻辑
	if strings.Contains(s, "_") {
		s = strings.ToLower(s)
		parts := strings.Split(s, "_")

		for i := 1; i < len(parts); i++ {
			if len(parts[i]) > 0 {
				parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
			}
		}
		return strings.Join(parts, "")
	}

	// 如果是驼峰（或普通字符串）→ 只把首字母小写
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// SubStr 安全截取字符串，支持中文
func SubStr(s string, start, length int) string {
	runes := []rune(s)
	total := len(runes)
	if start >= total {
		return ""
	}
	end := start + length
	if end > total {
		end = total
	}
	return string(runes[start:end])
}

// MaskPhone 手机号脱敏 138****1234
func MaskPhone(phone string) string {
	if len(phone) != 11 {
		return phone
	}
	return phone[:3] + "****" + phone[7:]
}

// MaskEmail 邮箱脱敏 t***@gmail.com
func MaskEmail(email string) string {
	idx := strings.Index(email, "@")
	if idx == -1 || idx <= 1 {
		return email
	}
	return email[:1] + "***" + email[idx:]
}

// IsDigit 是否纯数字
func IsDigit(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return s != ""
}

// IsLetter 是否纯字母
func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return s != ""
}

// TrimChar 去除首尾指定字符
func TrimChar(s string, c rune) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return r == c
	})
}

// StrToNumberArray 泛型工具：字符串分割 → 转成任意数字类型数组
// 支持：int, int64, uint, uint64 等
// sep 为空时默认使用逗号 , 分割
func StrToNumberArray[T int | int64 | uint | uint64](idsStr string, sep string) ([]T, error) {
	if idsStr == "" {
		return []T{}, nil
	}

	// 默认分隔符 ,
	if sep == "" {
		sep = ","
	}

	strArr := strings.Split(idsStr, sep)
	var res []T

	for _, s := range strArr {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		// 解析成 int64 再强转泛型（最稳）
		num, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}

		res = append(res, T(num))
	}

	return res, nil
}

// Split 逗号分隔字符串转切片，自动去空格、排除空值
// 支持：table1,table2  | table1 , table2  | table1,,table2 全部自动处理
func Split(s string, sep string) []string {
	if s == "" {
		return []string{}
	}

	var result []string
	// 使用传入的分隔符 sep，不是写死的 ,
	for _, item := range strings.Split(s, sep) {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

// Contains 判断字符串切片是否包含指定元素
func Contains(list []string, item string) bool {
	for _, s := range list {
		if s == item {
			return true
		}
	}
	return false
}

// ContainsAllMap 判断 target 里的所有元素是否都在 list 中
func ContainsAllMap(data map[string]bool, target []string) bool {
	// 遍历所有目标，必须全部存在
	for _, s := range target {
		if !data[s] {
			return false
		}
	}
	return true
}

// ContainsAll 判断 target 里的所有元素是否都在 list 中
func ContainsAll(list []string, target []string) bool {
	// 先把 list 转成 map 加速查询
	checkMap := make(map[string]bool)
	for _, s := range list {
		checkMap[s] = true
	}
	return ContainsAllMap(checkMap, target)
}

// Copy 结构体字段拷贝（支持 DTO → Model，同名字段自动复制）
func Copy(dst, src interface{}) {
	srcVal := getReflectValue(src)
	dstVal := getReflectValue(dst)

	// 必须是结构体才能拷贝字段
	if srcVal.Kind() != reflect.Struct || dstVal.Kind() != reflect.Struct {
		return
	}

	// 遍历源结构体的字段
	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Type().Field(i)
		srcValue := srcVal.Field(i)

		// 在目标结构体中找同名字段
		dstField := dstVal.FieldByName(srcField.Name)
		if !dstField.IsValid() || !dstField.CanSet() {
			continue
		}

		// 类型相同才赋值
		if srcValue.Type() == dstField.Type() {
			dstField.Set(srcValue)
		}
	}
}

// getReflectValue 获取反射值（自动解指针）
func getReflectValue(v interface{}) reflect.Value {
	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		val = val.Elem()
	}
	return val
}

// StructToMap 结构体转 map（用于更新）
func StructToMap(obj any) map[string]any {
	data := make(map[string]any)
	bytes, _ := json.Marshal(obj)
	_ = json.Unmarshal(bytes, &data)
	return data
}
func StructToMapIgnoreNil(obj any) map[string]any {
	if obj == nil {
		return nil
	}

	val := reflect.ValueOf(obj)

	// 解引用指针
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	typ := val.Type()
	result := make(map[string]any, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		// 跳过未导出字段（避免 panic）
		if fieldType.PkgPath != "" {
			continue
		}

		tag := parseJSONTag(fieldType)

		// 忽略 "-"
		if tag == "-" {
			continue
		}

		// 👉 处理指针
		if fieldVal.Kind() == reflect.Ptr {
			if fieldVal.IsNil() {
				continue
			}
			result[tag] = fieldVal.Elem().Interface()
			continue
		}

		// 👉 处理 struct（递归）
		if fieldVal.Kind() == reflect.Struct {
			result[tag] = StructToMapIgnoreNil(fieldVal.Interface())
			continue
		}

		// 👉 处理 slice/map（可选：忽略空）
		if fieldVal.Kind() == reflect.Slice || fieldVal.Kind() == reflect.Map {
			if fieldVal.Len() == 0 {
				continue
			}
		}

		// 👉 默认值直接放
		result[tag] = fieldVal.Interface()
	}

	return result
}

func parseJSONTag(f reflect.StructField) string {
	tag := f.Tag.Get("json")

	if tag == "" {
		return f.Name
	}

	name := strings.Split(tag, ",")[0]
	if name == "" {
		return f.Name
	}

	return name
}

// GetFileSuffix 获取文件后缀，不带点，例如：sql、go、vue
func GetFileSuffix(fileName string) string {
	ext := filepath.Ext(fileName)
	if len(ext) > 1 {
		return ext[1:] // 去掉前面的 .
	}
	return ""
}

// PtrVal 安全解引用指针，避免 nil pointer panic
func PtrVal[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

func IsSafeField(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if !(r == '_' ||
			(r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}
