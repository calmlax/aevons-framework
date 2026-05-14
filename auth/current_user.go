package auth

import (
	"context"
	"errors"

	"github.com/calmlax/aevons-framework/consts"
)

var (
	ErrNoLoginUser = errors.New("上下文中不存在登录用户")
	ErrNoUserId    = errors.New("上下文中不存在 user_id")
	ErrNoRoles     = errors.New("上下文中不存在角色列表")
	ErrNoPerms     = errors.New("上下文中不存在权限列表")
)

// GetCurrentUser 从 context.Context 中提取完整的 LoginUser 对象。
func GetCurrentUser(ctx context.Context) (*LoginUser, error) {
	val := ctx.Value(consts.LoginUserKey)
	if val == nil {
		return nil, ErrNoLoginUser
	}
	user, ok := val.(*LoginUser)
	if !ok || user == nil {
		return nil, ErrNoLoginUser
	}
	user.RefreshToken = ""
	return user, nil
}

// GetCurrentUserId 从 context.Context 中提取当前用户 Id。
func GetCurrentUserId(ctx context.Context) (int64, error) {
	val := ctx.Value(consts.UserIdKey)
	if val == nil {
		return 0, ErrNoUserId
	}
	id, ok := val.(int64)
	if !ok {
		return 0, ErrNoUserId
	}
	return id, nil
}

// GetCurrentRoleIds 从 context.Context 中提取当前用户的角色列表。
func GetCurrentRoleIds(ctx context.Context) ([]int64, error) {
	val := ctx.Value(consts.UserRoleKey)
	if val == nil {
		return nil, ErrNoRoles
	}
	roles, ok := val.([]Role)
	if !ok {
		return nil, ErrNoRoles
	}
	ids := make([]int64, len(roles))
	for i, role := range roles {
		ids[i] = role.Id
	}
	return ids, nil
}

// GetCurrentRoleKeys 从 context.Context 中提取当前用户的角色列表。
func GetCurrentRoleKeys(ctx context.Context) ([]string, error) {
	val := ctx.Value(consts.UserRoleKey)
	if val == nil {
		return nil, ErrNoRoles
	}
	roles, ok := val.([]Role)
	if !ok {
		return nil, ErrNoRoles
	}
	keys := make([]string, len(roles))
	for i, role := range roles {
		keys[i] = role.RoleKey
	}
	return keys, nil
}

// GetCurrentDepts 从 context.Context 中提取当前用户的部门列表。
func GetCurrentDepts(ctx context.Context) ([]Dept, error) {
	val := ctx.Value(consts.UserDeptKey)
	if val == nil {
		return nil, ErrNoRoles
	}
	depts, ok := val.([]Dept)
	if !ok {
		return nil, ErrNoRoles
	}
	return depts, nil
}

// GetCurrentDeptIds 从 context.Context 中提取当前用户的部门ID列表。
func GetCurrentDeptIds(ctx context.Context) ([]int64, error) {
	val := ctx.Value(consts.UserDeptKey)
	if val == nil {
		return nil, ErrNoRoles
	}
	depts, ok := val.([]Dept)
	if !ok {
		return nil, ErrNoRoles
	}
	ids := make([]int64, len(depts))
	for i, dept := range depts {
		ids[i] = dept.DeptId
	}
	return ids, nil
}

// GetCurrentPermissions 从 context.Context 中提取当前用户的权限列表。
// 权限可能以 []string 或 PermissionMap（map[string]struct{}）形式存储。
func GetCurrentPermissions(ctx context.Context) ([]string, error) {
	val := ctx.Value(consts.UserPermissionKey)
	if val == nil {
		return nil, ErrNoPerms
	}
	switch v := val.(type) {
	case []string:
		return v, nil
	case map[string]struct{}:
		perms := make([]string, 0, len(v))
		for p := range v {
			perms = append(perms, p)
		}
		return perms, nil
	default:
		return nil, ErrNoPerms
	}
}
