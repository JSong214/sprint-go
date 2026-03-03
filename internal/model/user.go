package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户实体模型
type User struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement;comment:用户ID"`
	Username  string         `json:"username" gorm:"type:varchar(50);uniqueIndex;not null;comment:用户名（唯一）"`
	Email     string         `json:"email" gorm:"type:varchar(100);uniqueIndex;not null;comment:邮箱（唯一）"`
	Password  string         `json:"-" gorm:"column:password_hash;type:varchar(255);not null;comment:密码哈希值"`
	CreatedAt time.Time      `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt time.Time      `json:"-" gorm:"not null;default:CURRENT_TIMESTAMP;comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:软删除时间"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
