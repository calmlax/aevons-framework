package base

import (
	"reflect"
	"strconv"

	apperr "github.com/calmlax/aevons-framework/errors"
	"github.com/calmlax/aevons-framework/response"
	"github.com/calmlax/aevons-framework/utils"

	"github.com/gin-gonic/gin"
)

// 定义两个接口，用来约束 DTO 必须有 ToMap 方法

// 泛型基础控制器
// M = Model 模型
// Q = Query 查询条件
// C = CreateDTO
// U = UpdateDTO
type BaseHandler[M any, Q BaseQueryInterface, C any, U any] struct {
	BaseHandlerCommon
	srv BaseService[M, Q]
}

func NewBaseHandler[M any, Q BaseQueryInterface, C any, U any](
	srv BaseService[M, Q],
) *BaseHandler[M, Q, C, U] {
	return &BaseHandler[M, Q, C, U]{srv: srv}
}

// ==================== 通用 CURD 接口 ====================

// List 列表（不分页）
func (h *BaseHandler[M, Q, C, U]) List(c *gin.Context) {
	q, _ := bindQuery[Q](c)
	list, err := h.srv.List(q)
	if err != nil {
		h.Fail(c, apperr.ErrQueryFailed)
		return
	}
	h.Success(c, list)
}

// Page 分页
func (h *BaseHandler[M, Q, C, U]) Page(c *gin.Context) {
	q, _ := bindQuery[Q](c)
	res, err := h.srv.Page(q)
	if err != nil {
		h.Fail(c, apperr.ErrQueryFailed)
		return
	}
	h.Success(c, res)
}

// Get 单条
func (h *BaseHandler[M, Q, C, U]) Get(c *gin.Context) {
	id, ok := h.GetId(c)
	if !ok {
		return
	}

	info, err := h.srv.GetById(id)
	if err != nil {
		h.Fail(c, apperr.ErrQueryFailed)
		return
	}
	h.Success(c, info)
}

// Create 创建
func (h *BaseHandler[M, Q, C, U]) Create(c *gin.Context) {
	var dto C
	if !h.BindJSON(c, &dto) {
		return
	}

	var m M
	utils.Copy(&m, dto) // DTO → Model

	err := h.srv.Create(&m)
	if err != nil {
		h.Fail(c, apperr.ErrCreateFailed)
		return
	}
	h.Success(c, m)
}

// Update 更新 ✅ 使用 DTO 自带 ToMap()
func (h *BaseHandler[M, Q, C, U]) Update(c *gin.Context) {
	id, ok := h.GetId(c)
	if !ok {
		return
	}
	var dto U
	if !h.BindJSON(c, &dto) {
		return
	}
	mp := utils.StructToMapIgnoreNil(dto)
	_, err := h.srv.Update(id, mp)
	if err != nil {
		h.Fail(c, apperr.ErrUpdateFailed)
		return
	}
	h.Success(c, id)
}

// Delete 删除
func (h *BaseHandler[M, Q, C, U]) Delete(c *gin.Context) {
	id, ok := h.GetId(c)
	if !ok {
		return
	}
	err := h.srv.Delete(id)

	if err != nil {
		h.Fail(c, apperr.ErrDeleteFailed)
		return
	}
	h.Success(c, id)
}

// BatchDelete 批量删除
func (h *BaseHandler[M, Q, C, U]) BatchDelete(c *gin.Context) {
	ids, ok := h.GetIds(c)
	if !ok {
		return
	}
	err := h.srv.BatchDelete(ids)
	if err != nil {
		h.Fail(c, apperr.ErrDeleteFailed)
		return
	}
	h.Success(c, ids)
}

// ==================== 公共基础方法 ====================
type BaseHandlerCommon struct{}

func bindQuery[Q BaseQueryInterface](c *gin.Context) (Q, bool) {
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

func (bh *BaseHandlerCommon) BindJSON(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return false
	}
	return true
}

func (bh *BaseHandlerCommon) GetId(c *gin.Context) (int64, bool) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailBy(c, apperr.ErrInvalidId)
		return 0, false
	}
	return id, true
}

func (bh *BaseHandlerCommon) GetIds(c *gin.Context) ([]int64, bool) {
	idsStr := c.Param("ids")
	ids, err := utils.StrToNumberArray[int64](idsStr, ",")
	if err != nil {
		response.FailBy(c, apperr.ErrInvalidId)
		return nil, false
	}
	return ids, true
}

func (bh *BaseHandlerCommon) Success(c *gin.Context, data any) {
	response.Success(c, data)
}

func (bh *BaseHandlerCommon) Fail(c *gin.Context, errDef apperr.ErrorDef) {
	response.FailBy(c, errDef)
}
