package repository

import (
	"context"
	"errors"
	"EduGo_servers/internal/models"
	"gorm.io/gorm"
)

type UserRelationRepository interface {
	CreateRelation(ctx context.Context, relation *models.UserRelation) error
	GetRelationByID(ctx context.Context, id int64) (*models.UserRelation, error)
	GetRelationsByUserID(ctx context.Context, userID int64, relationType string) ([]*models.UserRelation, error)
	GetRelationsByRelatedUserID(ctx context.Context, relatedUserID int64, relationType string) ([]*models.UserRelation, error)
	UpdateRelation(ctx context.Context, relation *models.UserRelation) error
	DeleteRelation(ctx context.Context, id int64) error
	
	// 特定关系类型的方法
	CreateAdminTeacherRelation(ctx context.Context, relation *models.AdminTeacherRelation) error
	CreateTeacherStudentRelation(ctx context.Context, relation *models.TeacherStudentRelation) error
	CreateStudentParentRelation(ctx context.Context, relation *models.StudentParentRelation) error
	
	GetAdminTeacherRelations(ctx context.Context, adminID int64) ([]*models.AdminTeacherRelation, error)
	GetTeacherStudentRelations(ctx context.Context, teacherID int64) ([]*models.TeacherStudentRelation, error)
	GetStudentParentRelations(ctx context.Context, studentID int64) ([]*models.StudentParentRelation, error)
	
	GetTeachersByAdminID(ctx context.Context, adminID int64) ([]*models.User, error)
	GetStudentsByTeacherID(ctx context.Context, teacherID int64) ([]*models.User, error)
	GetParentsByStudentID(ctx context.Context, studentID int64) ([]*models.User, error)
	
	GetAdminsByTeacherID(ctx context.Context, teacherID int64) ([]*models.User, error)
	GetTeachersByStudentID(ctx context.Context, studentID int64) ([]*models.User, error)
	GetStudentsByParentID(ctx context.Context, parentID int64) ([]*models.User, error)
}

type userRelationRepository struct {
	db *gorm.DB
}

func NewUserRelationRepository(db *gorm.DB) UserRelationRepository {
	return &userRelationRepository{db: db}
}

func (r *userRelationRepository) CreateRelation(ctx context.Context, relation *models.UserRelation) error {
	return r.db.WithContext(ctx).Create(relation).Error
}

func (r *userRelationRepository) GetRelationByID(ctx context.Context, id int64) (*models.UserRelation, error) {
	var relation models.UserRelation
	err := r.db.WithContext(ctx).First(&relation, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &relation, err
}

func (r *userRelationRepository) GetRelationsByUserID(ctx context.Context, userID int64, relationType string) ([]*models.UserRelation, error) {
	var relations []*models.UserRelation
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if relationType != "" {
		query = query.Where("relation_type = ?", relationType)
	}
	err := query.Find(&relations).Error
	return relations, err
}

func (r *userRelationRepository) GetRelationsByRelatedUserID(ctx context.Context, relatedUserID int64, relationType string) ([]*models.UserRelation, error) {
	var relations []*models.UserRelation
	query := r.db.WithContext(ctx).Where("related_user_id = ?", relatedUserID)
	if relationType != "" {
		query = query.Where("relation_type = ?", relationType)
	}
	err := query.Find(&relations).Error
	return relations, err
}

func (r *userRelationRepository) UpdateRelation(ctx context.Context, relation *models.UserRelation) error {
	return r.db.WithContext(ctx).Save(relation).Error
}

func (r *userRelationRepository) DeleteRelation(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.UserRelation{}, id).Error
}

// 特定关系类型的方法实现

func (r *userRelationRepository) CreateAdminTeacherRelation(ctx context.Context, relation *models.AdminTeacherRelation) error {
	relation.RelationType = models.RelationAdminTeacher
	return r.db.WithContext(ctx).Create(relation).Error
}

func (r *userRelationRepository) CreateTeacherStudentRelation(ctx context.Context, relation *models.TeacherStudentRelation) error {
	relation.RelationType = models.RelationTeacherStudent
	return r.db.WithContext(ctx).Create(relation).Error
}

func (r *userRelationRepository) CreateStudentParentRelation(ctx context.Context, relation *models.StudentParentRelation) error {
	relation.RelationType = models.RelationStudentParent
	return r.db.WithContext(ctx).Create(relation).Error
}

func (r *userRelationRepository) GetAdminTeacherRelations(ctx context.Context, adminID int64) ([]*models.AdminTeacherRelation, error) {
	var relations []*models.AdminTeacherRelation
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND relation_type = ?", adminID, models.RelationAdminTeacher).
		Find(&relations).Error
	return relations, err
}

func (r *userRelationRepository) GetTeacherStudentRelations(ctx context.Context, teacherID int64) ([]*models.TeacherStudentRelation, error) {
	var relations []*models.TeacherStudentRelation
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND relation_type = ?", teacherID, models.RelationTeacherStudent).
		Find(&relations).Error
	return relations, err
}

func (r *userRelationRepository) GetStudentParentRelations(ctx context.Context, studentID int64) ([]*models.StudentParentRelation, error) {
	var relations []*models.StudentParentRelation
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND relation_type = ?", studentID, models.RelationStudentParent).
		Find(&relations).Error
	return relations, err
}

func (r *userRelationRepository) GetTeachersByAdminID(ctx context.Context, adminID int64) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).
		Joins("JOIN user_relations ON users.id = user_relations.related_user_id").
		Where("user_relations.user_id = ? AND user_relations.relation_type = ?", adminID, models.RelationAdminTeacher).
		Find(&users).Error
	return users, err
}

func (r *userRelationRepository) GetStudentsByTeacherID(ctx context.Context, teacherID int64) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).
		Joins("JOIN user_relations ON users.id = user_relations.related_user_id").
		Where("user_relations.user_id = ? AND user_relations.relation_type = ?", teacherID, models.RelationTeacherStudent).
		Find(&users).Error
	return users, err
}

func (r *userRelationRepository) GetParentsByStudentID(ctx context.Context, studentID int64) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).
		Joins("JOIN user_relations ON users.id = user_relations.related_user_id").
		Where("user_relations.user_id = ? AND user_relations.relation_type = ?", studentID, models.RelationStudentParent).
		Find(&users).Error
	return users, err
}

func (r *userRelationRepository) GetAdminsByTeacherID(ctx context.Context, teacherID int64) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).
		Joins("JOIN user_relations ON users.id = user_relations.user_id").
		Where("user_relations.related_user_id = ? AND user_relations.relation_type = ?", teacherID, models.RelationAdminTeacher).
		Find(&users).Error
	return users, err
}

func (r *userRelationRepository) GetTeachersByStudentID(ctx context.Context, studentID int64) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).
		Joins("JOIN user_relations ON users.id = user_relations.user_id").
		Where("user_relations.related_user_id = ? AND user_relations.relation_type = ?", studentID, models.RelationTeacherStudent).
		Find(&users).Error
	return users, err
}

func (r *userRelationRepository) GetStudentsByParentID(ctx context.Context, parentID int64) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).
		Joins("JOIN user_relations ON users.id = user_relations.user_id").
		Where("user_relations.related_user_id = ? AND user_relations.relation_type = ?", parentID, models.RelationStudentParent).
		Find(&users).Error
	return users, err
}
