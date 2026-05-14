package base

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/calmlax/aevons-framework/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// BaseQueryInterface 分页查询接口定义
type BaseQueryInterface interface {
	GetOrder() (field string, isDesc bool, err error) // 这里已同步修正
	Normalize()
	Offset() int
	GetPageSize() int
}

// BaseRepository 泛型仓储基类接口
// 提供全功能通用 CRUD、分页、条件查询、事务等能力，所有业务仓储均可继承
type BaseRepository[Model any] interface {

	// Create 创建单条数据
	Create(entity *Model) error

	// CreateBatch 批量创建数据
	CreateBatch(entities []*Model) error

	// GetById 根据主键 ID 查询单条数据
	GetById(id int64) (*Model, error)

	// Update 根据 ID 更新数据，支持 map 形式批量更新字段
	Update(id int64, updates map[string]any) (*Model, error)

	// Delete 根据 ID 删除数据
	Delete(id int64) error

	// BatchDelete 根据 ID 列表批量删除
	BatchDelete(ids []int64) error

	// List 根据条件查询列表
	List(q BaseQueryInterface) ([]Model, error)

	// ListAll 查询全部数据
	ListAll() ([]Model, error)

	// Page 分页查询（默认按 id 倒序）
	Page(q BaseQueryInterface) ([]Model, int64, error)

	// Count 根据条件统计总条数
	Count(q BaseQueryInterface) (int64, error)

	// ListByIds 根据 ID 列表查询列表
	ListByIds(ids []int64) ([]Model, error)

	// GetByField 根据单个字段查询单条数据
	GetByField(field string, value any) (*Model, error)

	// ListByField 根据单个字段查询列表
	ListByField(field string, value any) ([]Model, error)

	// ListByFields 根据多个字段等值查询列表
	ListByFields(conds map[string]any) ([]Model, error)

	// ExistByField 判断指定字段值是否存在（用于唯一校验）
	ExistByField(field string, value any) (bool, error)

	// ExistByFieldExcludeId 判断字段值是否存在（排除 excludeId，编辑时使用）
	ExistByFieldExcludeId(field string, value any, excludeId int64) (bool, error)

	// Transaction 执行事务
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error

	// ModelDB 获取当前模型的 *gorm.DB 对象，用于自定义复杂查询
	ModelDB() *gorm.DB
}

// baseRepository 泛型仓储基类实现
// 嵌入 *gorm.DB，实现 BaseRepository 接口所有方法
type baseRepository[Model any] struct {
	db *gorm.DB // GORM 数据库连接
}

// NewBaseRepository 创建 BaseRepository 实例
func NewBaseRepository[Model any](db *gorm.DB) BaseRepository[Model] {
	return &baseRepository[Model]{db: db}
}

// ======================================================
// 基础 CRUD 实现
// ======================================================

// Create 创建单条数据
func (r *baseRepository[Model]) Create(entity *Model) error {
	return r.db.Create(entity).Error
}

// CreateBatch 批量创建数据
// 传入空切片时直接返回，避免执行无效 SQL
func (r *baseRepository[Model]) CreateBatch(entities []*Model) error {
	if len(entities) == 0 {
		return nil
	}
	return r.db.CreateInBatches(entities, len(entities)).Error
}

// GetById 根据 ID 查询
// 未找到时返回 nil, gorm.ErrRecordNotFound，不返回零值指针
func (r *baseRepository[Model]) GetById(id int64) (*Model, error) {
	var entity Model
	if err := r.db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// Update 根据 ID 更新字段
func (r *baseRepository[Model]) Update(id int64, updates map[string]any) (*Model, error) {
	var entity Model
	// 先查询数据是否存在
	if err := r.db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	// 执行更新
	if err := r.db.Model(&entity).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// Delete 根据 ID 删除
func (r *baseRepository[Model]) Delete(id int64) error {
	return r.db.Delete(new(Model), id).Error
}

// BatchDelete 批量删除
// 传入空切片时直接返回，避免执行 IN () 的非法 SQL
func (r *baseRepository[Model]) BatchDelete(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Where("id IN (?)", ids).Delete(new(Model)).Error
}

// List 条件查询列表
func (r *baseRepository[Model]) List(q BaseQueryInterface) ([]Model, error) {
	var list []Model
	query := r.ApplyQuery(r.db, q)
	field, isDesc, e := q.GetOrder()
	if e == nil {
		query = query.Order(clause.OrderBy{Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: field}, Desc: isDesc},
		}})
	}
	err := query.Find(&list).Error
	return list, err
}

// ListAll 查询所有数据
func (r *baseRepository[Model]) ListAll() ([]Model, error) {
	var list []Model
	err := r.db.Find(&list).Error
	return list, err
}

// Page 分页查询（默认排序 id DESC）
func (r *baseRepository[Model]) Page(q BaseQueryInterface) ([]Model, int64, error) {
	var list []Model
	var total int64

	// 构建查询条件（Count 与 Find 共用同一基础条件，互不干扰）
	query := r.ApplyQuery(r.db, q)

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 无数据直接返回
	if total == 0 {
		return list, 0, nil
	}

	// 分页查询，order 为空时不附加排序
	q2 := query.Offset(q.Offset()).Limit(q.GetPageSize())

	field, isDesc, e := q.GetOrder()
	if e == nil {
		q2 = q2.Order(clause.OrderBy{Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: field}, Desc: isDesc},
		}})
	}
	err := q2.Find(&list).Error
	return list, total, err
}

// Count 统计满足条件的数据量
func (r *baseRepository[Model]) Count(q BaseQueryInterface) (int64, error) {
	var total int64
	err := r.ApplyQuery(r.db, q).Count(&total).Error
	return total, err
}

func (r *baseRepository[Model]) ListByIds(ids []int64) ([]Model, error) {
	var list []Model
	err := r.db.Where("id IN (?)", ids).Find(&list).Error
	return list, err
}

// GetByField 根据字段查询单条
// 未找到时返回 nil, gorm.ErrRecordNotFound，不返回零值指针
func (r *baseRepository[Model]) GetByField(field string, value any) (*Model, error) {
	var entity Model
	if err := r.db.Where(field+" = ?", value).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// ListByField 根据字段查询列表
func (r *baseRepository[Model]) ListByField(field string, value any) ([]Model, error) {
	var list []Model
	err := r.db.Where(field+" = ?", value).Find(&list).Error
	return list, err
}

// ListByFields 多字段等值查询
func (r *baseRepository[Model]) ListByFields(conds map[string]any) ([]Model, error) {
	var list []Model
	err := r.db.Where(conds).Find(&list).Error
	return list, err
}

// ExistByField 判断字段值是否存在
func (r *baseRepository[Model]) ExistByField(field string, value any) (bool, error) {
	var count int64
	err := r.db.Model(new(Model)).Where(field+" = ?", value).Limit(1).Count(&count).Error
	return count > 0, err
}

// ExistByFieldExcludeId 判断字段值是否存在（排除 excludeId，编辑时使用）
func (r *baseRepository[Model]) ExistByFieldExcludeId(field string, value any, excludeId int64) (bool, error) {
	var count int64
	err := r.db.Model(new(Model)).
		Where(field+" = ?", value).
		Where("id != ?", excludeId).
		Limit(1).Count(&count).Error
	return count > 0, err
}

// Transaction 执行数据库事务
// 传入闭包函数，若返回 error 则自动回滚，否则提交
func (r *baseRepository[Model]) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

// ModelDB 获取当前模型的 GORM DB 对象
// 用于自定义复杂查询、原生 SQL、联表查询等
func (r *baseRepository[Model]) ModelDB() *gorm.DB {
	return r.db.Model(new(Model))
}

// ApplyQuery 根据结构体字段的 `q` tag 自动构建 GORM 查询条件
//
// 设计目标：
// 1. 支持通过 tag 定义查询操作符（eq / like / in 等）
// 2. 自动忽略 nil 字段（用于 PATCH / 动态查询）
// 3. 支持指针类型（区分"未传值"和"零值"）
// 4. 防 SQL 注入（统一使用占位符 ?）
//
// 使用示例：
//
//	type UserQuery struct {
//	    Name *string `q:"like"`
//	    Age  *int    `q:"gte"`
//	}
//
// 注意：
// - q 必须是 struct 或 *struct，传 nil 时直接返回不加条件
// - 推荐所有字段使用指针类型（避免零值歧义）
// - 非指针 string 字段为空串时将被跳过，不参与查询
func (r *baseRepository[Model]) ApplyQuery(db *gorm.DB, q BaseQueryInterface) *gorm.DB {
	db = db.Model(new(Model))

	// q 为 nil 时直接返回，不加任何条件
	if q == nil {
		return db
	}

	val := reflect.ValueOf(q)

	// 如果是指针，先判断是否为 nil，再取实际值
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return db
		}
		val = val.Elem()
	}

	// 非 struct 类型直接返回，避免 panic
	if val.Kind() != reflect.Struct {
		return db
	}

	typ := val.Type()

	// 遍历结构体字段
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		typeField := typ.Field(i)

		// 获取自定义查询操作符（如 eq / like）
		op := typeField.Tag.Get("q")
		if op == "" {
			continue // 没有定义查询规则，跳过
		}

		// 如果是指针且为 nil，说明未传值 → 不参与查询
		if field.Kind() == reflect.Ptr && field.IsNil() {
			continue
		}

		// 获取数据库字段名（优先 gorm column，其次 json，最后字段名）
		col := r.getColumn(typeField)

		// 获取实际值（解引用指针）
		var value any
		if field.Kind() == reflect.Ptr {
			value = field.Elem().Interface()
		} else {
			value = field.Interface()
		}

		if utils.IsEmpty(value) {
			continue
		}

		//fmt.Printf("-------ApplyQuery------- field: %v,value: %v,\n", col, value)

		// 根据操作符构建 SQL 条件
		switch op {

		// = 等于
		case "eq":
			db = db.Where(col+" = ?", value)

		// != 不等于
		case "ne":
			db = db.Where(col+" <> ?", value)

		// > 大于
		case "gt":
			db = db.Where(col+" > ?", value)

		// >= 大于等于
		case "gte":
			db = db.Where(col+" >= ?", value)

		// < 小于
		case "lt":
			db = db.Where(col+" < ?", value)

		// <= 小于等于
		case "lte":
			db = db.Where(col+" <= ?", value)

		// 模糊匹配：%value%
		case "like":
			db = db.Where(col+" LIKE ?", "%"+fmt.Sprint(value)+"%")

		// 左模糊：%value
		case "like_l":
			db = db.Where(col+" LIKE ?", "%"+fmt.Sprint(value))

		// 右模糊：value%
		case "like_r":
			db = db.Where(col+" LIKE ?", fmt.Sprint(value)+"%")

		// IN 查询（value 必须是 slice）
		// 空切片时跳过，避免生成 IN () 的非法 SQL
		case "in":
			v := reflect.ValueOf(value)
			if v.Kind() == reflect.Slice && v.Len() == 0 {
				continue
			}
			db = db.Where(col+" IN ?", value)

		// NOT IN 查询，空切片时跳过
		case "not_in":
			v := reflect.ValueOf(value)
			if v.Kind() == reflect.Slice && v.Len() == 0 {
				continue
			}
			db = db.Where(col+" NOT IN ?", value)

		// BETWEEN 区间查询（必须是长度为 2 的 slice/array，否则跳过）
		case "between":
			v := reflect.ValueOf(value)
			if (v.Kind() == reflect.Slice || v.Kind() == reflect.Array) && v.Len() == 2 {
				db = db.Where(
					col+" BETWEEN ? AND ?",
					v.Index(0).Interface(),
					v.Index(1).Interface(),
				)
			}

		// IS NULL
		case "is_null":
			db = db.Where(col + " IS NULL")

		// IS NOT NULL
		case "not_null":
			db = db.Where(col + " IS NOT NULL")
		}
	}

	return db
}

// getColumn 获取数据库字段名
//
// 优先级：
// 1. gorm tag: column:xxx
// 2. json tag（取逗号前第一段）
// 3. 结构体字段名
//
// 示例：
// `gorm:"column:resource_key"` → resource_key
// `json:"resourceKey,omitempty"` → resourceKey
func (r *baseRepository[Model]) getColumn(f reflect.StructField) string {

	// 优先解析 gorm tag 中的 column:xxx
	// ⚠️ 修复：gormTag[idx:] 形如 "column:created_at;type:timestamp"
	//    必须先取 column: 后的部分，再按 ";" 截断，而非直接 split(":")
	gormTag := f.Tag.Get("gorm")
	if gormTag != "" {
		if idx := strings.Index(gormTag, "column:"); idx != -1 {
			colStr := gormTag[idx+len("column:"):]
			// 遇到 ; 或 , 则截断（gorm tag 多属性用 ; 分隔）
			if end := strings.IndexAny(colStr, ";,"); end != -1 {
				return colStr[:end]
			}
			return colStr
		}
	}

	// fallback：json tag（取第一段，忽略 omitempty 等修饰符）
	jsonTag := f.Tag.Get("json")
	if jsonTag != "" && jsonTag != "-" {
		name := strings.SplitN(jsonTag, ",", 2)[0]
		if name != "" {
			return name
		}
	}

	// fallback：结构体字段名
	return f.Name
}
