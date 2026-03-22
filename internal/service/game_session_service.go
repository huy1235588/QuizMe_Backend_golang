package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/huy/quizme-backend/internal/domain/enums"
	"github.com/huy/quizme-backend/internal/dto/game"
	"github.com/huy/quizme-backend/internal/repository"
	"github.com/huy/quizme-backend/internal/websocket"
)

// GameSession represents an in-memory game session
type GameSession struct {
	RoomID               uint
	QuizID               uint
	Status               enums.GameStatus
	CurrentQuestionIndex int
	Participants         map[uint]*game.ParticipantSession // participantID -> session
	StartTime            time.Time
	EndTime              *time.Time
	Questions            []*game.QuestionGameDTO
	cancelFunc           context.CancelFunc
	mu                   sync.RWMutex
}

// GameSessionService manages game sessions
type GameSessionService interface {
	InitGameSession(roomID uint) (*GameSession, error)
	StartGame(roomID uint) error
	ProcessAnswer(roomID uint, participantID uint, answer *game.AnswerRequest) error
	GetGameState(roomID uint) *game.GameStateDTO
	DisconnectPlayer(roomID uint, participantID uint)
	ReconnectPlayer(roomID uint, participantID uint, clientID string) *game.GameStateDTO
	EndGame(roomID uint) *game.GameResultDTO
	GetSession(roomID uint) *GameSession
	IsGameActive(roomID uint) bool
}

type gameSessionService struct {
	hub                 *websocket.Hub
	progressService     GameProgressService
	resultService       GameResultService
	roomRepo            repository.RoomRepository
	participantRepo     repository.RoomParticipantRepository
	quizRepo            repository.QuizRepository
	sessions            map[uint]*GameSession
	mu                  sync.RWMutex
}

// NewGameSessionService creates a new game session service
func NewGameSessionService(
	hub *websocket.Hub,
	progressService GameProgressService,
	resultService GameResultService,
	roomRepo repository.RoomRepository,
	participantRepo repository.RoomParticipantRepository,
	quizRepo repository.QuizRepository,
) GameSessionService {
	return &gameSessionService{
		hub:             hub,
		progressService: progressService,
		resultService:   resultService,
		roomRepo:        roomRepo,
		participantRepo: participantRepo,
		quizRepo:        quizRepo,
		sessions:        make(map[uint]*GameSession),
	}
}

func (s *gameSessionService) InitGameSession(roomID uint) (*GameSession, error) {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		return nil, err
	}

	// Load questions
	questions, err := s.progressService.LoadQuizAndPrepareQuestions(room.QuizID)
	if err != nil {
		return nil, err
	}

	// Load participants
	participants, err := s.participantRepo.FindByRoomID(roomID)
	if err != nil {
		return nil, err
	}

	session := &GameSession{
		RoomID:               roomID,
		QuizID:               room.QuizID,
		Status:               enums.GameStatusWaiting,
		CurrentQuestionIndex: -1,
		Participants:         make(map[uint]*game.ParticipantSession),
		Questions:            questions,
	}

	for _, p := range participants {
		if p.LeftAt != nil {
			continue // Skip participants who left
		}

		username := "Guest"
		if p.User != nil {
			username = p.User.Username
		} else if p.GuestName != nil {
			username = *p.GuestName
		}

		session.Participants[p.ID] = &game.ParticipantSession{
			ParticipantID:    p.ID,
			UserID:           p.UserID,
			Username:         username,
			Score:            0,
			Rank:             0,
			Answers:          make(map[uint]*game.PlayerAnswer),
			ConnectionStatus: enums.ConnectionStatusActive.String(),
			JoinedAt:         p.JoinedAt,
			SessionIDs:       make(map[string]bool),
		}
	}

	s.mu.Lock()
	s.sessions[roomID] = session
	s.mu.Unlock()

	return session, nil
}

func (s *gameSessionService) StartGame(roomID uint) error {
	s.mu.RLock()
	session, ok := s.sessions[roomID]
	s.mu.RUnlock()

	if !ok {
		return ErrRoomNotFound
	}

	session.mu.Lock()
	session.Status = enums.GameStatusInProgress
	session.StartTime = time.Now()
	session.mu.Unlock()

	// Send game start event
	msg, _ := websocket.NewMessage(websocket.MessageTypeGameStart, game.GameStartEvent{
		Message: "Game started!",
	})
	s.hub.BroadcastToRoom(roomID, msg)

	// Start first question
	s.startQuestion(roomID, 0)

	return nil
}

func (s *gameSessionService) startQuestion(roomID uint, questionIndex int) {
	s.mu.RLock()
	session, ok := s.sessions[roomID]
	s.mu.RUnlock()

	if !ok {
		return
	}

	session.mu.Lock()
	if questionIndex >= len(session.Questions) {
		session.mu.Unlock()
		s.EndGame(roomID)
		return
	}

	session.CurrentQuestionIndex = questionIndex
	session.Status = enums.GameStatusInProgress
	session.StartTime = time.Now()

	// Cancel previous timer if any
	if session.cancelFunc != nil {
		session.cancelFunc()
	}

	question := session.Questions[questionIndex]
	session.mu.Unlock()

	// Send question to all players
	msg, _ := websocket.NewMessage(websocket.MessageTypeQuestion, question)
	s.hub.BroadcastToRoom(roomID, msg)

	log.Printf("Starting question %d for room %d. Time limit: %d seconds",
		questionIndex+1, roomID, question.TimeLimit)

	// Start timer with goroutine
	ctx, cancel := context.WithCancel(context.Background())
	session.mu.Lock()
	session.cancelFunc = cancel
	session.mu.Unlock()

	go s.runQuestionTimer(ctx, roomID, question.TimeLimit)
}

func (s *gameSessionService) runQuestionTimer(ctx context.Context, roomID uint, totalSeconds int) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	remaining := totalSeconds

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			remaining--

			// Send timer event
			msg, _ := websocket.NewMessage(websocket.MessageTypeTimer, websocket.TimerPayload{
				RemainingTime: remaining,
				TotalTime:     totalSeconds,
			})
			s.hub.BroadcastToRoom(roomID, msg)

			if remaining <= 0 {
				s.endCurrentQuestion(roomID)
				return
			}
		}
	}
}

func (s *gameSessionService) endCurrentQuestion(roomID uint) {
	s.mu.RLock()
	session, ok := s.sessions[roomID]
	s.mu.RUnlock()

	if !ok {
		return
	}

	session.mu.Lock()
	if session.cancelFunc != nil {
		session.cancelFunc()
		session.cancelFunc = nil
	}
	session.Status = enums.GameStatusQuestionEnd
	questionIndex := session.CurrentQuestionIndex
	session.mu.Unlock()

	log.Printf("Ending question %d for room %d", questionIndex+1, roomID)

	// Calculate results
	result := s.progressService.CalculateResults(session, questionIndex)

	// Send question result
	msg, _ := websocket.NewMessage(websocket.MessageTypeQuestionResult, result)
	s.hub.BroadcastToRoom(roomID, msg)

	// Start showing results phase (5 seconds)
	s.startShowingResults(roomID)
}

func (s *gameSessionService) startShowingResults(roomID uint) {
	s.mu.RLock()
	session, ok := s.sessions[roomID]
	s.mu.RUnlock()

	if !ok {
		return
	}

	session.mu.Lock()
	session.Status = enums.GameStatusShowingResults
	session.mu.Unlock()

	// Wait 5 seconds then show leaderboard
	time.AfterFunc(5*time.Second, func() {
		s.startShowingLeaderboard(roomID)
	})
}

func (s *gameSessionService) startShowingLeaderboard(roomID uint) {
	s.mu.RLock()
	session, ok := s.sessions[roomID]
	s.mu.RUnlock()

	if !ok {
		return
	}

	session.mu.Lock()
	session.Status = enums.GameStatusShowingLeaderboard
	questionIndex := session.CurrentQuestionIndex
	totalQuestions := len(session.Questions)
	session.mu.Unlock()

	// Generate and send leaderboard
	leaderboard := s.progressService.GenerateLeaderboard(session)
	msg, _ := websocket.NewMessage(websocket.MessageTypeLeaderboard, leaderboard)
	s.hub.BroadcastToRoom(roomID, msg)

	// Wait 5 seconds then proceed
	time.AfterFunc(5*time.Second, func() {
		nextIndex := questionIndex + 1
		if nextIndex >= totalQuestions {
			s.EndGame(roomID)
		} else {
			s.startNextQuestionCountdown(roomID, nextIndex)
		}
	})
}

func (s *gameSessionService) startNextQuestionCountdown(roomID uint, nextIndex int) {
	s.mu.RLock()
	session, ok := s.sessions[roomID]
	s.mu.RUnlock()

	if !ok {
		return
	}

	session.mu.Lock()
	session.Status = enums.GameStatusNextQuestion
	session.mu.Unlock()

	// Send next question event
	msg, _ := websocket.NewMessage(websocket.MessageTypeNextQuestion, game.NextQuestionEvent{
		QuestionNumber: nextIndex + 1,
	})
	s.hub.BroadcastToRoom(roomID, msg)

	// Wait 5 seconds then start next question
	time.AfterFunc(5*time.Second, func() {
		s.startQuestion(roomID, nextIndex)
	})
}

func (s *gameSessionService) ProcessAnswer(roomID uint, participantID uint, answer *game.AnswerRequest) error {
	s.mu.RLock()
	session, ok := s.sessions[roomID]
	s.mu.RUnlock()

	if !ok {
		return ErrRoomNotFound
	}

	session.mu.Lock()
	defer session.mu.Unlock()

	if session.Status != enums.GameStatusInProgress {
		return nil // Silently ignore answers when not in progress
	}

	participant, ok := session.Participants[participantID]
	if !ok {
		return ErrNotInRoom
	}

	// Check if already answered
	if _, exists := participant.Answers[answer.QuestionID]; exists {
		return nil // Already answered
	}

	// Store answer
	participant.Answers[answer.QuestionID] = &game.PlayerAnswer{
		QuestionID:      answer.QuestionID,
		SelectedOptions: answer.SelectedOptions,
		AnswerTime:      answer.AnswerTime,
	}

	return nil
}

func (s *gameSessionService) GetGameState(roomID uint) *game.GameStateDTO {
	s.mu.RLock()
	session, ok := s.sessions[roomID]
	s.mu.RUnlock()

	if !ok {
		return &game.GameStateDTO{GameActive: false}
	}

	session.mu.RLock()
	defer session.mu.RUnlock()

	if !s.isActiveStatus(session.Status) {
		return &game.GameStateDTO{GameActive: false}
	}

	var currentQuestion *game.QuestionGameDTO
	if session.CurrentQuestionIndex >= 0 && session.CurrentQuestionIndex < len(session.Questions) {
		currentQuestion = session.Questions[session.CurrentQuestionIndex]
	}

	questionNumber := session.CurrentQuestionIndex + 1
	leaderboard := s.progressService.GenerateLeaderboard(session)

	return &game.GameStateDTO{
		GameActive:      true,
		CurrentQuestion: currentQuestion,
		RemainingTime:   s.getRemainingTime(session),
		QuestionNumber:  &questionNumber,
		TotalQuestions:  len(session.Questions),
		Leaderboard:     leaderboard,
	}
}

func (s *gameSessionService) DisconnectPlayer(roomID uint, participantID uint) {
	s.mu.RLock()
	session, ok := s.sessions[roomID]
	s.mu.RUnlock()

	if !ok {
		return
	}

	session.mu.Lock()
	defer session.mu.Unlock()

	participant, ok := session.Participants[participantID]
	if !ok {
		return
	}

	now := time.Now()
	participant.ConnectionStatus = enums.ConnectionStatusDisconnected.String()
	participant.DisconnectedAt = &now

	log.Printf("Player %d disconnected from room %d", participantID, roomID)

	// Send disconnect event
	msg, _ := websocket.NewMessage(websocket.MessageTypePlayerDisconnect, game.PlayerDisconnectEvent{
		UserID:   participantID,
		Username: participant.Username,
	})
	s.hub.BroadcastToRoom(roomID, msg)
}

func (s *gameSessionService) ReconnectPlayer(roomID uint, participantID uint, clientID string) *game.GameStateDTO {
	s.mu.RLock()
	session, ok := s.sessions[roomID]
	s.mu.RUnlock()

	if !ok {
		return &game.GameStateDTO{GameActive: false}
	}

	session.mu.Lock()
	participant, ok := session.Participants[participantID]
	if ok {
		participant.ConnectionStatus = enums.ConnectionStatusActive.String()
		participant.DisconnectedAt = nil
		participant.SessionIDs[clientID] = true
	}
	session.mu.Unlock()

	log.Printf("Player %d reconnected to room %d", participantID, roomID)

	return s.GetGameState(roomID)
}

func (s *gameSessionService) EndGame(roomID uint) *game.GameResultDTO {
	s.mu.RLock()
	session, ok := s.sessions[roomID]
	s.mu.RUnlock()

	if !ok {
		return nil
	}

	session.mu.Lock()
	if session.cancelFunc != nil {
		session.cancelFunc()
		session.cancelFunc = nil
	}
	session.Status = enums.GameStatusCompleted
	now := time.Now()
	session.EndTime = &now
	session.mu.Unlock()

	log.Printf("Ending game for room %d", roomID)

	// Get quiz title
	quiz, _ := s.quizRepo.FindByID(session.QuizID)
	quizTitle := ""
	if quiz != nil {
		quizTitle = quiz.Title
	}

	// Calculate duration
	duration := 0
	if session.EndTime != nil {
		duration = int(session.EndTime.Sub(session.StartTime).Seconds())
	}

	// Generate final rankings
	finalRankings := s.progressService.GenerateFinalRankings(session)

	result := &game.GameResultDTO{
		RoomID:         roomID,
		QuizTitle:      quizTitle,
		TotalQuestions: len(session.Questions),
		Duration:       duration,
		FinalRankings:  finalRankings,
	}

	// Save result to database
	go s.resultService.SaveGameResult(session, result)

	// Send game end event
	msg, _ := websocket.NewMessage(websocket.MessageTypeGameEnd, game.GameEndEvent{
		Result: result,
	})
	s.hub.BroadcastToRoom(roomID, msg)

	// Clean up session after a delay
	time.AfterFunc(30*time.Second, func() {
		s.mu.Lock()
		delete(s.sessions, roomID)
		s.mu.Unlock()
		log.Printf("Cleaned up session for room %d", roomID)
	})

	return result
}

func (s *gameSessionService) GetSession(roomID uint) *GameSession {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessions[roomID]
}

func (s *gameSessionService) IsGameActive(roomID uint) bool {
	s.mu.RLock()
	session, ok := s.sessions[roomID]
	s.mu.RUnlock()

	if !ok {
		return false
	}

	session.mu.RLock()
	defer session.mu.RUnlock()

	return s.isActiveStatus(session.Status)
}

func (s *gameSessionService) isActiveStatus(status enums.GameStatus) bool {
	return status == enums.GameStatusInProgress ||
		status == enums.GameStatusQuestionEnd ||
		status == enums.GameStatusShowingResults ||
		status == enums.GameStatusShowingLeaderboard ||
		status == enums.GameStatusNextQuestion
}

func (s *gameSessionService) getRemainingTime(session *GameSession) int {
	if session.Status != enums.GameStatusInProgress {
		return 0
	}

	if session.CurrentQuestionIndex < 0 || session.CurrentQuestionIndex >= len(session.Questions) {
		return 0
	}

	question := session.Questions[session.CurrentQuestionIndex]
	elapsed := time.Since(session.StartTime).Seconds()
	remaining := question.TimeLimit - int(elapsed)

	if remaining < 0 {
		return 0
	}
	return remaining
}
