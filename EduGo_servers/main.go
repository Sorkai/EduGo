package main

import (
	"EduGo_servers/internal/controllers"
	"EduGo_servers/internal/database"
	"EduGo_servers/internal/middleware"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载.env文件
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// 初始化数据库连接
	dbErr := database.InitDB(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"), 
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	if dbErr != nil {
		log.Fatalf("Failed to connect to database: %v", dbErr)
	}

	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 公开路由
		v1.POST("/register", controllers.Register)
		v1.POST("/login", controllers.Login)

		// 需要认证的路由
		auth := v1.Group("/")
		auth.Use(middleware.JWTMiddleware())
		{
			// 用户相关
			auth.GET("/user", controllers.GetUserProfile)
			auth.PUT("/user", controllers.UpdateUser)
			auth.PUT("/user/password", controllers.ResetPassword)
			auth.POST("/logout", controllers.Logout)
			auth.POST("/refresh", controllers.RefreshToken)

			// 管理员路由
			admin := auth.Group("/admin")
			admin.Use(middleware.AdminOnly())
			{
				// TODO: 添加管理员路由
			}
		}
	}

	// 获取端口配置，默认为10086
	port := os.Getenv("PORT")
	if port == "" {
		port = "10086"
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
