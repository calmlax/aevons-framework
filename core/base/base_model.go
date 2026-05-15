package base

import "time"

// BaseModel 定义所有实体共享的时间字段。
// CreatedAt / UpdatedAt 会由 GORM 按约定自动维护。
type BaseModel struct {
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp" json:"updatedAt,omitempty"`
}

// DefaultModel 定义常见的审计字段。
// 当模型包含这些字段时，框架会在 db.Init() 注册的全局回调中自动填充：
// 1. 创建时自动补 CreatedBy / UpdatedBy
// 2. 更新时自动刷新 UpdatedBy
// 3. 当前上下文没有登录用户时，默认写入 0
type DefaultModel struct {
	CreatedBy int64 `gorm:"column:created_by" json:"createdBy,omitempty,string"`
	UpdatedBy int64 `gorm:"column:updated_by" json:"updatedBy,omitempty,string"`
	BaseModel
}

// LogicDeleteModel 定义逻辑删除标记。
// 创建数据时如果模型包含该字段，框架会自动保持为 false（数据库中即 0）。
type LogicDeleteModel struct {
	IsDeleted bool `gorm:"column:is_deleted;type:tinyint(1);default:0" json:"-"`
}

// CommonModel 聚合逻辑删除与审计字段，适合作为大多数业务模型的通用基类。
type CommonModel struct {
	LogicDeleteModel
	DefaultModel
}
