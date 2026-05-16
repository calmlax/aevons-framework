package base

import (
	"reflect"
	"strconv"

	apperr "github.com/calmlax/aevons-framework/errors"
	"github.com/calmlax/aevons-framework/response"
	"github.com/calmlax/aevons-framework/utils"

	"github.com/gin-gonic/gin"
)

// BaseHandler 是一个轻量级的泛型控制器辅助器。
// 它不直接作为业务 handler 被继承，而是由各业务 handler 在构造阶段注入对应的
// Model、Query、CreateDTO、UpdateDTO 和 Service，
// 以便复用标准的列表、分页、详情、创建、更新、删除等公共处理流程。
type BaseHandler[M any, Q BaseQueryInterface, C any, U any] struct {
	srv BaseService[M, Q]
}

// NewBaseHandler 创建一个绑定了具体模型与服务的通用控制器辅助器。
func NewBaseHandler[M any, Q BaseQueryInterface, C any, U any](srv BaseService[M, Q]) *BaseHandler[M, Q, C, U] {
	return &BaseHandler[M, Q, C, U]{srv: srv}
}

// BindQuery 将查询参数绑定到指定的 Query DTO。
// Q 通常是实现了 BaseQueryInterface 的查询结构体指针类型。
func BindQuery[Q BaseQueryInterface](c *gin.Context) (Q, bool) {
	// 获取接口的真实类型
	qType := reflect.TypeOf((*Q)(nil)).Elem()
	if qType.Kind() == reflect.Ptr {
		qType = qType.Elem()
	}

	// 创建实例
	qVal := reflect.New(qType)

	// 绑定参数
	if err := c.ShouldBindQuery(qVal.Interface()); err != nil {
		response.FailBy(c, apperr.ErrInvalidQuery)
		return *new(Q), false
	}

	// 转换为接口 Q 并返回
	return qVal.Interface().(Q), true
}

// HandleList 处理标准的不分页列表查询。
func (h *BaseHandler[M, Q, C, U]) HandleList(c *gin.Context) {
	q, ok := BindQuery[Q](c)
	if !ok {
		return
	}
	list, err := h.srv.List(q)
	if err != nil {
		response.FailBy(c, apperr.ErrQueryFailed)
		return
	}
	response.Success(c, list)
}

// HandlePage 处理标准的分页查询。
func (h *BaseHandler[M, Q, C, U]) HandlePage(c *gin.Context) {
	q, ok := BindQuery[Q](c)
	if !ok {
		return
	}
	res, err := h.srv.Page(q)
	if err != nil {
		response.FailBy(c, apperr.ErrQueryFailed)
		return
	}
	response.Success(c, res)
}

// HandleGet 按路径参数 id 查询单条记录详情。
func (h *BaseHandler[M, Q, C, U]) HandleGet(c *gin.Context) {
	id, ok := GetId(c)
	if !ok {
		return
	}
	info, err := h.srv.GetById(id)
	if err != nil {
		response.FailBy(c, apperr.ErrQueryFailed)
		return
	}
	response.Success(c, info)
}

// HandleCreate 处理标准创建逻辑。
// 它会将 CreateDTO 绑定到请求体后复制到模型，再调用对应的 service.Create。
func (h *BaseHandler[M, Q, C, U]) HandleCreate(c *gin.Context) {
	var dto C
	if !BindJSON(c, &dto) {
		return
	}
	var m M
	utils.Copy(&m, dto)
	if err := h.srv.Create(&m); err != nil {
		response.FailBy(c, apperr.ErrCreateFailed)
		return
	}
	response.Success(c, m)
}

// HandleUpdate 处理标准更新逻辑。
// 它会将 UpdateDTO 转换为 map，并调用对应的 service.Update。
func (h *BaseHandler[M, Q, C, U]) HandleUpdate(c *gin.Context) {
	id, ok := GetId(c)
	if !ok {
		return
	}
	var dto U
	if !BindJSON(c, &dto) {
		return
	}
	mp := utils.StructToMapIgnoreNil(dto)
	if _, err := h.srv.Update(id, mp); err != nil {
		response.FailBy(c, apperr.ErrUpdateFailed)
		return
	}
	response.Success(c, id)
}

// HandleDelete 按路径参数 id 删除单条记录。
func (h *BaseHandler[M, Q, C, U]) HandleDelete(c *gin.Context) {
	id, ok := GetId(c)
	if !ok {
		return
	}
	if err := h.srv.Delete(id); err != nil {
		response.FailBy(c, apperr.ErrDeleteFailed)
		return
	}
	response.Success(c, id)
}

// HandleBatchDelete 按路径参数 ids 批量删除记录。
// ids 约定使用逗号分隔，例如：1,2,3
func (h *BaseHandler[M, Q, C, U]) HandleBatchDelete(c *gin.Context) {
	ids, ok := GetIds(c)
	if !ok {
		return
	}
	if err := h.srv.BatchDelete(ids); err != nil {
		response.FailBy(c, apperr.ErrDeleteFailed)
		return
	}
	response.Success(c, ids)
}

// BindJSON 绑定 JSON 请求体，失败时返回统一错误响应。
func BindJSON(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return false
	}
	return true
}

// GetId 从路径参数 id 中解析单个 int64 主键。
func GetId(c *gin.Context) (int64, bool) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailBy(c, apperr.ErrInvalidId)
		return 0, false
	}
	return id, true
}

// GetIds 从路径参数 ids 中解析批量主键。
// ids 约定使用逗号分隔，例如：1,2,3
func GetIds(c *gin.Context) ([]int64, bool) {
	idsStr := c.Param("ids")
	ids, err := utils.StrToNumberArray[int64](idsStr, ",")
	if err != nil {
		response.FailBy(c, apperr.ErrInvalidId)
		return nil, false
	}
	return ids, true
}
