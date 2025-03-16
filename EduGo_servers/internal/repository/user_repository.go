package repository

import (
	"context"
	"errors"
	"EduGo_servers/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id int64) error
	UserExists(username string, email string) bool
	IsFirstUser() bool
	GetUsersByRole(ctx context.Context, role string) ([]*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	
	// 用户关系相关
	CreateUserRelation(ctx context.Context, relation *models.UserRelation) error
	GetUserRelations(ctx context.Context, userID int64, relationType string) ([]*models.UserRelation, error)
	GetRelatedUsers(ctx context.Context, userID int64, relationType string) ([]*models.User, error)
	UpdateUserRelation(ctx context.Context, relation *models.UserRelation) error
	DeleteUserRelation(ctx context.Context, id int64) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) DeleteUser(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (r *userRepository) UserExists(username string, email string) bool {
	var count int64
	r.db.Model(&models.User{}).
		Where("username = ? OR email = ?", username, email).
		Count(&count)
	return count > 0
}

// IsFirstUser 检查是否是第一个用户
func (r *userRepository) IsFirstUser() bool {
	var count int64
	r.db.Model(&models.User{}).Count(&count)
	return count == 0
}

// GetUsersByRole 根据角色获取用户列表
func (r *userRepository) GetUsersByRole(ctx context.Context, role string) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).Where("role = ?", role).Find(&users).Error
	return users, err
}

// GetAllUsers 获取所有用户
func (r *userRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}

// CreateUserRelation 创建用户关系
func (r *userRepository) CreateUserRelation(ctx context.Context, relation *models.UserRelation) error {
	return r.db.WithContext(ctx).Create(relation).Error
}

// GetUserRelations 获取用户的关系列表
func (r *userRepository) GetUserRelations(ctx context.Context, userID int64, relationType string) ([]*models.UserRelation, error) {
	var relations []*models.UserRelation
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if relationType != "" {
		query = query.Where("relation_type = ?", relationType)
	}
	err := query.Find(&relations).Error
	return relations, err
}

// GetRelatedUsers 获取与用户有关系的用户列表
func (r *userRepository) GetRelatedUsers(ctx context.Context, userID int64, relationType string) ([]*models.User, error) {
	var users []*models.User
	query := r.db.WithContext(ctx).
		Joins("JOIN user_relations ON users.id = user_relations.related_user_id").
		Where("user_relations.user_id = ?", userID)
	
	if relationType != "" {
		query = query.Where("user_relations.relation_type = ?", relationType)
	}
	
	err := query.Find(&users).Error
	return users, err
}

// UpdateUserRelation 更新用户关系
func (r *userRepository) UpdateUserRelation(ctx context.Context, relation *models.UserRelation) error {
	return r.db.WithContext(ctx).Save(relation).Error
}

// DeleteUserRelation 删除用户关系
func (r *userRepository) DeleteUserRelation(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.UserRelation{}, id).Error
}
