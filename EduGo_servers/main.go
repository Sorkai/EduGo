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

			// 超级管理员路由
			superAdmin := auth.Group("/super-admin")
			superAdmin.Use(middleware.SuperAdminOnly())
			{
				superAdmin.GET("/users", controllers.GetAllUsers)
				superAdmin.PUT("/users/:id/role", controllers.UpdateUserRole)
			}

			// 管理员路由
			admin := auth.Group("/admin")
			admin.Use(middleware.AdminOnly())
			{
				admin.GET("/users/role/:role", controllers.GetUsersByRole)
				admin.GET("/users/:id", controllers.GetUserByID)
				admin.PUT("/users/:id/status", controllers.UpdateUserStatus)
				
				// 管理员-教师关系
				admin.POST("/relations/teacher", controllers.CreateAdminTeacherRelation)
				admin.GET("/relations/teachers", controllers.GetTeachersByAdmin)
			}
			
			// 教师路由
			teacher := auth.Group("/teacher")
			teacher.Use(middleware.TeacherOnly())
			{
				// 教师-学生关系
				teacher.POST("/relations/student", controllers.CreateTeacherStudentRelation)
				teacher.GET("/relations/students", controllers.GetStudentsByTeacher)
			}
			
			// 学生路由
			student := auth.Group("/student")
			student.Use(middleware.StudentOnly())
			{
				// 学生-家长关系
				student.POST("/relations/parent", controllers.CreateStudentParentRelation)
				student.GET("/relations/parents", controllers.GetParentsByStudent)
			}
			
			// 用户管理页面API
			userManagement := auth.Group("/user-management")
			userManagement.Use(middleware.TeacherOnly()) // 教师及以上角色可访问
			{
				// 根据当前用户角色返回不同的用户列表
				userManagement.GET("/users", controllers.GetAllUsers) // 实际会根据角色权限过滤
				userManagement.GET("/users/role/:role", controllers.GetUsersByRole)
				userManagement.GET("/users/:id", controllers.GetUserByID)
			}
			
			// 智能测评模块API
			assessmentController := controllers.NewAssessmentController()
			assessment := auth.Group("/assessment")
			{
				// 教师路由（教师及以上角色可访问）
				assessmentTeacher := assessment.Group("/teacher")
				assessmentTeacher.Use(middleware.TeacherOnly())
				{
					// 创建测评
					assessmentTeacher.POST("", assessmentController.CreateAssessment)
					// 获取教师创建的所有测评
					assessmentTeacher.GET("", assessmentController.GetTeacherAssessments)
					// 获取测评详情
					assessmentTeacher.GET("/:id", assessmentController.GetTeacherAssessmentDetail)
					// 添加题目到测评
					assessmentTeacher.POST("/:id/question", assessmentController.AddQuestionToAssessment)
					// 发布测评
					assessmentTeacher.PUT("/:id/publish", assessmentController.PublishAssessment)
					// 关闭测评
					assessmentTeacher.PUT("/:id/close", assessmentController.CloseAssessment)
					// 获取测评学生列表
					assessmentTeacher.GET("/:id/students", assessmentController.GetAssessmentStudents)
				}
				
				// 学生路由（学生角色可访问）
				assessmentStudent := assessment.Group("/student")
				assessmentStudent.Use(middleware.StudentOnly())
				{
					// 获取学生可参与的所有测评
					assessmentStudent.GET("", assessmentController.GetStudentAssessments)
					// 获取测评详情
					assessmentStudent.GET("/:id", assessmentController.GetStudentAssessmentDetail)
					// 开始测评
					assessmentStudent.POST("/:id/start", assessmentController.StartAssessment)
					// 提交答案
					assessmentStudent.POST("/:id/submit", assessmentController.SubmitAssessment)
					// 获取测评结果
					assessmentStudent.GET("/:id/result", assessmentController.GetAssessmentResult)
				}
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
