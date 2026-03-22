package service

import (
	"log"
	"time"

	"github.com/huy/quizme-backend/internal/domain"
	"github.com/huy/quizme-backend/internal/dto/game"
	"github.com/huy/quizme-backend/internal/repository"
)

// GameResultService handles saving and retrieving game results
type GameResultService interface {
	SaveGameResult(session *GameSession, result *game.GameResultDTO) error
	GetGameResult(gameResultID uint) (*domain.GameResult, error)
	GetGameResultsByRoom(roomID uint) ([]*domain.GameResult, error)
}

type gameResultService struct {
	gameResultRepo       repository.GameResultRepository
	gamePlayerAnswerRepo repository.GamePlayerAnswerRepository
	participantRepo      repository.RoomParticipantRepository
}

// NewGameResultService creates a new game result service
func NewGameResultService(
	gameResultRepo repository.GameResultRepository,
	gamePlayerAnswerRepo repository.GamePlayerAnswerRepository,
	participantRepo repository.RoomParticipantRepository,
) GameResultService {
	return &gameResultService{
		gameResultRepo:       gameResultRepo,
		gamePlayerAnswerRepo: gamePlayerAnswerRepo,
		participantRepo:      participantRepo,
	}
}

func (s *gameResultService) SaveGameResult(session *GameSession, result *game.GameResultDTO) error {
	session.mu.RLock()
	defer session.mu.RUnlock()

	// Calculate statistics
	var totalScore int
	var highestScore, lowestScore *int
	correctAnswersTotal := 0

	for _, p := range session.Participants {
		totalScore += p.Score

		if highestScore == nil || p.Score > *highestScore {
			score := p.Score
			highestScore = &score
		}
		if lowestScore == nil || p.Score < *lowestScore {
			score := p.Score
			lowestScore = &score
		}

		for _, answer := range p.Answers {
			if answer.IsCorrect {
				correctAnswersTotal++
			}
		}
	}

	participantCount := len(session.Participants)
	avgScore := 0.0
	if participantCount > 0 {
		avgScore = float64(totalScore) / float64(participantCount)
	}

	totalAnswersPossible := participantCount * len(session.Questions)
	completionRate := 0.0
	if totalAnswersPossible > 0 {
		totalAnswers := 0
		for _, p := range session.Participants {
			totalAnswers += len(p.Answers)
		}
		completionRate = float64(totalAnswers) / float64(totalAnswersPossible) * 100
	}

	// Create game result
	gameResult := &domain.GameResult{
		RoomID:           session.RoomID,
		QuizID:           session.QuizID,
		StartTime:        session.StartTime,
		EndTime:          session.EndTime,
		ParticipantCount: participantCount,
		QuestionCount:    len(session.Questions),
		AvgScore:         &avgScore,
		HighestScore:     highestScore,
		LowestScore:      lowestScore,
		CompletionRate:   &completionRate,
	}

	if err := s.gameResultRepo.Create(gameResult); err != nil {
		log.Printf("Error saving game result: %v", err)
		return err
	}

	// Save question results
	for i, question := range session.Questions {
		correctCount := 0
		incorrectCount := 0
		var totalTime float64
		answerCount := 0

		for _, p := range session.Participants {
			if answer, ok := p.Answers[question.QuestionID]; ok {
				if answer.IsCorrect {
					correctCount++
				} else {
					incorrectCount++
				}
				totalTime += answer.AnswerTime
				answerCount++
			}
		}

		avgTime := 0.0
		if answerCount > 0 {
			avgTime = totalTime / float64(answerCount)
		}

		questionResult := &domain.GameResultQuestion{
			GameResultID:   gameResult.ID,
			QuestionID:     question.QuestionID,
			CorrectCount:   correctCount,
			IncorrectCount: incorrectCount,
			AvgTime:        &avgTime,
		}

		gameResult.GameResultQuestions = append(gameResult.GameResultQuestions, *questionResult)
		_ = i // Suppress unused variable warning
	}

	// Save player answers
	var playerAnswers []*domain.GamePlayerAnswer
	for _, p := range session.Participants {
		for questionID, answer := range p.Answers {
			playerAnswer := &domain.GamePlayerAnswer{
				GameResultID:  gameResult.ID,
				ParticipantID: p.ParticipantID,
				QuestionID:    questionID,
				IsCorrect:     answer.IsCorrect,
				AnswerTime:    answer.AnswerTime,
				Score:         answer.Score,
				CreatedAt:     time.Now(),
			}

			// Add selected options
			for _, optID := range answer.SelectedOptions {
				playerAnswer.SelectedOptions = append(playerAnswer.SelectedOptions, domain.GamePlayerAnswerOption{
					OptionID: optID,
				})
			}

			playerAnswers = append(playerAnswers, playerAnswer)
		}
	}

	if len(playerAnswers) > 0 {
		if err := s.gamePlayerAnswerRepo.CreateBatch(playerAnswers); err != nil {
			log.Printf("Error saving player answers: %v", err)
		}
	}

	log.Printf("Saved game result %d for room %d", gameResult.ID, session.RoomID)
	return nil
}

func (s *gameResultService) GetGameResult(gameResultID uint) (*domain.GameResult, error) {
	return s.gameResultRepo.FindByID(gameResultID)
}

func (s *gameResultService) GetGameResultsByRoom(roomID uint) ([]*domain.GameResult, error) {
	return s.gameResultRepo.FindByRoomID(roomID)
}
