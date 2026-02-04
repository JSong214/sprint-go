package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/JSong214/sprint-go/internal/config"
	"github.com/JSong214/sprint-go/internal/database"
	"github.com/JSong214/sprint-go/internal/router/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	configPath := "configs/config.yaml"
	if err := config.LoadConfig(configPath); err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	// 初始化数据库
	if err := database.InitDB(&config.AppConfig.Database); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer database.CloseDB()

	// 设置Gin模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 创建Gin引擎
	r := gin.Default()

	// 配置CORS
	r.Use(middleware.CORS())

	// 检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// 启动服务器
	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	log.Printf("服务器启动在端口 %s", addr)

	// 优雅关闭
	go func() {
		if err := r.Run(addr); err != nil {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("服务器正在关闭...")
}
