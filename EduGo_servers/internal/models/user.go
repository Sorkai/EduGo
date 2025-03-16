package models

import "time"

import (
	"golang.org/x/crypto/bcrypt"
)

// 用户类型常量
const (
	RoleSuperAdmin = "super_admin" // 超级管理员
	RoleAdmin      = "admin"       // 管理员
	RoleTeacher    = "teacher"     // 教师
	RoleStudent    = "student"     // 学生
	RoleParent     = "parent"      // 家长
)

type User struct {
	ID        int64  `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Email     string `gorm:"unique;not null"`
	Role      string `gorm:"default:'student'"`
	Status    string `gorm:"default:'active'"`
	FirstName string `gorm:"size:50"`
	LastName  string `gorm:"size:50"`
	CreatedAt time.Time
	UpdatedAt time.Time
	LastLoginAt *time.Time
}

// HashPassword 加密密码
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
