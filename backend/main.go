package main

import (
	"fmt"
	"log"
	"os"

	"github.com/007Secret/007Secret/api"
	"github.com/007Secret/007Secret/config"
	"github.com/007Secret/007Secret/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 显式按 GIN_MODE 设置运行模式（release=关闭调试日志，未设置则 debug）。
	// gin 只在包 init 时读一次该变量；若 .env 里配置了 GIN_MODE，需要在这里重新应用。
	gin.SetMode(os.Getenv(gin.EnvGinMode))

	// 初始化配置
	if err := config.Init(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化Redis连接
	if err := utils.InitRedis(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	// 启动服务器
	addr := ":" + config.AppConfig.Port
	fmt.Printf("Server is running on http://localhost%s\n", addr)
	log.Fatal(setupRouter().Run(addr))
}

// setupRouter 构建 Gin 路由。单独拆出以便测试直接校验路由表。
func setupRouter() *gin.Engine {
	r := gin.Default()

	// 安全响应头（API 直接访问时也生效；SPA 文档的头由 nginx 设置）
	r.Use(api.NewSecurityHeaders())

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	// 注册路由处理器
	signalingHub := api.NewSignalingHub()
	r.POST("/api/secret", api.NewBodyLimit(api.MaxRequestBodyBytes), api.CreateSecretHandler)
	r.POST("/api/secret/get", api.NewBodyLimit(api.MaxRequestBodyBytes), api.NewSecretGetGuard(), api.GetSecretHandler)
	r.GET("/api/transfer/signal", api.NewSignalGuard(), signalingHub.Handle)

	return r
}
