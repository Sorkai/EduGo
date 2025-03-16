package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"EduGo_servers/internal/database"
	"EduGo_servers/internal/models"
	"EduGo_servers/internal/repository"
)

// GetAllUsers 获取所有用户（超级管理员权限）
func GetAllUsers(c *gin.Context) {
	userRepo := repository.NewUserRepository(database.DB)
	users, err := userRepo.GetAllUsers(c.Request.Context())
	if err != nil {
		log.Printf("获取用户列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	var userList []gin.H
	for _, user := range users {
		userList = append(userList, gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"role":      user.Role,
			"status":    user.Status,
			"createdAt": user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"users": userList,
	})
}

// GetUsersByRole 根据角色获取用户（管理员及以上权限）
func GetUsersByRole(c *gin.Context) {
	role := c.Param("role")
	
	// 验证角色是否有效
	validRoles := map[string]bool{
		models.RoleSuperAdmin: true,
		models.RoleAdmin:      true,
		models.RoleTeacher:    true,
		models.RoleStudent:    true,
		models.RoleParent:     true,
	}
	
	if !validRoles[role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户角色"})
		return
	}
	
	userRepo := repository.NewUserRepository(database.DB)
	users, err := userRepo.GetUsersByRole(c.Request.Context(), role)
	if err != nil {
		log.Printf("获取用户列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	var userList []gin.H
	for _, user := range users {
		userList = append(userList, gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"role":      user.Role,
			"status":    user.Status,
			"createdAt": user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"users": userList,
	})
}

// GetUserByID 根据ID获取用户（管理员及以上权限）
func GetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}
	
	userRepo := repository.NewUserRepository(database.DB)
	user, err := userRepo.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("获取用户信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}
	
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"role":      user.Role,
			"status":    user.Status,
			"createdAt": user.CreatedAt,
		},
	})
}

// UpdateUserRole 更新用户角色（超级管理员权限）
func UpdateUserRole(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}
	
	var input struct {
		Role string `json:"role" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的输入数据"})
		return
	}
	
	// 验证角色是否有效
	validRoles := map[string]bool{
		models.RoleSuperAdmin: true,
		models.RoleAdmin:      true,
		models.RoleTeacher:    true,
		models.RoleStudent:    true,
		models.RoleParent:     true,
	}
	
	if !validRoles[input.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户角色"})
		return
	}
	
	userRepo := repository.NewUserRepository(database.DB)
	user, err := userRepo.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("获取用户信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}
	
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	
	// 不允许修改超级管理员的角色
	if user.Role == models.RoleSuperAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "不能修改超级管理员的角色"})
		return
	}
	
	user.Role = input.Role
	if err := userRepo.UpdateUser(c.Request.Context(), user); err != nil {
		log.Printf("更新用户角色失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "用户角色更新成功",
		"user": gin.H{
			"id":   user.ID,
			"role": user.Role,
		},
	})
}

// UpdateUserStatus 更新用户状态（管理员及以上权限）
func UpdateUserStatus(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}
	
	var input struct {
		Status string `json:"status" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的输入数据"})
		return
	}
	
	// 验证状态是否有效
	validStatus := map[string]bool{
		"active":   true,
		"inactive": true,
		"blocked":  true,
	}
	
	if !validStatus[input.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户状态"})
		return
	}
	
	userRepo := repository.NewUserRepository(database.DB)
	user, err := userRepo.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("获取用户信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}
	
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	
	// 不允许修改超级管理员的状态
	if user.Role == models.RoleSuperAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "不能修改超级管理员的状态"})
		return
	}
	
	// 检查当前用户的角色
	currentUserRole := c.GetString("role")
	
	// 管理员不能修改其他管理员的状态
	if currentUserRole == models.RoleAdmin && user.Role == models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "管理员不能修改其他管理员的状态"})
		return
	}
	
	user.Status = input.Status
	if err := userRepo.UpdateUser(c.Request.Context(), user); err != nil {
		log.Printf("更新用户状态失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "用户状态更新成功",
		"user": gin.H{
			"id":     user.ID,
			"status": user.Status,
		},
	})
}

// 用户关系管理

// CreateAdminTeacherRelation 创建管理员-教师关系（管理员及以上权限）
func CreateAdminTeacherRelation(c *gin.Context) {
	var input struct {
		TeacherID  int64  `json:"teacher_id" binding:"required"`
		Department string `json:"department"`
		Position   string `json:"position"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的输入数据"})
		return
	}
	
	adminID := c.GetInt64("userID")
	
	// 验证教师是否存在且角色是否为教师
	userRepo := repository.NewUserRepository(database.DB)
	teacher, err := userRepo.GetUserByID(c.Request.Context(), input.TeacherID)
	if err != nil {
		log.Printf("获取教师信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}
	
	if teacher == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "教师不存在"})
		return
	}
	
	if teacher.Role != models.RoleTeacher {
		c.JSON(http.StatusBadRequest, gin.H{"error": "指定用户不是教师"})
		return
	}
	
	// 创建关系
	relation := &models.AdminTeacherRelation{
		UserRelation: models.UserRelation{
			UserID:       adminID,
			RelatedUserID: input.TeacherID,
			RelationType: models.RelationAdminTeacher,
			Status:       "active",
		},
		Department: input.Department,
		Position:   input.Position,
	}
	
	relationRepo := repository.NewUserRelationRepository(database.DB)
	if err := relationRepo.CreateAdminTeacherRelation(c.Request.Context(), relation); err != nil {
		log.Printf("创建管理员-教师关系失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "管理员-教师关系创建成功",
		"relation": gin.H{
			"id":         relation.ID,
			"admin_id":   relation.UserID,
			"teacher_id": relation.RelatedUserID,
			"department": relation.Department,
			"position":   relation.Position,
		},
	})
}

// CreateTeacherStudentRelation 创建教师-学生关系（教师及以上权限）
func CreateTeacherStudentRelation(c *gin.Context) {
	var input struct {
		StudentID  int64  `json:"student_id" binding:"required"`
		CourseID   int64  `json:"course_id"`
		CourseName string `json:"course_name"`
		Semester   string `json:"semester"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的输入数据"})
		return
	}
	
	teacherID := c.GetInt64("userID")
	
	// 验证学生是否存在且角色是否为学生
	userRepo := repository.NewUserRepository(database.DB)
	student, err := userRepo.GetUserByID(c.Request.Context(), input.StudentID)
	if err != nil {
		log.Printf("获取学生信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}
	
	if student == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "学生不存在"})
		return
	}
	
	if student.Role != models.RoleStudent {
		c.JSON(http.StatusBadRequest, gin.H{"error": "指定用户不是学生"})
		return
	}
	
	// 创建关系
	relation := &models.TeacherStudentRelation{
		UserRelation: models.UserRelation{
			UserID:       teacherID,
			RelatedUserID: input.StudentID,
			RelationType: models.RelationTeacherStudent,
			Status:       "active",
		},
		CourseID:   input.CourseID,
		CourseName: input.CourseName,
		Semester:   input.Semester,
	}
	
	relationRepo := repository.NewUserRelationRepository(database.DB)
	if err := relationRepo.CreateTeacherStudentRelation(c.Request.Context(), relation); err != nil {
		log.Printf("创建教师-学生关系失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "教师-学生关系创建成功",
		"relation": gin.H{
			"id":          relation.ID,
			"teacher_id":  relation.UserID,
			"student_id":  relation.RelatedUserID,
			"course_id":   relation.CourseID,
			"course_name": relation.CourseName,
			"semester":    relation.Semester,
		},
	})
}

// CreateStudentParentRelation 创建学生-家长关系（学生及以上权限）
func CreateStudentParentRelation(c *gin.Context) {
	var input struct {
		ParentID     int64  `json:"parent_id" binding:"required"`
		Relationship string `json:"relationship"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的输入数据"})
		return
	}
	
	studentID := c.GetInt64("userID")
	
	// 验证家长是否存在且角色是否为家长
	userRepo := repository.NewUserRepository(database.DB)
	parent, err := userRepo.GetUserByID(c.Request.Context(), input.ParentID)
	if err != nil {
		log.Printf("获取家长信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}
	
	if parent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "家长不存在"})
		return
	}
	
	if parent.Role != models.RoleParent {
		c.JSON(http.StatusBadRequest, gin.H{"error": "指定用户不是家长"})
		return
	}
	
	// 创建关系
	relation := &models.StudentParentRelation{
		UserRelation: models.UserRelation{
			UserID:       studentID,
			RelatedUserID: input.ParentID,
			RelationType: models.RelationStudentParent,
			Status:       "active",
		},
		Relationship: input.Relationship,
	}
	
	relationRepo := repository.NewUserRelationRepository(database.DB)
	if err := relationRepo.CreateStudentParentRelation(c.Request.Context(), relation); err != nil {
		log.Printf("创建学生-家长关系失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "学生-家长关系创建成功",
		"relation": gin.H{
			"id":           relation.ID,
			"student_id":   relation.UserID,
			"parent_id":    relation.RelatedUserID,
			"relationship": relation.Relationship,
		},
	})
}

// GetTeachersByAdmin 获取管理员管理的教师列表（管理员及以上权限）
func GetTeachersByAdmin(c *gin.Context) {
	adminID := c.GetInt64("userID")
	
	relationRepo := repository.NewUserRelationRepository(database.DB)
	teachers, err := relationRepo.GetTeachersByAdminID(c.Request.Context(), adminID)
	if err != nil {
		log.Printf("获取教师列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}
	
	var teacherList []gin.H
	for _, teacher := range teachers {
		teacherList = append(teacherList, gin.H{
			"id":        teacher.ID,
			"username":  teacher.Username,
			"email":     teacher.Email,
			"firstName": teacher.FirstName,
			"lastName":  teacher.LastName,
			"role":      teacher.Role,
			"status":    teacher.Status,
		})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"teachers": teacherList,
	})
}

// GetStudentsByTeacher 获取教师教授的学生列表（教师及以上权限）
func GetStudentsByTeacher(c *gin.Context) {
	teacherID := c.GetInt64("userID")
	
	relationRepo := repository.NewUserRelationRepository(database.DB)
	students, err := relationRepo.GetStudentsByTeacherID(c.Request.Context(), teacherID)
	if err != nil {
		log.Printf("获取学生列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}
	
	var studentList []gin.H
	for _, student := range students {
		studentList = append(studentList, gin.H{
			"id":        student.ID,
			"username":  student.Username,
			"email":     student.Email,
			"firstName": student.FirstName,
			"lastName":  student.LastName,
			"role":      student.Role,
			"status":    student.Status,
		})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"students": studentList,
	})
}

// GetParentsByStudent 获取学生的家长列表（学生及以上权限）
func GetParentsByStudent(c *gin.Context) {
	studentID := c.GetInt64("userID")
	
	relationRepo := repository.NewUserRelationRepository(database.DB)
	parents, err := relationRepo.GetParentsByStudentID(c.Request.Context(), studentID)
	if err != nil {
		log.Printf("获取家长列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}
	
	var parentList []gin.H
	for _, parent := range parents {
		parentList = append(parentList, gin.H{
			"id":        parent.ID,
			"username":  parent.Username,
			"email":     parent.Email,
			"firstName": parent.FirstName,
			"lastName":  parent.LastName,
			"role":      parent.Role,
			"status":    parent.Status,
		})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"parents": parentList,
	})
}
