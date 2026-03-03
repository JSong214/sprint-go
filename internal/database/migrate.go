package database

import (
	"github.com/JSong214/sprint-go/internal/model"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() error {
	return DB.AutoMigrate(
		&model.User{},
		// 后续可添加其他模型
	)
}
