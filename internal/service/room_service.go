package service

import (
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"github.com/huy/quizme-backend/internal/domain"
	"github.com/huy/quizme-backend/internal/domain/enums"
	"github.com/huy/quizme-backend/internal/dto/request"
	"github.com/huy/quizme-backend/internal/dto/response"
	"github.com/huy/quizme-backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrRoomNotFound       = errors.New("room not found")
	ErrRoomFull           = errors.New("room is full")
	ErrWrongPassword      = errors.New("wrong room password")
	ErrNotRoomHost        = errors.New("you are not the host of this room")
	ErrRoomNotWaiting     = errors.New("room is not in waiting status")
	ErrAlreadyInRoom      = errors.New("you are already in this room")
	ErrNotInRoom          = errors.New("you are not in this room")
)

// RoomService handles room-related operations
type RoomService interface {
	CreateRoom(hostID uint, req *request.RoomRequest) (*response.RoomResponse, error)
	GetRoomByCode(code string) (*response.RoomResponse, error)
	GetRoomByID(id uint) (*response.RoomResponse, error)
	GetWaitingRooms() ([]*response.RoomResponse, error)
	GetAvailableRooms(search *string, quizID *uint, page, pageSize int) ([]*response.RoomResponse, int64, error)
	JoinRoom(roomID uint, userID *uint, guestName, password *string) (*response.ParticipantResponse, error)
	JoinRoomByCode(code string, userID *uint, guestName, password *string) (*response.RoomResponse, *response.ParticipantResponse, error)
	LeaveRoom(roomID uint, userID *uint, guestName *string) error
	CloseRoom(roomID, hostID uint) error
	UpdateRoom(roomID, hostID uint, req *request.UpdateRoomRequest) (*response.RoomResponse, error)
	StartGame(roomID, hostID uint) (*response.RoomResponse, error)
}

type roomService struct {
	roomRepo        repository.RoomRepository
	participantRepo repository.RoomParticipantRepository
	quizRepo        repository.QuizRepository
	userRepo        repository.UserRepository
}

// NewRoomService creates a new room service
func NewRoomService(
	roomRepo repository.RoomRepository,
	participantRepo repository.RoomParticipantRepository,
	quizRepo repository.QuizRepository,
	userRepo repository.UserRepository,
) RoomService {
	return &roomService{
		roomRepo:        roomRepo,
		participantRepo: participantRepo,
		quizRepo:        quizRepo,
		userRepo:        userRepo,
	}
}

func (s *roomService) CreateRoom(hostID uint, req *request.RoomRequest) (*response.RoomResponse, error) {
	// Verify quiz exists
	_, err := s.quizRepo.FindByID(req.QuizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQuizNotFound
		}
		return nil, err
	}

	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	maxPlayers := 10
	if req.MaxPlayers > 0 {
		maxPlayers = req.MaxPlayers
	}

	room := &domain.Room{
		Name:       req.Name,
		Code:       generateRoomCode(),
		QuizID:     req.QuizID,
		HostID:     hostID,
		Password:   req.Password,
		IsPublic:   isPublic,
		MaxPlayers: maxPlayers,
		Status:     enums.RoomStatusWaiting,
	}

	if err := s.roomRepo.Create(room); err != nil {
		return nil, err
	}

	// Add host as participant
	participant := &domain.RoomParticipant{
		RoomID:   room.ID,
		UserID:   &hostID,
		IsHost:   true,
		JoinedAt: time.Now(),
	}
	if err := s.participantRepo.Create(participant); err != nil {
		return nil, err
	}

	return s.GetRoomByID(room.ID)
}

func (s *roomService) GetRoomByCode(code string) (*response.RoomResponse, error) {
	room, err := s.roomRepo.FindByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoomNotFound
		}
		return nil, err
	}
	return response.FromRoom(room), nil
}

func (s *roomService) GetRoomByID(id uint) (*response.RoomResponse, error) {
	room, err := s.roomRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoomNotFound
		}
		return nil, err
	}
	return response.FromRoom(room), nil
}

func (s *roomService) GetWaitingRooms() ([]*response.RoomResponse, error) {
	rooms, err := s.roomRepo.FindWaiting()
	if err != nil {
		return nil, err
	}
	return response.FromRooms(rooms), nil
}

func (s *roomService) GetAvailableRooms(search *string, quizID *uint, page, pageSize int) ([]*response.RoomResponse, int64, error) {
	rooms, total, err := s.roomRepo.FindAvailable(search, quizID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return response.FromRooms(rooms), total, nil
}

func (s *roomService) JoinRoom(roomID uint, userID *uint, guestName, password *string) (*response.ParticipantResponse, error) {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoomNotFound
		}
		return nil, err
	}

	// Check room status
	if room.Status != enums.RoomStatusWaiting {
		return nil, ErrRoomNotWaiting
	}

	// Check password
	if room.HasPassword() {
		if password == nil || *password != *room.Password {
			return nil, ErrWrongPassword
		}
	}

	// Check if room is full
	count, _ := s.participantRepo.CountByRoomID(roomID)
	if int(count) >= room.MaxPlayers {
		return nil, ErrRoomFull
	}

	// Check if user already in room
	if userID != nil {
		existing, _ := s.participantRepo.FindByRoomAndUser(roomID, *userID)
		if existing != nil && existing.LeftAt == nil {
			return nil, ErrAlreadyInRoom
		}
	}

	isGuest := userID == nil
	participant := &domain.RoomParticipant{
		RoomID:    roomID,
		UserID:    userID,
		IsGuest:   isGuest,
		GuestName: guestName,
		JoinedAt:  time.Now(),
	}

	if err := s.participantRepo.Create(participant); err != nil {
		return nil, err
	}

	// Load user if not guest
	if userID != nil {
		user, _ := s.userRepo.FindByID(*userID)
		participant.User = user
	}

	return response.FromRoomParticipant(participant), nil
}

func (s *roomService) JoinRoomByCode(code string, userID *uint, guestName, password *string) (*response.RoomResponse, *response.ParticipantResponse, error) {
	room, err := s.roomRepo.FindByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrRoomNotFound
		}
		return nil, nil, err
	}

	participant, err := s.JoinRoom(room.ID, userID, guestName, password)
	if err != nil {
		return nil, nil, err
	}

	roomResp, _ := s.GetRoomByID(room.ID)
	return roomResp, participant, nil
}

func (s *roomService) LeaveRoom(roomID uint, userID *uint, guestName *string) error {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoomNotFound
		}
		return err
	}

	var participant *domain.RoomParticipant
	if userID != nil {
		participant, _ = s.participantRepo.FindByRoomAndUser(roomID, *userID)
	}

	if participant == nil {
		return ErrNotInRoom
	}

	// Mark as left
	now := time.Now()
	participant.LeftAt = &now
	if err := s.participantRepo.Update(participant); err != nil {
		return err
	}

	// If host leaves and room is waiting, cancel the room
	if participant.IsHost && room.Status == enums.RoomStatusWaiting {
		room.Status = enums.RoomStatusCancelled
		return s.roomRepo.Update(room)
	}

	return nil
}

func (s *roomService) CloseRoom(roomID, hostID uint) error {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoomNotFound
		}
		return err
	}

	if room.HostID != hostID {
		return ErrNotRoomHost
	}

	room.Status = enums.RoomStatusCancelled
	now := time.Now()
	room.EndTime = &now

	return s.roomRepo.Update(room)
}

func (s *roomService) UpdateRoom(roomID, hostID uint, req *request.UpdateRoomRequest) (*response.RoomResponse, error) {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoomNotFound
		}
		return nil, err
	}

	if room.HostID != hostID {
		return nil, ErrNotRoomHost
	}

	if req.Name != nil {
		room.Name = *req.Name
	}
	if req.MaxPlayers != nil {
		room.MaxPlayers = *req.MaxPlayers
	}
	if req.Password != nil {
		room.Password = req.Password
	}
	if req.IsPublic != nil {
		room.IsPublic = *req.IsPublic
	}

	if err := s.roomRepo.Update(room); err != nil {
		return nil, err
	}

	return s.GetRoomByID(roomID)
}

func (s *roomService) StartGame(roomID, hostID uint) (*response.RoomResponse, error) {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoomNotFound
		}
		return nil, err
	}

	if room.HostID != hostID {
		return nil, ErrNotRoomHost
	}

	if room.Status != enums.RoomStatusWaiting {
		return nil, ErrRoomNotWaiting
	}

	room.Status = enums.RoomStatusInProgress
	now := time.Now()
	room.StartTime = &now

	if err := s.roomRepo.Update(room); err != nil {
		return nil, err
	}

	return s.GetRoomByID(roomID)
}

// generateRoomCode generates a random 6-character room code
func generateRoomCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		code[i] = charset[n.Int64()]
	}
	return string(code)
}
