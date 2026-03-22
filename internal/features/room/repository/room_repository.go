package repository

import (
	"github.com/huy/quizme-backend/internal/domain/enums"
	"github.com/huy/quizme-backend/internal/features/room/domain"
	"gorm.io/gorm"
)

type roomRepository struct {
	db *gorm.DB
}

// NewRoomRepository creates a new room repository
func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(room *domain.Room) error {
	return r.db.Create(room).Error
}

func (r *roomRepository) FindByID(id uint) (*domain.Room, error) {
	var room domain.Room
	err := r.db.Preload("Quiz").Preload("Host").Preload("Participants").
		Preload("Participants.User").First(&room, id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) FindByCode(code string) (*domain.Room, error) {
	var room domain.Room
	err := r.db.Preload("Quiz").Preload("Host").Preload("Participants").
		Preload("Participants.User").Where("code = ?", code).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) FindByHostID(hostID uint) ([]domain.Room, error) {
	var rooms []domain.Room
	err := r.db.Preload("Quiz").Preload("Host").
		Where("host_id = ?", hostID).
		Order("created_at DESC").Find(&rooms).Error
	return rooms, err
}

func (r *roomRepository) FindWaiting() ([]domain.Room, error) {
	var rooms []domain.Room
	err := r.db.Preload("Quiz").Preload("Host").Preload("Participants").
		Where("status = ? AND is_public = ?", enums.RoomStatusWaiting, true).
		Order("created_at DESC").Find(&rooms).Error
	return rooms, err
}

func (r *roomRepository) FindAvailable(search *string, quizID *uint, page, pageSize int) ([]domain.Room, int64, error) {
	var rooms []domain.Room
	var total int64

	query := r.db.Model(&domain.Room{}).
		Where("status = ? AND is_public = ?", enums.RoomStatusWaiting, true)

	if search != nil && *search != "" {
		searchPattern := "%" + *search + "%"
		query = query.Where("LOWER(name) LIKE LOWER(?)", searchPattern)
	}

	if quizID != nil {
		query = query.Where("quiz_id = ?", *quizID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Quiz").Preload("Host").Preload("Participants").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).Find(&rooms).Error

	return rooms, total, err
}

func (r *roomRepository) Update(room *domain.Room) error {
	return r.db.Save(room).Error
}

func (r *roomRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Room{}, id).Error
}

func (r *roomRepository) UpdateStatus(id uint, status enums.RoomStatus) error {
	return r.db.Model(&domain.Room{}).Where("id = ?", id).Update("status", status).Error
}

type roomParticipantRepository struct {
	db *gorm.DB
}

// NewRoomParticipantRepository creates a new room participant repository
func NewRoomParticipantRepository(db *gorm.DB) RoomParticipantRepository {
	return &roomParticipantRepository{db: db}
}

func (r *roomParticipantRepository) Create(participant *domain.RoomParticipant) error {
	return r.db.Create(participant).Error
}

func (r *roomParticipantRepository) FindByID(id uint) (*domain.RoomParticipant, error) {
	var participant domain.RoomParticipant
	err := r.db.Preload("User").First(&participant, id).Error
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

func (r *roomParticipantRepository) FindByRoomID(roomID uint) ([]domain.RoomParticipant, error) {
	var participants []domain.RoomParticipant
	err := r.db.Preload("User").Where("room_id = ?", roomID).
		Order("score DESC").Find(&participants).Error
	return participants, err
}

func (r *roomParticipantRepository) FindByRoomAndUser(roomID, userID uint) (*domain.RoomParticipant, error) {
	var participant domain.RoomParticipant
	err := r.db.Where("room_id = ? AND user_id = ?", roomID, userID).First(&participant).Error
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

func (r *roomParticipantRepository) Update(participant *domain.RoomParticipant) error {
	return r.db.Save(participant).Error
}

func (r *roomParticipantRepository) Delete(id uint) error {
	return r.db.Delete(&domain.RoomParticipant{}, id).Error
}

func (r *roomParticipantRepository) DeleteByRoomID(roomID uint) error {
	return r.db.Where("room_id = ?", roomID).Delete(&domain.RoomParticipant{}).Error
}

func (r *roomParticipantRepository) CountByRoomID(roomID uint) (int64, error) {
	var count int64
	err := r.db.Model(&domain.RoomParticipant{}).Where("room_id = ? AND left_at IS NULL", roomID).Count(&count).Error
	return count, err
}

func (r *roomParticipantRepository) UpdateScore(id uint, score int) error {
	return r.db.Model(&domain.RoomParticipant{}).Where("id = ?", id).Update("score", score).Error
}

type roomChatRepository struct {
	db *gorm.DB
}

// NewRoomChatRepository creates a new room chat repository
func NewRoomChatRepository(db *gorm.DB) RoomChatRepository {
	return &roomChatRepository{db: db}
}

func (r *roomChatRepository) Create(chat *domain.RoomChat) error {
	return r.db.Create(chat).Error
}

func (r *roomChatRepository) FindByRoomID(roomID uint, limit int) ([]domain.RoomChat, error) {
	var chats []domain.RoomChat
	query := r.db.Preload("User").Where("room_id = ?", roomID).Order("sent_at ASC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&chats).Error
	return chats, err
}

func (r *roomChatRepository) DeleteByRoomID(roomID uint) error {
	return r.db.Where("room_id = ?", roomID).Delete(&domain.RoomChat{}).Error
}
