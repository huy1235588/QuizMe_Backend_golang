package domain

import (
	"time"

	"github.com/huy/quizme-backend/internal/domain/enums"
)

// Room represents a game room
type Room struct {
	ID         uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string           `gorm:"size:100;not null" json:"name"`
	Code       string           `gorm:"uniqueIndex;size:10;not null" json:"code"`
	QuizID     uint             `gorm:"column:quiz_id;not null" json:"quizId"`
	HostID     uint             `gorm:"column:host_id;not null" json:"hostId"`
	Password   *string          `gorm:"size:255" json:"-"`
	IsPublic   bool             `gorm:"column:is_public;default:true" json:"isPublic"`
	MaxPlayers int              `gorm:"column:max_players;default:10" json:"maxPlayers"`
	Status     enums.RoomStatus `gorm:"column:status;type:varchar(20);default:'WAITING'" json:"status"`
	StartTime  *time.Time       `gorm:"column:start_time" json:"startTime,omitempty"`
	EndTime    *time.Time       `gorm:"column:end_time" json:"endTime,omitempty"`
	CreatedAt  time.Time        `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time        `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// Associations
	Quiz         *Quiz             `gorm:"foreignKey:QuizID" json:"quiz,omitempty"`
	Host         *User             `gorm:"foreignKey:HostID" json:"host,omitempty"`
	Participants []RoomParticipant `gorm:"foreignKey:RoomID" json:"participants,omitempty"`
}

// TableName specifies the table name for Room
func (Room) TableName() string {
	return "room"
}

// HasPassword checks if the room has a password
func (r *Room) HasPassword() bool {
	return r.Password != nil && *r.Password != ""
}

// RoomParticipant represents a participant in a room
type RoomParticipant struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	RoomID    uint       `gorm:"column:room_id;not null;index" json:"roomId"`
	UserID    *uint      `gorm:"column:user_id;index" json:"userId,omitempty"`
	Score     int        `gorm:"default:0" json:"score"`
	IsHost    bool       `gorm:"column:is_host;default:false" json:"isHost"`
	JoinedAt  time.Time  `gorm:"column:joined_at;autoCreateTime" json:"joinedAt"`
	LeftAt    *time.Time `gorm:"column:left_at" json:"leftAt,omitempty"`
	IsGuest   bool       `gorm:"column:is_guest;default:false" json:"isGuest"`
	GuestName *string    `gorm:"column:guest_name;size:50" json:"guestName,omitempty"`

	// Associations
	Room *Room `gorm:"foreignKey:RoomID" json:"-"`
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for RoomParticipant
func (RoomParticipant) TableName() string {
	return "room_participant"
}

// GetDisplayName returns the display name for the participant
func (rp *RoomParticipant) GetDisplayName() string {
	if rp.IsGuest && rp.GuestName != nil {
		return *rp.GuestName
	}
	if rp.User != nil {
		return rp.User.Username
	}
	return "Unknown"
}

// RoomChat represents a chat message in a room
type RoomChat struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RoomID    uint      `gorm:"column:room_id;not null;index" json:"roomId"`
	UserID    *uint     `gorm:"column:user_id;index" json:"userId,omitempty"`
	IsGuest   bool      `gorm:"column:is_guest;default:false" json:"isGuest"`
	GuestName *string   `gorm:"column:guest_name;size:50" json:"guestName,omitempty"`
	Message   string    `gorm:"column:message;type:text;not null" json:"message"`
	SentAt    time.Time `gorm:"column:sent_at;autoCreateTime" json:"sentAt"`

	// Associations
	Room *Room `gorm:"foreignKey:RoomID" json:"-"`
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for RoomChat
func (RoomChat) TableName() string {
	return "room_chat"
}
