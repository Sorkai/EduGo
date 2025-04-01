package repository

import (
	"EduGo_servers/internal/database"
	"EduGo_servers/internal/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

// AssessmentRepository 测评仓库
type AssessmentRepository struct {
	db *gorm.DB
}

// NewAssessmentRepository 创建测评仓库
func NewAssessmentRepository() *AssessmentRepository {
	return &AssessmentRepository{
		db: database.DB,
	}
}

// CreateAssessment 创建测评
func (r *AssessmentRepository) CreateAssessment(assessment *models.Assessment) error {
	return r.db.Create(assessment).Error
}

// GetAssessmentByID 根据ID获取测评
func (r *AssessmentRepository) GetAssessmentByID(id uint) (*models.Assessment, error) {
	var assessment models.Assessment
	err := r.db.Preload("Questions").Preload("Students").Preload("Students.Student").First(&assessment, id).Error
	if err != nil {
		return nil, err
	}

	// 计算总分
	assessment.CalculateTotalScore()

	return &assessment, nil
}

// GetAssessmentWithQuestionsAndStudents 获取测评及其题目和学生
func (r *AssessmentRepository) GetAssessmentWithQuestionsAndStudents(id uint) (*models.Assessment, error) {
	var assessment models.Assessment
	err := r.db.Preload("Questions").Preload("Students").Preload("Students.Student").First(&assessment, id).Error
	if err != nil {
		return nil, err
	}

	// 计算总分
	assessment.CalculateTotalScore()

	return &assessment, nil
}

// GetAssessmentsByCreatorID 获取创建者的所有测评
func (r *AssessmentRepository) GetAssessmentsByCreatorID(creatorID uint) ([]models.Assessment, error) {
	var assessments []models.Assessment
	err := r.db.Where("creator_id = ?", creatorID).Find(&assessments).Error
	if err != nil {
		return nil, err
	}

	// 计算每个测评的总分
	for i := range assessments {
		// 加载题目
		r.db.Model(&assessments[i]).Association("Questions").Find(&assessments[i].Questions)
		assessments[i].CalculateTotalScore()
	}

	return assessments, nil
}

// GetAssessmentsForStudent 获取学生可参与的测评
func (r *AssessmentRepository) GetAssessmentsForStudent(studentID uint) ([]map[string]interface{}, error) {
	// 获取分配给学生的测评
	var assessmentStudents []models.AssessmentStudent
	err := r.db.Where("student_id = ?", studentID).Find(&assessmentStudents).Error
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, as := range assessmentStudents {
		// 获取测评详情
		var assessment models.Assessment
		err := r.db.Preload("Questions").First(&assessment, as.AssessmentID).Error
		if err != nil {
			continue
		}

		// 计算总分
		assessment.CalculateTotalScore()

		// 构建结果
		assessmentInfo := map[string]interface{}{
			"id":             assessment.ID,
			"title":          assessment.Title,
			"description":    assessment.Description,
			"status":         assessment.Status,
			"start_time":     assessment.StartTime,
			"end_time":       assessment.EndTime,
			"total_score":    assessment.TotalScore,
			"student_status": as.Status,
			"student_score":  as.Score,
			"started_at":     as.StartedAt,
			"completed_at":   as.CompletedAt,
		}

		result = append(result, assessmentInfo)
	}

	return result, nil
}

// AddQuestionToAssessment 添加题目到测评
func (r *AssessmentRepository) AddQuestionToAssessment(question *models.Question) error {
	// 检查测评是否存在
	var assessment models.Assessment
	err := r.db.First(&assessment, question.AssessmentID).Error
	if err != nil {
		return err
	}

	// 检查测评状态
	if assessment.Status != "draft" {
		return errors.New("只能在草稿状态下添加题目")
	}

	// 添加题目
	return r.db.Create(question).Error
}

// UpdateAssessmentStatus 更新测评状态
func (r *AssessmentRepository) UpdateAssessmentStatus(id uint, status string) error {
	return r.db.Model(&models.Assessment{}).Where("id = ?", id).Update("status", status).Error
}

// PublishAssessment 发布测评
func (r *AssessmentRepository) PublishAssessment(id uint, startTime, endTime time.Time, studentIDs []uint) error {
	// 开始事务
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新测评状态和时间
	err := tx.Model(&models.Assessment{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     "published",
		"start_time": startTime,
		"end_time":   endTime,
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 为每个学生创建关联
	for _, studentID := range studentIDs {
		// 检查是否已存在关联
		var count int64
		tx.Model(&models.AssessmentStudent{}).Where("assessment_id = ? AND student_id = ?", id, studentID).Count(&count)
		if count > 0 {
			continue
		}

		// 创建关联
		assessmentStudent := models.AssessmentStudent{
			AssessmentID: id,
			StudentID:    studentID,
			Status:       "assigned",
		}
		err := tx.Create(&assessmentStudent).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交事务
	return tx.Commit().Error
}

// StartAssessment 学生开始测评
func (r *AssessmentRepository) StartAssessment(assessmentID, studentID uint) error {
	// 获取测评
	assessment, err := r.GetAssessmentWithQuestionsAndStudents(assessmentID)
	if err != nil {
		return err
	}

	// 检查学生是否可以开始测评
	if !assessment.CanStudentStartAssessment(studentID) {
		return errors.New("无法开始测评")
	}

	// 更新学生测评状态
	return r.db.Model(&models.AssessmentStudent{}).
		Where("assessment_id = ? AND student_id = ?", assessmentID, studentID).
		Updates(map[string]interface{}{
			"status":     "started",
			"started_at": time.Now(),
		}).Error
}

// SubmitAssessment 学生提交测评
func (r *AssessmentRepository) SubmitAssessment(assessmentID, studentID uint, answers []map[string]interface{}) (map[string]interface{}, error) {
	// 获取测评
	assessment, err := r.GetAssessmentWithQuestionsAndStudents(assessmentID)
	if err != nil {
		return nil, err
	}

	// 检查学生是否可以提交测评
	if !assessment.CanStudentSubmitAssessment(studentID) {
		return nil, errors.New("无法提交测评")
	}

	// 开始事务
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取学生测评关联
	var assessmentStudent models.AssessmentStudent
	err = tx.Where("assessment_id = ? AND student_id = ?", assessmentID, studentID).First(&assessmentStudent).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 处理答案
	totalScore := 0
	for _, answerData := range answers {
		questionID := uint(answerData["question_id"].(float64))
		answer := answerData["answer"].(string)

		// 获取题目
		var question models.Question
		err := tx.First(&question, questionID).Error
		if err != nil {
			continue
		}

		// 判断答案是否正确
		isCorrect := question.Answer == answer
		if isCorrect {
			totalScore += question.Score
		}

		// 保存学生答案
		studentAnswer := models.StudentAnswer{
			AssessmentStudentID: assessmentStudent.ID,
			QuestionID:          questionID,
			Answer:              answer,
			IsCorrect:           isCorrect,
		}
		err = tx.Create(&studentAnswer).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 更新学生测评状态和分数
	err = tx.Model(&assessmentStudent).Updates(map[string]interface{}{
		"status":       "completed",
		"score":        totalScore,
		"completed_at": time.Now(),
	}).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	// 重新获取测评和学生答案
	assessment, err = r.GetAssessmentWithQuestionsAndStudents(assessmentID)
	if err != nil {
		return nil, err
	}

	// 加载学生答案
	var studentAnswers []models.StudentAnswer
	err = r.db.Where("assessment_student_id = ?", assessmentStudent.ID).Preload("Question").Find(&studentAnswers).Error
	if err != nil {
		return nil, err
	}
	assessmentStudent.Answers = studentAnswers

	// 构建结果
	result := map[string]interface{}{
		"assessment_id": assessment.ID,
		"title":         assessment.Title,
		"total_score":   assessment.TotalScore,
		"your_score":    totalScore,
		"answers":       []map[string]interface{}{},
	}

	// 添加答案详情
	for _, answer := range assessmentStudent.Answers {
		answerDetail := map[string]interface{}{
			"question_id":    answer.QuestionID,
			"content":        answer.Question.Content,
			"your_answer":    answer.Answer,
			"correct_answer": answer.Question.Answer,
			"is_correct":     answer.IsCorrect,
			"score":          answer.Question.Score,
		}
		if answer.Question.Explanation != "" {
			answerDetail["explanation"] = answer.Question.Explanation
		}
		result["answers"] = append(result["answers"].([]map[string]interface{}), answerDetail)
	}

	// 生成AI分析
	result["ai_analysis"] = generateAIAnalysis(&assessmentStudent)

	return result, nil
}

// GetAssessmentResult 获取测评结果
func (r *AssessmentRepository) GetAssessmentResult(assessmentID, studentID uint) (map[string]interface{}, error) {
	// 检查测评是否存在
	_, err := r.GetAssessmentWithQuestionsAndStudents(assessmentID)
	if err != nil {
		return nil, err
	}

	// 获取学生测评关联
	var assessmentStudent models.AssessmentStudent
	err = r.db.Where("assessment_id = ? AND student_id = ?", assessmentID, studentID).First(&assessmentStudent).Error
	if err != nil {
		return nil, err
	}

	// 检查学生是否已完成测评
	if assessmentStudent.Status != "completed" {
		return nil, errors.New("学生尚未完成测评")
	}

	// 加载学生答案
	var studentAnswers []models.StudentAnswer
	err = r.db.Where("assessment_student_id = ?", assessmentStudent.ID).Preload("Question").Find(&studentAnswers).Error
	if err != nil {
		return nil, err
	}
	assessmentStudent.Answers = studentAnswers

	// 构建结果
	result := map[string]interface{}{
		"your_score":   assessmentStudent.Score,
		"completed_at": assessmentStudent.CompletedAt,
		"answers":      []map[string]interface{}{},
	}

	// 添加答案详情
	for _, answer := range assessmentStudent.Answers {
		answerDetail := map[string]interface{}{
			"question_id":    answer.QuestionID,
			"content":        answer.Question.Content,
			"options":        answer.Question.Options,
			"your_answer":    answer.Answer,
			"correct_answer": answer.Question.Answer,
			"is_correct":     answer.IsCorrect,
			"score":          answer.Question.Score,
		}
		if answer.Question.Explanation != "" {
			answerDetail["explanation"] = answer.Question.Explanation
		}
		result["answers"] = append(result["answers"].([]map[string]interface{}), answerDetail)
	}

	// 生成AI分析
	result["ai_analysis"] = generateAIAnalysis(&assessmentStudent)

	return result, nil
}

// GetAssessmentStudents 获取测评的学生列表
func (r *AssessmentRepository) GetAssessmentStudents(assessmentID uint) ([]map[string]interface{}, error) {
	// 检查测评是否存在
	_, err := r.GetAssessmentByID(assessmentID)
	if err != nil {
		return nil, err
	}

	// 获取学生测评关联
	var assessmentStudents []models.AssessmentStudent
	err = r.db.Where("assessment_id = ?", assessmentID).Preload("Student").Find(&assessmentStudents).Error
	if err != nil {
		return nil, err
	}

	// 构建结果
	var result []map[string]interface{}
	for _, as := range assessmentStudents {
		studentInfo := map[string]interface{}{
			"id":           as.Student.ID,
			"username":     as.Student.Username,
			"email":        as.Student.Email,
			"firstName":    as.Student.FirstName,
			"lastName":     as.Student.LastName,
			"status":       as.Status,
			"score":        as.Score,
			"started_at":   as.StartedAt,
			"completed_at": as.CompletedAt,
		}
		result = append(result, studentInfo)
	}

	return result, nil
}

// generateAIAnalysis 生成AI分析
func generateAIAnalysis(as *models.AssessmentStudent) string {
	// 计算正确率
	correctCount := 0
	for _, answer := range as.Answers {
		if answer.IsCorrect {
			correctCount++
		}
	}
	correctRate := float64(correctCount) / float64(len(as.Answers))

	// 根据正确率生成分析
	var analysis string
	if correctRate >= 0.9 {
		analysis = "您的表现非常出色！您已经掌握了大部分知识点，建议您可以尝试更具挑战性的内容。"
	} else if correctRate >= 0.7 {
		analysis = "您的表现良好，对大部分知识点有很好的理解。建议您重点关注错误题目，巩固相关知识点。"
	} else if correctRate >= 0.5 {
		analysis = "您的表现一般，对部分知识点理解不够深入。建议您重新学习相关内容，并多做练习。"
	} else {
		analysis = "您的表现需要提升，对大部分知识点理解不够。建议您从基础知识开始，系统性地学习相关内容。"
	}

	// 添加错题分析
	if correctRate < 1.0 {
		analysis += "\n\n以下是您需要重点关注的知识点："
		for _, answer := range as.Answers {
			if !answer.IsCorrect {
				analysis += "\n- " + answer.Question.Content
			}
		}
	}

	return analysis
}
