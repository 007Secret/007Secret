package main

import (
	"fmt"
	"log"

	"github.com/007Secret/007Secret/api"
	"github.com/007Secret/007Secret/config"
	"github.com/007Secret/007Secret/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化Redis连接
	if err := utils.InitRedis(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	// 创建 Gin 路由
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	// 注册路由处理器
	r.POST("/api/secret", api.CreateSecretHandler)
	r.GET("/api/secret/get", api.GetSecretHandler)

	// 启动服务器
	addr := ":" + config.AppConfig.Port
	fmt.Printf("Server is running on http://localhost%s\n", addr)
	log.Fatal(r.Run(addr))
}
