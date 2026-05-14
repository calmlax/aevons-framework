package base

import (
	"errors"

	errdef "github.com/calmlax/aevons-framework/err"

	"gorm.io/gorm"
)

// ==================== 泛型 BaseService 基类 ====================

type BaseService[Model any, Q BaseQueryInterface] interface {
	// 列表
	List(q Q) ([]Model, error)
	// 分页
	Page(q Q) (*PageResult[Model], error)
	// 根据ID查询
	GetById(id int64) (*Model, error)
	// 创建
	Create(entity *Model) error
	// 批量创建
	CreateBatch(entities []*Model) error
	// 更新
	Update(id int64, updates map[string]any) (*Model, error)
	// 删除
	Delete(id int64) error
	// 批量删除
	BatchDelete(ids []int64) error

	// ==================== 扩展常用方法 ====================
	// 统计数量
	Count(q Q) (int64, error)
	// 判断字段是否存在（新增去重）
	ExistField(field string, value any) (bool, error)
	// 判断字段是否存在（排除ID，编辑去重）
	ExistFieldExcludeId(field string, value any, id int64) (bool, error)
	// 根据单字段查询单条
	GetByField(field string, value any) (*Model, error)
	// 根据单字段查询列表
	ListByField(field string, value any) ([]Model, error)
	// 多字段等值查询
	ListByFields(conds map[string]any) ([]Model, error)
}

// baseService 实现
type baseService[Model any, Q BaseQueryInterface] struct {
	repo BaseRepository[Model]
}

func NewBaseService[Model any, Q BaseQueryInterface](repo BaseRepository[Model]) BaseService[Model, Q] {
	return &baseService[Model, Q]{repo: repo}
}

// ==================== 原有方法实现 ====================

func (s *baseService[Model, Q]) List(q Q) ([]Model, error) {
	return s.repo.List(q)
}

func (s *baseService[Model, Q]) Page(q Q) (*PageResult[Model], error) {
	q.Normalize()
	list, total, err := s.repo.Page(q)
	if err != nil {
		return nil, err
	}
	return &PageResult[Model]{Rows: list, Total: total}, nil
}

func (s *baseService[Model, Q]) GetById(id int64) (*Model, error) {
	entity, err := s.repo.GetById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errdef.ErrNotFound
		}
		return nil, err
	}
	return entity, nil
}

func (s *baseService[Model, Q]) Create(entity *Model) error {
	return s.repo.Create(entity)
}

func (s *baseService[Model, Q]) Update(id int64, updates map[string]any) (*Model, error) {
	if len(updates) == 0 {
		return nil, errdef.ErrNoUpdateField
	}
	entity, err := s.repo.Update(id, updates)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errdef.ErrNotFound
		}
		return nil, err
	}
	return entity, nil
}

func (s *baseService[Model, Q]) Delete(id int64) error {
	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errdef.ErrNotFound
		}
		return err
	}
	return nil
}

func (s *baseService[Model, Q]) BatchDelete(ids []int64) error {
	return s.repo.BatchDelete(ids)
}

// ==================== 新增常用方法 ====================

// CreateBatch 批量创建
func (s *baseService[Model, Q]) CreateBatch(entities []*Model) error {
	return s.repo.CreateBatch(entities)
}

// Count 统计数量
func (s *baseService[Model, Q]) Count(q Q) (int64, error) {
	return s.repo.Count(q)
}

// ExistField 判断字段是否存在
func (s *baseService[Model, Q]) ExistField(field string, value any) (bool, error) {
	return s.repo.ExistByField(field, value)
}

// ExistFieldExcludeId 判断字段是否存在（排除当前ID）
func (s *baseService[Model, Q]) ExistFieldExcludeId(field string, value any, id int64) (bool, error) {
	return s.repo.ExistByFieldExcludeId(field, value, id)
}

// GetByField 根据单个字段查询单条
func (s *baseService[Model, Q]) GetByField(field string, value any) (*Model, error) {
	return s.repo.GetByField(field, value)
}

// ListByField 根据单个字段查询列表
func (s *baseService[Model, Q]) ListByField(field string, value any) ([]Model, error) {
	return s.repo.ListByField(field, value)
}

// ListByFields 多字段等值查询
func (s *baseService[Model, Q]) ListByFields(conds map[string]any) ([]Model, error) {
	return s.repo.ListByFields(conds)
}
