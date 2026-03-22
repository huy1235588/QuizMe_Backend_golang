package repository

import (
	"github.com/huy/quizme-backend/internal/domain"
	"github.com/huy/quizme-backend/internal/domain/enums"
	"gorm.io/gorm"
)

// QuizRepository defines the interface for quiz data access
type QuizRepository interface {
	Create(quiz *domain.Quiz) error
	FindByID(id uint) (*domain.Quiz, error)
	FindByIDWithQuestions(id uint) (*domain.Quiz, error)
	FindAll() ([]domain.Quiz, error)
	FindPublic() ([]domain.Quiz, error)
	FindByDifficulty(difficulty enums.Difficulty) ([]domain.Quiz, error)
	FindByCreatorID(creatorID uint) ([]domain.Quiz, error)
	FindWithFilters(categoryID *uint, difficulty *string, isPublic *bool, search *string, page, pageSize int, sortBy, sortDir string) ([]domain.Quiz, int64, error)
	Update(quiz *domain.Quiz) error
	Delete(id uint) error
	IncrementPlayCount(id uint) error
	UpdateQuestionCount(id uint, count int) error
	AddCategories(quizID uint, categoryIDs []uint) error
	RemoveCategories(quizID uint) error
}

type quizRepository struct {
	db *gorm.DB
}

// NewQuizRepository creates a new quiz repository
func NewQuizRepository(db *gorm.DB) QuizRepository {
	return &quizRepository{db: db}
}

func (r *quizRepository) Create(quiz *domain.Quiz) error {
	return r.db.Create(quiz).Error
}

func (r *quizRepository) FindByID(id uint) (*domain.Quiz, error) {
	var quiz domain.Quiz
	err := r.db.Preload("Creator").Preload("Categories").First(&quiz, id).Error
	if err != nil {
		return nil, err
	}
	return &quiz, nil
}

func (r *quizRepository) FindByIDWithQuestions(id uint) (*domain.Quiz, error) {
	var quiz domain.Quiz
	err := r.db.Preload("Creator").
		Preload("Categories").
		Preload("Questions", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_number ASC")
		}).
		Preload("Questions.Options").
		First(&quiz, id).Error
	if err != nil {
		return nil, err
	}
	return &quiz, nil
}

func (r *quizRepository) FindAll() ([]domain.Quiz, error) {
	var quizzes []domain.Quiz
	err := r.db.Preload("Creator").Preload("Categories").
		Order("created_at DESC").Find(&quizzes).Error
	return quizzes, err
}

func (r *quizRepository) FindPublic() ([]domain.Quiz, error) {
	var quizzes []domain.Quiz
	err := r.db.Preload("Creator").Preload("Categories").
		Where("is_public = ?", true).
		Order("created_at DESC").Find(&quizzes).Error
	return quizzes, err
}

func (r *quizRepository) FindByDifficulty(difficulty enums.Difficulty) ([]domain.Quiz, error) {
	var quizzes []domain.Quiz
	err := r.db.Preload("Creator").Preload("Categories").
		Where("difficulty = ? AND is_public = ?", difficulty, true).
		Order("created_at DESC").Find(&quizzes).Error
	return quizzes, err
}

func (r *quizRepository) FindByCreatorID(creatorID uint) ([]domain.Quiz, error) {
	var quizzes []domain.Quiz
	err := r.db.Preload("Creator").Preload("Categories").
		Where("creator_id = ?", creatorID).
		Order("created_at DESC").Find(&quizzes).Error
	return quizzes, err
}

func (r *quizRepository) FindWithFilters(
	categoryID *uint,
	difficulty *string,
	isPublic *bool,
	search *string,
	page, pageSize int,
	sortBy, sortDir string,
) ([]domain.Quiz, int64, error) {
	var quizzes []domain.Quiz
	var total int64

	query := r.db.Model(&domain.Quiz{})

	// Apply category filter (join with quiz_category)
	if categoryID != nil {
		query = query.Joins("JOIN quiz_category ON quiz.id = quiz_category.quiz_id").
			Where("quiz_category.category_id = ?", *categoryID)
	}

	// Apply difficulty filter
	if difficulty != nil && *difficulty != "" {
		query = query.Where("quiz.difficulty = ?", *difficulty)
	}

	// Apply public filter
	if isPublic != nil {
		query = query.Where("quiz.is_public = ?", *isPublic)
	}

	// Apply search filter
	if search != nil && *search != "" {
		searchPattern := "%" + *search + "%"
		query = query.Where("LOWER(quiz.title) LIKE LOWER(?) OR LOWER(quiz.description) LIKE LOWER(?)",
			searchPattern, searchPattern)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if sortBy != "" {
		order := "quiz." + sortBy
		if sortDir == "desc" {
			order += " DESC"
		} else {
			order += " ASC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("quiz.created_at DESC")
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Preload("Creator").Preload("Categories").
		Offset(offset).Limit(pageSize).Find(&quizzes).Error

	return quizzes, total, err
}

func (r *quizRepository) Update(quiz *domain.Quiz) error {
	return r.db.Save(quiz).Error
}

func (r *quizRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Quiz{}, id).Error
}

func (r *quizRepository) IncrementPlayCount(id uint) error {
	return r.db.Model(&domain.Quiz{}).Where("id = ?", id).
		Update("play_count", gorm.Expr("play_count + 1")).Error
}

func (r *quizRepository) UpdateQuestionCount(id uint, count int) error {
	return r.db.Model(&domain.Quiz{}).Where("id = ?", id).
		Update("question_count", count).Error
}

func (r *quizRepository) AddCategories(quizID uint, categoryIDs []uint) error {
	var quiz domain.Quiz
	if err := r.db.First(&quiz, quizID).Error; err != nil {
		return err
	}

	var categories []domain.Category
	if err := r.db.Where("id IN ?", categoryIDs).Find(&categories).Error; err != nil {
		return err
	}

	return r.db.Model(&quiz).Association("Categories").Replace(categories)
}

func (r *quizRepository) RemoveCategories(quizID uint) error {
	var quiz domain.Quiz
	if err := r.db.First(&quiz, quizID).Error; err != nil {
		return err
	}
	return r.db.Model(&quiz).Association("Categories").Clear()
}
