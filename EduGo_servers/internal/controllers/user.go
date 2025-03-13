package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"EduGo_servers/internal/models"
	"EduGo_servers/internal/repository"
	"EduGo_servers/internal/database"
)

// Login 处理用户登录请求
func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的输入数据"})
		return
	}

	userRepo := repository.NewUserRepository(database.DB)
	user, err := userRepo.GetUserByUsername(c.Request.Context(), input.Username)
	if err != nil {
		log.Printf("获取用户信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	token, err := generateJWT(user.ID, user.Username)
	if err != nil {
		log.Printf("生成JWT失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
	})
}

// Register 处理用户注册请求
func Register(c *gin.Context) {
	var input struct {
		Username  string `json:"username" binding:"required"`
		Password  string `json:"password" binding:"required,min=8"`
		Email     string `json:"email" binding:"required,email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRepo := repository.NewUserRepository(database.DB)
	
	if exists := userRepo.UserExists(input.Username, input.Email); exists {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名或邮箱已被注册"})
		return
	}

	if err := validatePasswordStrength(input.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Username:  input.Username,
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	if err := user.HashPassword(input.Password); err != nil {
		log.Printf("密码加密失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	if err := userRepo.CreateUser(c.Request.Context(), &user); err != nil {
		log.Printf("创建用户失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "用户注册成功",
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
		},
	})
}

// generateJWT 生成JWT令牌
func generateJWT(userID int64, username string) (string, error) {
	claims := jwt.MapClaims{
		"id":       userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		return "", fmt.Errorf("JWT_SECRET环境变量未设置")
	}
	return token.SignedString(secret)
}

// RefreshToken 刷新JWT令牌
func RefreshToken(c *gin.Context) {
	userID := c.GetInt64("userID")
	username := c.GetString("username")

	token, err := generateJWT(userID, username)
	if err != nil {
		log.Printf("刷新JWT失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "令牌刷新成功",
		"token":   token,
	})
}

// UpdateUser 更新用户信息
func UpdateUser(c *gin.Context) {
	userID := c.GetInt64("userID")

	var input struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的输入数据"})
		return
	}

	userRepo := repository.NewUserRepository(database.DB)
	user, err := userRepo.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("获取用户信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	if input.Email != "" {
		user.Email = input.Email
	}
	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}

	if err := userRepo.UpdateUser(c.Request.Context(), user); err != nil {
		log.Printf("更新用户信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户信息更新成功",
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
		},
	})
}

// validatePasswordStrength 验证密码强度
func validatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("密码至少需要8个字符")
	}

	if matched, _ := regexp.MatchString(`[A-Z]`, password); !matched {
		return fmt.Errorf("密码必须包含至少一个大写字母")
	}

	if matched, _ := regexp.MatchString(`[a-z]`, password); !matched {
		return fmt.Errorf("密码必须包含至少一个小写字母")
	}

	if matched, _ := regexp.MatchString(`[0-9]`, password); !matched {
		return fmt.Errorf("密码必须包含至少一个数字")
	}

	if matched, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`, password); !matched {
		return fmt.Errorf("密码必须包含至少一个特殊字符")
	}

	return nil
}

// ResetPassword 重置用户密码
func ResetPassword(c *gin.Context) {
	userID := c.GetInt64("userID")

	var input struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的输入数据"})
		return
	}

	userRepo := repository.NewUserRepository(database.DB)
	user, err := userRepo.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("获取用户信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "旧密码错误"})
		return
	}

	if err := validatePasswordStrength(input.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := user.HashPassword(input.NewPassword); err != nil {
		log.Printf("密码加密失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	if err := userRepo.UpdateUser(c.Request.Context(), user); err != nil {
		log.Printf("更新用户密码失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "密码重置成功",
	})
}

// Logout 用户注销
func Logout(c *gin.Context) {
	// 在实际应用中，这里可以添加token黑名单机制
	c.JSON(http.StatusOK, gin.H{
		"message": "注销成功",
	})
}
