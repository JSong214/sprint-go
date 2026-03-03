package repository

import (
	"github.com/JSong214/sprint-go/internal/model"
	"gorm.io/gorm"
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	// Create 创建用户
	Create(user *model.User) error

	// FindByID 根据 ID 查找用户
	FindByID(id int64) (*model.User, error)

	// FindByUsername 根据用户名查找用户
	FindByUsername(username string) (*model.User, error)

	// FindByEmail 根据邮箱查找用户
	FindByEmail(email string) (*model.User, error)

	// ExistsByUsername 检查用户名是否已存在
	ExistsByUsername(username string) (bool, error)

	// ExistsByEmail 检查邮箱是否已存在
	ExistsByEmail(email string) (bool, error)
}

// userRepositoryImpl UserRepository 的实现
type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建 UserRepository 实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

// Create 创建用户
func (r *userRepositoryImpl) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByID 根据 ID 查找用户
func (r *userRepositoryImpl) FindByID(id int64) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (r *userRepositoryImpl) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *userRepositoryImpl) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ExistsByUsername 检查用户名是否已存在
func (r *userRepositoryImpl) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail 检查邮箱是否已存在
func (r *userRepositoryImpl) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
