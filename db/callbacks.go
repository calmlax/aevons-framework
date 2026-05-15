package db

import (
	"github.com/calmlax/aevons-framework/auth"
	"gorm.io/gorm"
)

const (
	auditCreateCallbackName = "aevons:audit_before_create"
	auditUpdateCallbackName = "aevons:audit_before_update"
)

// registerCallbacks 注册全局 GORM 回调：
// 1. 创建数据时自动补充 CreatedBy / UpdatedBy / IsDeleted
// 2. 更新数据时自动刷新 UpdatedBy
// 当前用户从 db.WithContext(ctx) 传入的上下文中读取，未登录时默认写入 0。
func registerCallbacks(db *gorm.DB) error {
	if err := db.Callback().Create().Before("gorm:create").Register(auditCreateCallbackName, fillAuditForCreate); err != nil {
		return err
	}
	if err := db.Callback().Update().Before("gorm:update").Register(auditUpdateCallbackName, fillAuditForUpdate); err != nil {
		return err
	}
	return nil
}

func fillAuditForCreate(tx *gorm.DB) {
	if tx == nil || tx.Statement == nil || tx.Statement.Schema == nil {
		return
	}

	userID := currentUserID(tx)
	setFieldIfPresent(tx, "CreatedBy", userID, true)
	setFieldIfPresent(tx, "UpdatedBy", userID, false)
	setFieldIfPresent(tx, "IsDeleted", false, true)
}

func fillAuditForUpdate(tx *gorm.DB) {
	if tx == nil || tx.Statement == nil || tx.Statement.Schema == nil {
		return
	}

	userID := currentUserID(tx)
	setFieldIfPresent(tx, "UpdatedBy", userID, false)
}

func currentUserID(tx *gorm.DB) int64 {
	if tx == nil || tx.Statement == nil || tx.Statement.Context == nil {
		return 0
	}
	userID, err := auth.GetCurrentUserId(tx.Statement.Context)
	if err != nil {
		return 0
	}
	return userID
}

func setFieldIfPresent(tx *gorm.DB, fieldName string, value any, onlyWhenZero bool) {
	field := tx.Statement.Schema.LookUpField(fieldName)
	if field == nil {
		return
	}

	if onlyWhenZero {
		current, zero := field.ValueOf(tx.Statement.Context, tx.Statement.ReflectValue)
		if !zero {
			switch v := current.(type) {
			case int64:
				if v != 0 {
					return
				}
			case bool:
				if v {
					return
				}
			default:
				return
			}
		}
	}

	_ = field.Set(tx.Statement.Context, tx.Statement.ReflectValue, value)
}
