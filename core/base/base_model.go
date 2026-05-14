package base

import "time"

type BaseModel struct {
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp" json:"updatedAt,omitempty"`
}

type DefaultModel struct {
	CreatedBy int64 `gorm:"column:created_by" json:"createdBy,omitempty,string"`
	UpdatedBy int64 `gorm:"column:updated_by" json:"updatedBy,omitempty,string"`
	BaseModel
}

type LogicDeleteModel struct {
	IsDeleted bool `gorm:"column:is_deleted;type:tinyint(1);default:0" json:"-"`
}

type CommonModel struct {
	LogicDeleteModel
	DefaultModel
}
