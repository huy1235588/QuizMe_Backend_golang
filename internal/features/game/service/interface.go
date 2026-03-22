package service

import (
	"github.com/huy/quizme-backend/internal/dto/game"
	roomDTO "github.com/huy/quizme-backend/internal/features/room/dto"
)

// GameSessionService handles game session operations
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

// GameProgressService handles game progress tracking
type GameProgressService interface {
	LoadQuizAndPrepareQuestions(quizID uint) ([]*game.QuestionGameDTO, error)
	CalculateResults(session *GameSession, questionIndex int) *game.QuestionResultDTO
	GenerateLeaderboard(session *GameSession) *game.LeaderboardDTO
	GenerateFinalRankings(session *GameSession) []game.FinalPlayerRankingDTO
}

// GameResultService handles game result operations
type GameResultService interface {
	SaveGameResult(session *GameSession, result *game.GameResultDTO) error
	GetGameResult(gameResultID uint) (*game.GameResultDTO, error)
	GetGameResultsByRoom(roomID uint) ([]*game.GameResultDTO, error)
}

// ChatService handles real-time chat events (game context)
type ChatService interface {
	GetChatHistory(roomID uint, limit int) ([]*roomDTO.ChatMessageResponse, error)
	SendMessage(req *roomDTO.ChatMessageRequest, userID *uint) (*roomDTO.ChatMessageResponse, error)
}

