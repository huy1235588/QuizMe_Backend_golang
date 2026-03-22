package enums

// QuestionType represents different types of quiz questions
type QuestionType string

const (
	QuestionTypeQuiz       QuestionType = "QUIZ"
	QuestionTypeTrueFalse  QuestionType = "TRUE_FALSE"
	QuestionTypeTypeAnswer QuestionType = "TYPE_ANSWER"
	QuestionTypeQuizAudio  QuestionType = "QUIZ_AUDIO"
	QuestionTypeQuizVideo  QuestionType = "QUIZ_VIDEO"
	QuestionTypeCheckbox   QuestionType = "CHECKBOX"
	QuestionTypePoll       QuestionType = "POLL"
)

// IsValid checks if the question type is valid
func (qt QuestionType) IsValid() bool {
	switch qt {
	case QuestionTypeQuiz, QuestionTypeTrueFalse, QuestionTypeTypeAnswer,
		QuestionTypeQuizAudio, QuestionTypeQuizVideo, QuestionTypeCheckbox, QuestionTypePoll:
		return true
	}
	return false
}

// String returns the string representation
func (qt QuestionType) String() string {
	return string(qt)
}
