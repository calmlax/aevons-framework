package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/calmlax/aevons-framework/auth"
	"github.com/calmlax/aevons-framework/consts"
	"github.com/calmlax/aevons-framework/response"
	"github.com/calmlax/aevons-framework/utils"

	"github.com/gin-gonic/gin"
)

// PermissionMap 以集合形式存储权限标识，便于 O(1) 查找。
type PermissionMap map[string]struct{}

// isExcluded 判断请求路径是否命中排除规则列表。
// 支持两种格式：
//   - 精确全路径："/api/v1/auth/login"（完全相等）
//   - 路径段通配符："/api/*/auth/login"（* 匹配单个路径段，如 v1、v2）
//   - 前缀通配符："/api/v1/public/*"（末尾 * 匹配该前缀下的所有子路径）
func isExcluded(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if matchPath(path, pattern) {
			return true
		}
	}
	return false
}

// matchPath 将请求路径与单条规则进行匹配。
// * 作为路径段通配符，仅匹配一段（不含 /）；末尾 * 匹配所有子路径。
func matchPath(path, pattern string) bool {
	// 末尾通配符：/api/v1/public/* 匹配 /api/v1/public/ 开头的所有路径
	if strings.HasSuffix(pattern, "/*") {
		prefix := pattern[:len(pattern)-1] // 保留末尾的 /
		return strings.HasPrefix(path, prefix)
	}

	// 无通配符：精确全路径匹配
	if !strings.Contains(pattern, "*") {
		return path == pattern
	}

	// 路径段通配符：逐段比较
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")
	if len(patternParts) != len(pathParts) {
		return false
	}
	for i, p := range patternParts {
		if p == "*" {
			continue // 通配符匹配任意单段
		}
		if p != pathParts[i] {
			return false
		}
	}
	return true
}

// SetPermission 将权限标识列表以 PermissionMap 形式存入 gin.Context。
func SetPermission(c *gin.Context, rawPermissions ...string) PermissionMap {
	permissionMap := make(PermissionMap, len(rawPermissions))
	for _, p := range rawPermissions {
		permissionMap[p] = struct{}{}
	}
	c.Set(consts.UserPermissionKey, permissionMap)
	return permissionMap
}

// AuthMiddleware 对每个请求进行 Bearer Token 验证。
// excludes 来自配置文件 auth.excludes，命中规则的路径直接放行。
func AuthMiddleware(store auth.TokenStore, excludes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 排除路径直接放行
		if isExcluded(c.Request.URL.Path, excludes) {
			c.Next()
			return
		}

		// 从 Authorization 请求头提取 Bearer Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Fail(c, http.StatusUnauthorized, http.StatusUnauthorized, consts.ErrTokenMissing)
			c.Abort()
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			response.Fail(c, http.StatusUnauthorized, http.StatusUnauthorized, consts.ErrTokenMissing)
			c.Abort()
			return
		}

		// 从 Redis 加载 LoginUser
		loginUser, err := store.GetLoginUser(c.Request.Context(), token)
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, http.StatusUnauthorized, consts.ErrTokenExpired)
			c.Abort()
			return
		}

		// 将用户数据注入 Context
		c.Set(consts.LoginUserKey, loginUser)
		c.Set(consts.UserIdKey, loginUser.UserId)
		c.Set(consts.UserRoleKey, loginUser.Roles)
		c.Set(consts.UserDeptKey, loginUser.Depts)
		SetPermission(c, loginUser.Permissions...)

		// 同步注入到标准 Request.Context，供 Service 层直接传递和解析使用
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, consts.LoginUserKey, loginUser)
		ctx = context.WithValue(ctx, consts.UserIdKey, loginUser.UserId)
		ctx = context.WithValue(ctx, consts.UserRoleKey, loginUser.Roles)
		ctx = context.WithValue(ctx, consts.UserDeptKey, loginUser.Depts)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// HasPermission 校验当前用户是否持有指定权限标识。
func HasPermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if utils.IsEmpty(permission) {
			c.Next()
			return
		}
		val, exists := c.Get(consts.UserPermissionKey)
		if !exists {
			response.Fail(c, http.StatusUnauthorized, http.StatusUnauthorized, consts.ErrUnauthorized)
			c.Abort()
			return
		}
		permissionMap := val.(PermissionMap)

		_, hasAll := permissionMap[consts.AllPermission]
		_, hasSpecific := permissionMap[permission]

		if !hasAll && !hasSpecific {
			response.Fail(c, http.StatusForbidden, http.StatusForbidden, consts.ErrPermissionDenied, map[string]any{"permission": permission})
			c.Abort()
			return
		}
		c.Next()
	}
}

// HasAnyPermission 校验当前用户是否持有给定权限列表中的至少一个。
func HasAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(permissions) == 0 {
			c.Next()
			return
		}
		val, exists := c.Get(consts.UserPermissionKey)
		if !exists {
			response.Fail(c, http.StatusUnauthorized, http.StatusUnauthorized, consts.ErrUnauthorized)
			c.Abort()
			return
		}
		permissionMap := val.(PermissionMap)

		_, hasAll := permissionMap[consts.AllPermission]
		if hasAll {
			c.Next()
			return
		}

		for _, perm := range permissions {
			if _, has := permissionMap[perm]; has {
				c.Next()
				return
			}
		}
		response.Fail(c, http.StatusForbidden, http.StatusForbidden, consts.ErrPermissionDenied, map[string]any{"permissions": permissions})
		c.Abort()
	}
}

// HasRole 校验当前用户是否持有指定角色。
func HasRole(roleKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if utils.IsEmpty(roleKey) {
			c.Next()
			return
		}
		val, exists := c.Get(consts.UserRoleKey)
		if !exists {
			response.Fail(c, http.StatusUnauthorized, http.StatusUnauthorized, consts.ErrUnauthorized)
			c.Abort()
			return
		}
		keys, ok := val.([]string)
		if !ok {
			response.Fail(c, http.StatusForbidden, http.StatusForbidden, consts.ErrRoleInvalid)
			c.Abort()
			return
		}

		for _, key := range keys {
			if key == consts.SuperAdminRoleKey || key == roleKey {
				c.Next()
				return
			}
		}
		response.Fail(c, http.StatusForbidden, http.StatusForbidden, consts.ErrRoleDenied, map[string]any{"role": roleKey})
		c.Abort()
	}
}

// HasAnyRole 校验当前用户是否持有给定角色列表中的至少一个。
func HasAnyRole(roleKeys ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(roleKeys) == 0 {
			c.Next()
			return
		}
		val, exists := c.Get(consts.UserRoleKey)
		if !exists {
			response.Fail(c, http.StatusUnauthorized, http.StatusUnauthorized, consts.ErrUnauthorized)
			c.Abort()
			return
		}
		keys, ok := val.([]string)
		if !ok {
			response.Fail(c, http.StatusForbidden, http.StatusForbidden, consts.ErrRoleInvalid)
			c.Abort()
			return
		}

		requiredRoleMap := make(map[string]bool, len(roleKeys))
		for _, r := range roleKeys {
			requiredRoleMap[r] = true
		}

		for _, roleKey := range keys {
			if roleKey == consts.SuperAdminRoleKey || requiredRoleMap[roleKey] {
				c.Next()
				return
			}
		}

		response.Fail(c, http.StatusForbidden, http.StatusForbidden, consts.ErrRoleDenied, map[string]any{"roles": roleKeys})
		c.Abort()
	}
}
