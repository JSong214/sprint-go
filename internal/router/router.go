package router

import (
	"github.com/JSong214/sprint-go/internal/handler"
	"github.com/JSong214/sprint-go/internal/repository"
	"github.com/JSong214/sprint-go/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes 配置所有路由
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// 初始化依赖注入
	// Repository 层
	userRepo := repository.NewUserRepository(db)

	// Service 层
	authService := service.NewAuthService(userRepo)

	// Handler 层
	authHandler := handler.NewAuthHandler(authService)

	// API 路由组
	api := r.Group("/api")
	{
		// 认证路由组
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout) // TODO: 需要 JWT 中间件保护
		}
	}
}
