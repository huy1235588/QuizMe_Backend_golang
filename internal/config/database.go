package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/huy/quizme-backend/internal/domain"
)

// InitDatabase initializes the database connection and runs migrations
func InitDatabase(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	// Configure GORM logger
	gormLogger := logger.Default
	if cfg.Server.Mode == "debug" {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Run auto migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database connected successfully")
	return db, nil
}

// runMigrations runs GORM auto migrations for all domain models
func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		// Phase 1 entities
		&domain.User{},
		&domain.UserProfile{},
		&domain.RefreshToken{},
		// Phase 2 entities
		&domain.Category{},
		&domain.Quiz{},
		&domain.Question{},
		&domain.QuestionOption{},
		&domain.Room{},
		&domain.RoomParticipant{},
		&domain.RoomChat{},
		// Phase 3 entities
		&domain.GameResult{},
		&domain.GameResultQuestion{},
		&domain.GamePlayerAnswer{},
		&domain.GamePlayerAnswerOption{},
	)
}
