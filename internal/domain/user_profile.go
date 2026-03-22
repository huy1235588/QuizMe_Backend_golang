package domain

import (
	"time"
)

// UserProfile contains extended user profile information
type UserProfile struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint       `gorm:"column:user_id;uniqueIndex;not null" json:"userId"`
	DateOfBirth    *time.Time `gorm:"column:date_of_birth;type:date" json:"dateOfBirth,omitempty"`
	City           *string    `gorm:"size:100" json:"city,omitempty"`
	PhoneNumber    *string    `gorm:"column:phone_number;size:20" json:"phoneNumber,omitempty"`
	TotalScore     int        `gorm:"column:total_score;not null;default:0" json:"totalScore"`
	QuizzesPlayed  int        `gorm:"column:quizzes_played;not null;default:0" json:"quizzesPlayed"`
	QuizzesCreated int        `gorm:"column:quizzes_created;not null;default:0" json:"quizzesCreated"`
	TotalQuizPlays int        `gorm:"column:total_quiz_plays;not null;default:0" json:"totalQuizPlays"`

	// Association
	User *User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name for UserProfile
func (UserProfile) TableName() string {
	return "user_profile"
}
