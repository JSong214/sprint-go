package dto

import "time"

// RegisterRequest 注册请求 DTO
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email,max=100"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

// LoginRequest 登录请求 DTO
type LoginRequest struct {
	Login    string `json:"login" binding:"required"`    // 用户名或邮箱
	Password string `json:"password" binding:"required"` // 密码
}

// UserResponse 用户响应 DTO（不包含敏感信息）
type UserResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginResponse 登录响应 DTO
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
