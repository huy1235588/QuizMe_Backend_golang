package enums

// GameStatus represents the status of a game session
type GameStatus string

const (
	GameStatusWaiting          GameStatus = "WAITING"
	GameStatusInProgress       GameStatus = "IN_PROGRESS"
	GameStatusQuestionEnd      GameStatus = "QUESTION_END"
	GameStatusShowingResults   GameStatus = "SHOWING_RESULTS"
	GameStatusShowingLeaderboard GameStatus = "SHOWING_LEADERBOARD"
	GameStatusNextQuestion     GameStatus = "NEXT_QUESTION"
	GameStatusCompleted        GameStatus = "COMPLETED"
)

func (s GameStatus) String() string {
	return string(s)
}
