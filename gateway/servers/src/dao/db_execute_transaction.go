package dao

import (
	"gateway/global"
	"gateway/protoc/pbs"
	"gorm.io/gorm"
)

// ExecuteTransaction 定义通用事务处理函数
func ExecuteTransaction(txFunc func(tx *gorm.DB) int32) int32 {
	db := global.GVA_USER_DB
	// 开启事务
	tx := db.Begin()
	if err := tx.Error; err != nil {
		global.GVA_LOG.Infof("事务开启失败: %v", err)
		//return int32(pbs.Code_ServerError)
	}

	// 执行事务逻辑
	code := txFunc(tx)
	if code != int32(pbs.Code_OK) {
		// 如果事务逻辑失败，回滚事务
		if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
			global.GVA_LOG.Infof("事务回滚失败 %v 原始错误  %v", rollbackErr, code)
			//return int32(pbs.Code_ServerError)
		}
		return code
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		global.GVA_LOG.Infof("事务提交失败: %v", err)
		//return int32(pbs.Code_ServerError)
	}

	return int32(pbs.Code_OK)
}
