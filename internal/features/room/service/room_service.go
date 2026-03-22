package service

import (
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"github.com/huy/quizme-backend/internal/domain/enums"
	roomdto "github.com/huy/quizme-backend/internal/features/room/dto"
	roomdomain "github.com/huy/quizme-backend/internal/features/room/domain"
	roomrepo "github.com/huy/quizme-backend/internal/features/room/repository"
	quizrepo "github.com/huy/quizme-backend/internal/features/quiz/repository"
	userrepo "github.com/huy/quizme-backend/internal/features/user/repository"
	"gorm.io/gorm"
)

var (
	ErrRoomNotFound       = errors.New("room not found")
	ErrQuizNotFound       = errors.New("quiz not found")
	ErrRoomFull           = errors.New("room is full")
	ErrWrongPassword      = errors.New("wrong room password")
	ErrNotRoomHost        = errors.New("you are not the host of this room")
	ErrRoomNotWaiting     = errors.New("room is not in waiting status")
	ErrAlreadyInRoom      = errors.New("you are already in this room")
	ErrNotInRoom          = errors.New("you are not in this room")
)

type roomService struct {
	roomRepo        roomrepo.RoomRepository
	participantRepo roomrepo.RoomParticipantRepository
	quizRepo        quizrepo.QuizRepository
	userRepo        userrepo.UserRepository
}

// NewRoomService creates a new room service
func NewRoomService(
	roomRepo roomrepo.RoomRepository,
	participantRepo roomrepo.RoomParticipantRepository,
	quizRepo quizrepo.QuizRepository,
	userRepo userrepo.UserRepository,
) RoomService {
	return &roomService{
		roomRepo:        roomRepo,
		participantRepo: participantRepo,
		quizRepo:        quizRepo,
		userRepo:        userRepo,
	}
}

func (s *roomService) CreateRoom(hostID uint, req *roomdto.RoomRequest) (*roomdto.RoomResponse, error) {
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

	room := &roomdomain.Room{
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
	participant := &roomdomain.RoomParticipant{
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

func (s *roomService) GetRoomByCode(code string) (*roomdto.RoomResponse, error) {
	room, err := s.roomRepo.FindByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoomNotFound
		}
		return nil, err
	}
	return roomdto.FromRoom(room), nil
}

func (s *roomService) GetRoomByID(id uint) (*roomdto.RoomResponse, error) {
	room, err := s.roomRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoomNotFound
		}
		return nil, err
	}
	return roomdto.FromRoom(room), nil
}

func (s *roomService) GetWaitingRooms() ([]*roomdto.RoomResponse, error) {
	rooms, err := s.roomRepo.FindWaiting()
	if err != nil {
		return nil, err
	}
	return roomdto.FromRooms(rooms), nil
}

func (s *roomService) GetAvailableRooms(search *string, quizID *uint, page, pageSize int) ([]*roomdto.RoomResponse, int64, error) {
	rooms, total, err := s.roomRepo.FindAvailable(search, quizID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return roomdto.FromRooms(rooms), total, nil
}

func (s *roomService) JoinRoom(roomID uint, userID *uint, guestName, password *string) (*roomdto.ParticipantResponse, error) {
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
	participant := &roomdomain.RoomParticipant{
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

	return roomdto.FromRoomParticipant(participant), nil
}

func (s *roomService) JoinRoomByCode(code string, userID *uint, guestName, password *string) (*roomdto.RoomResponse, *roomdto.ParticipantResponse, error) {
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

	var participant *roomdomain.RoomParticipant
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

func (s *roomService) UpdateRoom(roomID, hostID uint, req *roomdto.UpdateRoomRequest) (*roomdto.RoomResponse, error) {
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

func (s *roomService) StartGame(roomID, hostID uint) (*roomdto.RoomResponse, error) {
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
