package models

import (
	"time"

	"gorm.io/gorm"
)

// Assessment 测评模型
type Assessment struct {
	gorm.Model
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status" gorm:"default:draft"` // draft, published, closed
	CreatorID   uint      `json:"creator_id"`
	Creator     User      `json:"creator" gorm:"foreignKey:CreatorID"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	TotalScore  int       `json:"total_score" gorm:"-"` // 不存储在数据库中，由计算得出
	Questions   []Question
	Students    []AssessmentStudent `gorm:"foreignKey:AssessmentID"`
}

// Question 题目模型
type Question struct {
	gorm.Model
	AssessmentID uint     `json:"assessment_id"`
	Content      string   `json:"content"`
	Type         string   `json:"type" gorm:"default:single_choice"` // single_choice, multiple_choice, etc.
	Options      []string `json:"options" gorm:"type:json"`
	Answer       string   `json:"answer"`
	Score        int      `json:"score"`
	Explanation  string   `json:"explanation"`
}

// AssessmentStudent 测评-学生关联模型
type AssessmentStudent struct {
	gorm.Model
	AssessmentID uint      `json:"assessment_id"`
	StudentID    uint      `json:"student_id"`
	Student      User      `json:"student" gorm:"foreignKey:StudentID"`
	Status       string    `json:"status" gorm:"default:assigned"` // assigned, started, completed
	Score        int       `json:"score" gorm:"default:0"`
	StartedAt    time.Time `json:"started_at"`
	CompletedAt  time.Time `json:"completed_at"`
	Answers      []StudentAnswer `gorm:"-"` // 不存储在数据库中，通过关联查询获取
}

// StudentAnswer 学生答案模型
type StudentAnswer struct {
	gorm.Model
	AssessmentStudentID uint     `json:"assessment_student_id"`
	QuestionID          uint     `json:"question_id"`
	Question            Question `json:"question" gorm:"foreignKey:QuestionID"`
	Answer              string   `json:"answer"`
	IsCorrect           bool     `json:"is_correct"`
}

// CalculateTotalScore 计算测评总分
func (a *Assessment) CalculateTotalScore() {
	totalScore := 0
	for _, question := range a.Questions {
		totalScore += question.Score
	}
	a.TotalScore = totalScore
}

// CanStudentStartAssessment 检查学生是否可以开始测评
func (a *Assessment) CanStudentStartAssessment(studentID uint) bool {
	// 检查测评状态
	if a.Status != "published" {
		return false
	}

	// 检查测评时间
	now := time.Now()
	if now.Before(a.StartTime) || now.After(a.EndTime) {
		return false
	}

	// 检查学生是否被分配
	for _, as := range a.Students {
		if as.StudentID == studentID {
			// 检查学生测评状态
			return as.Status == "assigned"
		}
	}

	return false
}

// CanStudentSubmitAssessment 检查学生是否可以提交测评
func (a *Assessment) CanStudentSubmitAssessment(studentID uint) bool {
	// 检查测评状态
	if a.Status != "published" {
		return false
	}

	// 检查测评时间
	now := time.Now()
	if now.After(a.EndTime) {
		return false
	}

	// 检查学生是否被分配
	for _, as := range a.Students {
		if as.StudentID == studentID {
			// 检查学生测评状态
			return as.Status == "started"
		}
	}

	return false
}
