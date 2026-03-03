package service

import (
	"errors"
	"strings"

	"github.com/JSong214/sprint-go/internal/model"
	"github.com/JSong214/sprint-go/internal/repository"
	"github.com/JSong214/sprint-go/internal/service/dto"
	myjwt "github.com/JSong214/sprint-go/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

// AuthService 认证服务接口（方法签名严格对应 OpenAPI 定义）
type AuthService interface {
	// Register 用户注册
	Register(req *dto.RegisterRequest) (*dto.UserResponse, error)

	// Login 用户登录
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)

	// Logout 用户登出
	Logout(userID int64) error
}

// authServiceImpl AuthService 的实现结构体
type authServiceImpl struct {
	userRepo repository.UserRepository
}

// NewAuthService 创建 AuthService 实例
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authServiceImpl{
		userRepo: userRepo,
	}
}

// Register 用户注册
func (s *authServiceImpl) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
	// 1. 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// 2. 检查邮箱是否已存在
	exists, err = s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// 3. 使用 bcrypt 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// 4. 创建 User 实体
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 5. 转换为 UserResponse DTO
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

// Login 用户登录
func (s *authServiceImpl) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// 1. 根据 login 字段判断是用户名还是邮箱，查询用户
	var user *model.User
	var err error

	if strings.Contains(req.Login, "@") {
		user, err = s.userRepo.FindByEmail(req.Login)
	} else {
		user, err = s.userRepo.FindByUsername(req.Login)
	}
	if err != nil {
		return nil, err // gorm.ErrRecordNotFound 会被 handler 层捕获
	}

	// 2. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 3. 生成 JWT token
	token, err := myjwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// 4. 组装 LoginResponse
	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// Logout 用户登出
// MVP 阶段：前端删除 token 即可，后端暂不需要额外操作
// 后续扩展：可在 Redis 中添加 token 黑名单
func (s *authServiceImpl) Logout(userID int64) error {
	return nil
}
