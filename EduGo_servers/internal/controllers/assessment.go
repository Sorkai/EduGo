package controllers

import (
	"EduGo_servers/internal/models"
	"EduGo_servers/internal/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// AssessmentController 测评控制器
type AssessmentController struct {
	assessmentRepo *repository.AssessmentRepository
}

// NewAssessmentController 创建测评控制器
func NewAssessmentController() *AssessmentController {
	return &AssessmentController{
		assessmentRepo: repository.NewAssessmentRepository(),
	}
}

// CreateAssessment 创建测评
func (c *AssessmentController) CreateAssessment(ctx *gin.Context) {
	// 获取当前用户
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 解析请求
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建测评
	assessment := models.Assessment{
		Title:       req.Title,
		Description: req.Description,
		Status:      "draft",
		CreatorID:   userID.(uint),
	}
	err := c.assessmentRepo.CreateAssessment(&assessment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建测评失败"})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusCreated, gin.H{
		"message":    "测评创建成功",
		"assessment": assessment,
	})
}

// GetTeacherAssessments 获取教师创建的所有测评
func (c *AssessmentController) GetTeacherAssessments(ctx *gin.Context) {
	// 获取当前用户
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取测评列表
	assessments, err := c.assessmentRepo.GetAssessmentsByCreatorID(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取测评列表失败"})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"assessments": assessments,
	})
}

// GetTeacherAssessmentDetail 获取测评详情（教师视图）
func (c *AssessmentController) GetTeacherAssessmentDetail(ctx *gin.Context) {
	// 获取当前用户
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取测评ID
	assessmentID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的测评ID"})
		return
	}

	// 获取测评详情
	assessment, err := c.assessmentRepo.GetAssessmentWithQuestionsAndStudents(uint(assessmentID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "测评不存在"})
		return
	}

	// 检查权限
	if assessment.CreatorID != userID.(uint) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权访问此测评"})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"assessment": assessment,
		"questions":  assessment.Questions,
	})
}

// AddQuestionToAssessment 添加题目到测评
func (c *AssessmentController) AddQuestionToAssessment(ctx *gin.Context) {
	// 获取当前用户
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取测评ID
	assessmentID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的测评ID"})
		return
	}

	// 获取测评详情
	assessment, err := c.assessmentRepo.GetAssessmentByID(uint(assessmentID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "测评不存在"})
		return
	}

	// 检查权限
	if assessment.CreatorID != userID.(uint) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权访问此测评"})
		return
	}

	// 检查测评状态
	if assessment.Status != "draft" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "只能在草稿状态下添加题目"})
		return
	}

	// 解析请求
	var req struct {
		Content     string   `json:"content" binding:"required"`
		Type        string   `json:"type" binding:"required"`
		Options     []string `json:"options" binding:"required"`
		Answer      string   `json:"answer" binding:"required"`
		Score       int      `json:"score" binding:"required"`
		Explanation string   `json:"explanation"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建题目
	question := models.Question{
		AssessmentID: uint(assessmentID),
		Content:      req.Content,
		Type:         req.Type,
		Options:      req.Options,
		Answer:       req.Answer,
		Score:        req.Score,
		Explanation:  req.Explanation,
	}
	err = c.assessmentRepo.AddQuestionToAssessment(&question)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "添加题目失败: " + err.Error()})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusCreated, gin.H{
		"message":  "题目添加成功",
		"question": question,
	})
}

// PublishAssessment 发布测评
func (c *AssessmentController) PublishAssessment(ctx *gin.Context) {
	// 获取当前用户
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取测评ID
	assessmentID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的测评ID"})
		return
	}

	// 获取测评详情
	assessment, err := c.assessmentRepo.GetAssessmentByID(uint(assessmentID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "测评不存在"})
		return
	}

	// 检查权限
	if assessment.CreatorID != userID.(uint) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权访问此测评"})
		return
	}

	// 检查测评状态
	if assessment.Status != "draft" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "只能发布草稿状态的测评"})
		return
	}

	// 检查是否有题目
	if len(assessment.Questions) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "测评必须包含至少一个题目"})
		return
	}

	// 解析请求
	var req struct {
		StartTime  string `json:"start_time" binding:"required"`
		EndTime    string `json:"end_time" binding:"required"`
		StudentIDs []uint `json:"student_ids" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 解析时间
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的开始时间格式"})
		return
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的结束时间格式"})
		return
	}

	// 检查时间范围
	if startTime.After(endTime) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "开始时间不能晚于结束时间"})
		return
	}

	// 发布测评
	err = c.assessmentRepo.PublishAssessment(uint(assessmentID), startTime, endTime, req.StudentIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "发布测评失败: " + err.Error()})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"message": "测评发布成功",
		"assessment": gin.H{
			"id":         assessmentID,
			"status":     "published",
			"start_time": startTime,
			"end_time":   endTime,
		},
	})
}

// CloseAssessment 关闭测评
func (c *AssessmentController) CloseAssessment(ctx *gin.Context) {
	// 获取当前用户
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取测评ID
	assessmentID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的测评ID"})
		return
	}

	// 获取测评详情
	assessment, err := c.assessmentRepo.GetAssessmentByID(uint(assessmentID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "测评不存在"})
		return
	}

	// 检查权限
	if assessment.CreatorID != userID.(uint) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权访问此测评"})
		return
	}

	// 检查测评状态
	if assessment.Status != "published" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "只能关闭已发布的测评"})
		return
	}

	// 关闭测评
	err = c.assessmentRepo.UpdateAssessmentStatus(uint(assessmentID), "closed")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "关闭测评失败"})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"message": "测评已关闭",
		"assessment": gin.H{
			"id":     assessmentID,
			"status": "closed",
		},
	})
}

// GetAssessmentStudents 获取测评学生列表
func (c *AssessmentController) GetAssessmentStudents(ctx *gin.Context) {
	// 获取当前用户
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取测评ID
	assessmentID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的测评ID"})
		return
	}

	// 获取测评详情
	assessment, err := c.assessmentRepo.GetAssessmentByID(uint(assessmentID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "测评不存在"})
		return
	}

	// 检查权限
	if assessment.CreatorID != userID.(uint) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权访问此测评"})
		return
	}

	// 获取学生列表
	students, err := c.assessmentRepo.GetAssessmentStudents(uint(assessmentID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取学生列表失败"})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"assessment": gin.H{
			"id":          assessment.ID,
			"title":       assessment.Title,
			"total_score": assessment.TotalScore,
		},
		"students": students,
	})
}

// GetStudentAssessments 获取学生可参与的测评列表
func (c *AssessmentController) GetStudentAssessments(ctx *gin.Context) {
	// 获取当前用户
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取测评列表
	assessments, err := c.assessmentRepo.GetAssessmentsForStudent(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取测评列表失败"})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"assessments": assessments,
	})
}

// GetStudentAssessmentDetail 获取测评详情（学生视图）
func (c *AssessmentController) GetStudentAssessmentDetail(ctx *gin.Context) {
	// 获取当前用户
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取测评ID
	assessmentID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的测评ID"})
		return
	}

	// 获取测评详情
	assessment, err := c.assessmentRepo.GetAssessmentWithQuestionsAndStudents(uint(assessmentID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "测评不存在"})
		return
	}

	// 检查学生是否被分配
	found := false
	for _, as := range assessment.Students {
		if as.StudentID == userID.(uint) {
			found = true
			break
		}
	}
	if !found {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权访问此测评"})
		return
	}

	// 准备学生视图的题目（不包含答案）
	studentQuestions := make([]map[string]interface{}, len(assessment.Questions))
	for i, q := range assessment.Questions {
		studentQuestions[i] = map[string]interface{}{
			"id":      q.ID,
			"content": q.Content,
			"type":    q.Type,
			"options": q.Options,
			"score":   q.Score,
		}
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"assessment": assessment,
		"questions":  studentQuestions,
	})
}

// StartAssessment 开始测评
func (c *AssessmentController) StartAssessment(ctx *gin.Context) {
	// 获取当前用户
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取测评ID
	assessmentID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的测评ID"})
		return
	}

	// 开始测评
	err = c.assessmentRepo.StartAssessment(uint(assessmentID), userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取测评详情
	assessment, err := c.assessmentRepo.GetAssessmentByID(uint(assessmentID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取测评详情失败"})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"message": "测评已开始",
		"assessment": gin.H{
			"id":        assessment.ID,
			"title":     assessment.Title,
			"end_time":  assessment.EndTime,
			"total_score": assessment.TotalScore,
		},
	})
}

// SubmitAssessment 提交测评答案
func (c *AssessmentController) SubmitAssessment(ctx *gin.Context) {
	// 获取当前用户
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取测评ID
	assessmentID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的测评ID"})
		return
	}

	// 解析请求
	var req struct {
		Answers []map[string]interface{} `json:"answers" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 提交答案
	result, err := c.assessmentRepo.SubmitAssessment(uint(assessmentID), userID.(uint), req.Answers)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"message": "答案提交成功",
		"result":  result,
	})
}

// GetAssessmentResult 获取测评结果
func (c *AssessmentController) GetAssessmentResult(ctx *gin.Context) {
	// 获取当前用户
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取测评ID
	assessmentID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的测评ID"})
		return
	}

	// 获取测评详情
	assessment, err := c.assessmentRepo.GetAssessmentByID(uint(assessmentID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "测评不存在"})
		return
	}

	// 获取测评结果
	result, err := c.assessmentRepo.GetAssessmentResult(uint(assessmentID), userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"assessment": gin.H{
			"id":          assessment.ID,
			"title":       assessment.Title,
			"description": assessment.Description,
			"total_score": assessment.TotalScore,
		},
		"result": result,
	})
}
