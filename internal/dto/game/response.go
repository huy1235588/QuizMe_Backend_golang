package game

import "time"

// QuestionGameDTO is the question data sent during game
type QuestionGameDTO struct {
	QuestionID   uint            `json:"questionId"`
	Content      string          `json:"content"`
	ImageURL     *string         `json:"imageUrl,omitempty"`
	QuestionType string          `json:"questionType"`
	TimeLimit    int             `json:"timeLimit"`
	Points       int             `json:"points"`
	Options      []OptionGameDTO `json:"options"`
}

// OptionGameDTO is the option data sent during game (without correct flag)
type OptionGameDTO struct {
	OptionID uint   `json:"optionId"`
	Content  string `json:"content"`
}

// GameStateDTO represents the current state of a game
type GameStateDTO struct {
	GameActive      bool             `json:"gameActive"`
	CurrentQuestion *QuestionGameDTO `json:"currentQuestion,omitempty"`
	RemainingTime   int              `json:"remainingTime"`
	QuestionNumber  *int             `json:"questionNumber,omitempty"`
	TotalQuestions  int              `json:"totalQuestions"`
	Leaderboard     *LeaderboardDTO  `json:"leaderboard,omitempty"`
}

// Inactive returns an inactive game state
func (g *GameStateDTO) Inactive() *GameStateDTO {
	return &GameStateDTO{
		GameActive: false,
	}
}

// LeaderboardDTO represents the game leaderboard
type LeaderboardDTO struct {
	Rankings []PlayerRankingDTO `json:"rankings"`
}

// PlayerRankingDTO represents a player's ranking
type PlayerRankingDTO struct {
	UserID    *uint  `json:"userId,omitempty"`
	Username  string `json:"username"`
	Score     int    `json:"score"`
	Rank      int    `json:"rank"`
	IsCorrect *bool  `json:"isCorrect,omitempty"`
}

// QuestionResultDTO represents the result of a question
type QuestionResultDTO struct {
	QuestionID     uint               `json:"questionId"`
	CorrectOptions []uint             `json:"correctOptions"`
	Statistics     QuestionStatsDTO   `json:"statistics"`
	PlayerResults  []PlayerResultDTO  `json:"playerResults"`
}

// QuestionStatsDTO represents statistics for a question
type QuestionStatsDTO struct {
	TotalAnswers   int                `json:"totalAnswers"`
	CorrectCount   int                `json:"correctCount"`
	IncorrectCount int                `json:"incorrectCount"`
	AvgTime        float64            `json:"avgTime"`
	OptionStats    []OptionStatsDTO   `json:"optionStats"`
}

// OptionStatsDTO represents statistics for an option
type OptionStatsDTO struct {
	OptionID   uint    `json:"optionId"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// PlayerResultDTO represents a player's result for a question
type PlayerResultDTO struct {
	UserID          *uint  `json:"userId,omitempty"`
	Username        string `json:"username"`
	IsCorrect       bool   `json:"isCorrect"`
	Score           int    `json:"score"`
	AnswerTime      float64 `json:"answerTime"`
	SelectedOptions []uint  `json:"selectedOptions"`
}

// GameResultDTO represents the final result of a game
type GameResultDTO struct {
	RoomID         uint                    `json:"roomId"`
	QuizTitle      string                  `json:"quizTitle"`
	TotalQuestions int                     `json:"totalQuestions"`
	Duration       int                     `json:"duration"`
	FinalRankings  []FinalPlayerRankingDTO `json:"finalRankings"`
}

// FinalPlayerRankingDTO represents a player's final ranking
type FinalPlayerRankingDTO struct {
	UserID         *uint   `json:"userId,omitempty"`
	Username       string  `json:"username"`
	TotalScore     int     `json:"totalScore"`
	CorrectAnswers int     `json:"correctAnswers"`
	AvgAnswerTime  float64 `json:"avgAnswerTime"`
	Rank           int     `json:"rank"`
}

// GameStartEvent is sent when game starts
type GameStartEvent struct {
	Message string `json:"message"`
}

// NextQuestionEvent is sent before next question
type NextQuestionEvent struct {
	QuestionNumber int `json:"questionNumber"`
}

// AnswerResultEvent is sent after submitting an answer
type AnswerResultEvent struct {
	IsCorrect bool `json:"isCorrect"`
	Score     int  `json:"score"`
}

// PlayerDisconnectEvent is sent when a player disconnects
type PlayerDisconnectEvent struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
}

// GameEndEvent is sent when game ends
type GameEndEvent struct {
	Reason  string        `json:"reason,omitempty"`
	Message string        `json:"message,omitempty"`
	Result  *GameResultDTO `json:"result,omitempty"`
}

// ParticipantSession represents a participant in a game session (in-memory)
type ParticipantSession struct {
	ParticipantID    uint
	UserID           *uint
	Username         string
	Score            int
	Rank             int
	Answers          map[uint]*PlayerAnswer // questionID -> answer
	ConnectionStatus string
	JoinedAt         time.Time
	DisconnectedAt   *time.Time
	SessionIDs       map[string]bool
}

// PlayerAnswer represents a player's answer in memory
type PlayerAnswer struct {
	QuestionID      uint
	SelectedOptions []uint
	AnswerTime      float64
	IsCorrect       bool
	Score           int
}
