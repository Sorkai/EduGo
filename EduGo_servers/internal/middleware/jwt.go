package middleware

import (
	"EduGo_servers/internal/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your-secret-key") // 生产环境应从环境变量获取

type Claims struct {
	UserID int64 `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		// 尝试使用Claims结构体解析
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			// 如果使用Claims结构体解析失败，尝试使用MapClaims解析
			mapClaims := jwt.MapClaims{}
			token, err = jwt.ParseWithClaims(tokenString, mapClaims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err != nil || !token.Valid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}

			// 从MapClaims中获取用户ID和角色
			if userID, ok := mapClaims["user_id"].(float64); ok {
				c.Set("userID", int64(userID))
			} else if userID, ok := mapClaims["id"].(float64); ok {
				c.Set("userID", int64(userID))
			}

			if role, ok := mapClaims["role"].(string); ok {
				c.Set("role", role)
			}
		} else if token.Valid {
			// 如果使用Claims结构体解析成功，直接设置用户ID和角色
			c.Set("userID", claims.UserID)
			c.Set("role", claims.Role)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// 设置用户名（如果有）
		// 根据token.Claims的实际类型进行不同的处理
		switch claims := token.Claims.(type) {
		case *Claims:
			// 如果是自定义Claims类型，目前没有username字段，可以根据需要添加
			// 如果将来Claims结构体中添加了username字段，可以在这里设置
		case jwt.MapClaims:
			// 如果是MapClaims类型，检查username字段
			if claims["username"] != nil {
				c.Set("username", claims["username"])
			}
		}

		c.Next()
	}
}

// 权限控制中间件

// SuperAdminOnly 只允许超级管理员访问
func SuperAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != models.RoleSuperAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "超级管理员权限需要"})
			return
		}
		c.Next()
	}
}

// AdminOnly 只允许管理员及以上角色访问
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "管理员权限需要"})
			return
		}
		
		roleStr, ok := role.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "角色类型错误"})
			return
		}
		
		if roleStr != models.RoleSuperAdmin && roleStr != models.RoleAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "管理员权限需要"})
			return
		}
		
		c.Next()
	}
}

// TeacherOnly 只允许教师及以上角色访问
func TeacherOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "教师权限需要"})
			return
		}
		
		roleStr, ok := role.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "角色类型错误"})
			return
		}
		
		if roleStr != models.RoleSuperAdmin && roleStr != models.RoleAdmin && roleStr != models.RoleTeacher {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "教师权限需要"})
			return
		}
		
		c.Next()
	}
}

// StudentOnly 只允许学生及以上角色访问
func StudentOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "学生权限需要"})
			return
		}
		
		roleStr, ok := role.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "角色类型错误"})
			return
		}
		
		if roleStr != models.RoleSuperAdmin && roleStr != models.RoleAdmin && 
		   roleStr != models.RoleTeacher && roleStr != models.RoleStudent {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "学生权限需要"})
			return
		}
		
		c.Next()
	}
}

// ParentOnly 只允许家长及以上角色访问
func ParentOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "家长权限需要"})
			return
		}
		
		roleStr, ok := role.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "角色类型错误"})
			return
		}
		
		if roleStr != models.RoleSuperAdmin && roleStr != models.RoleAdmin && 
		   roleStr != models.RoleTeacher && roleStr != models.RoleParent {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "家长权限需要"})
			return
		}
		
		c.Next()
	}
}
