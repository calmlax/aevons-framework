package base

import (
	"fmt"

	"github.com/calmlax/aevons-framework/utils"
)

// BaseQuery 分页+排序基础查询参数
// 所有列表查询接口可直接嵌套使用，统一接收分页、排序参数
type BaseQuery struct {
	// 排序方式 ascend(升序)/descend(降序) 非必填，传值仅限枚举
	Direction string `form:"direction" binding:"omitempty,oneof=asc desc ascend descend"`
	// 排序字段 非必填
	Field string `form:"field" binding:"omitempty"`
	// 当前页码 非必填，最小值为1
	PageNum int `form:"pageNum" binding:"omitempty,min=1"`
	// 每页条数 非必填，范围1~1000
	PageSize int `form:"pageSize" binding:"omitempty,min=1,max=1000"`
}

// GetOrder 获取排序字段与是否降序
// 返回：排序字段、是否降序(true=降序/false=升序)、错误信息
func (b *BaseQuery) GetOrder() (field string, isDesc bool, err error) {
	// 设置默认排序方式
	direction := b.Direction
	if direction == "" {
		direction = "desc"
	}

	// 校验排序字段安全性（防SQL注入）
	field = b.Field
	if !utils.IsSafeField(field) {
		return "", true, fmt.Errorf("无效的排序字段")
	}

	// 返回是否为降序
	return utils.ToSnake(field), direction == "desc" || direction == "descend", nil
}

// Normalize 设置分页默认值
func (p *BaseQuery) Normalize() {
	if p.PageNum <= 0 {
		p.PageNum = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
}

// Offset 计算数据库偏移量
func (p *BaseQuery) Offset() int {
	return (p.PageNum - 1) * p.PageSize
}

// GetPageSize 获取每页条数
func (p *BaseQuery) GetPageSize() int {
	return p.PageSize
}
