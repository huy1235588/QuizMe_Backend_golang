package game

// AnswerRequest is the request for submitting an answer
type AnswerRequest struct {
	QuestionID      uint    `json:"questionId" binding:"required"`
	SelectedOptions []uint  `json:"selectedOptions" binding:"required"`
	AnswerTime      float64 `json:"answerTime" binding:"required"`
}

// StartGameRequest is the request for starting a game
type StartGameRequest struct {
	RoomID uint `json:"roomId" binding:"required"`
}
