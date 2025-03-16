package models

import "time"

// 关系类型常量
const (
	RelationAdminTeacher  = "admin_teacher"   // 管理员-教师关系
	RelationTeacherStudent = "teacher_student" // 教师-学生关系
	RelationStudentParent = "student_parent"   // 学生-家长关系
)

// UserRelation 用户关系模型
type UserRelation struct {
	ID           int64     `gorm:"primaryKey"`
	UserID       int64     `gorm:"not null;index"` // 关系发起者ID
	RelatedUserID int64    `gorm:"not null;index"` // 关系接收者ID
	RelationType string    `gorm:"not null"`       // 关系类型
	Status       string    `gorm:"default:'active'"` // 关系状态：active, inactive
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// AdminTeacherRelation 管理员-教师关系
// 管理员可以管理多个教师，教师归属于一个或多个管理员
type AdminTeacherRelation struct {
	UserRelation
	Department string `gorm:"size:100"` // 部门
	Position   string `gorm:"size:100"` // 职位
}

// TeacherStudentRelation 教师-学生关系
// 教师可以教授多个学生，学生可以有多个教师
type TeacherStudentRelation struct {
	UserRelation
	CourseID   int64  `gorm:"index"` // 课程ID
	CourseName string `gorm:"size:100"` // 课程名称
	Semester   string `gorm:"size:50"`  // 学期
}

// StudentParentRelation 学生-家长关系
// 学生可以有多个家长，家长可以有多个孩子
type StudentParentRelation struct {
	UserRelation
	Relationship string `gorm:"size:50"` // 关系：father, mother, guardian等
}
