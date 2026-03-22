package repository

import (
	"github.com/huy/quizme-backend/internal/domain"
	"gorm.io/gorm"
)

// QuestionRepository defines the interface for question data access
type QuestionRepository interface {
	Create(question *domain.Question) error
	CreateBatch(questions []domain.Question) error
	FindByID(id uint) (*domain.Question, error)
	FindByQuizID(quizID uint) ([]domain.Question, error)
	Update(question *domain.Question) error
	Delete(id uint) error
	DeleteByQuizID(quizID uint) error
	CountByQuizID(quizID uint) (int64, error)
}

type questionRepository struct {
	db *gorm.DB
}

// NewQuestionRepository creates a new question repository
func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) Create(question *domain.Question) error {
	return r.db.Create(question).Error
}

func (r *questionRepository) CreateBatch(questions []domain.Question) error {
	return r.db.Create(&questions).Error
}

func (r *questionRepository) FindByID(id uint) (*domain.Question, error) {
	var question domain.Question
	err := r.db.Preload("Options").First(&question, id).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *questionRepository) FindByQuizID(quizID uint) ([]domain.Question, error) {
	var questions []domain.Question
	err := r.db.Preload("Options").
		Where("quiz_id = ?", quizID).
		Order("order_number ASC").
		Find(&questions).Error
	return questions, err
}

func (r *questionRepository) Update(question *domain.Question) error {
	return r.db.Save(question).Error
}

func (r *questionRepository) Delete(id uint) error {
	// Delete options first
	if err := r.db.Where("question_id = ?", id).Delete(&domain.QuestionOption{}).Error; err != nil {
		return err
	}
	return r.db.Delete(&domain.Question{}, id).Error
}

func (r *questionRepository) DeleteByQuizID(quizID uint) error {
	// Get all question IDs for this quiz
	var questionIDs []uint
	r.db.Model(&domain.Question{}).Where("quiz_id = ?", quizID).Pluck("id", &questionIDs)

	// Delete options
	if len(questionIDs) > 0 {
		if err := r.db.Where("question_id IN ?", questionIDs).Delete(&domain.QuestionOption{}).Error; err != nil {
			return err
		}
	}

	// Delete questions
	return r.db.Where("quiz_id = ?", quizID).Delete(&domain.Question{}).Error
}

func (r *questionRepository) CountByQuizID(quizID uint) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Question{}).Where("quiz_id = ?", quizID).Count(&count).Error
	return count, err
}

// QuestionOptionRepository defines the interface for question option data access
type QuestionOptionRepository interface {
	Create(option *domain.QuestionOption) error
	CreateBatch(options []domain.QuestionOption) error
	FindByID(id uint) (*domain.QuestionOption, error)
	FindByQuestionID(questionID uint) ([]domain.QuestionOption, error)
	Update(option *domain.QuestionOption) error
	Delete(id uint) error
	DeleteByQuestionID(questionID uint) error
}

type questionOptionRepository struct {
	db *gorm.DB
}

// NewQuestionOptionRepository creates a new question option repository
func NewQuestionOptionRepository(db *gorm.DB) QuestionOptionRepository {
	return &questionOptionRepository{db: db}
}

func (r *questionOptionRepository) Create(option *domain.QuestionOption) error {
	return r.db.Create(option).Error
}

func (r *questionOptionRepository) CreateBatch(options []domain.QuestionOption) error {
	return r.db.Create(&options).Error
}

func (r *questionOptionRepository) FindByID(id uint) (*domain.QuestionOption, error) {
	var option domain.QuestionOption
	err := r.db.First(&option, id).Error
	if err != nil {
		return nil, err
	}
	return &option, nil
}

func (r *questionOptionRepository) FindByQuestionID(questionID uint) ([]domain.QuestionOption, error) {
	var options []domain.QuestionOption
	err := r.db.Where("question_id = ?", questionID).Find(&options).Error
	return options, err
}

func (r *questionOptionRepository) Update(option *domain.QuestionOption) error {
	return r.db.Save(option).Error
}

func (r *questionOptionRepository) Delete(id uint) error {
	return r.db.Delete(&domain.QuestionOption{}, id).Error
}

func (r *questionOptionRepository) DeleteByQuestionID(questionID uint) error {
	return r.db.Where("question_id = ?", questionID).Delete(&domain.QuestionOption{}).Error
}
