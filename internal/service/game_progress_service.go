package service

import (
	"github.com/huy/quizme-backend/internal/dto/game"
	"github.com/huy/quizme-backend/internal/repository"
)

// GameProgressService handles game progress calculations
type GameProgressService interface {
	LoadQuizAndPrepareQuestions(quizID uint) ([]*game.QuestionGameDTO, error)
	CalculateResults(session *GameSession, questionIndex int) *game.QuestionResultDTO
	GenerateLeaderboard(session *GameSession) *game.LeaderboardDTO
	GenerateFinalRankings(session *GameSession) []game.FinalPlayerRankingDTO
}

type gameProgressService struct {
	quizRepo           repository.QuizRepository
	questionRepo       repository.QuestionRepository
	questionOptionRepo repository.QuestionOptionRepository
}

// NewGameProgressService creates a new game progress service
func NewGameProgressService(
	quizRepo repository.QuizRepository,
	questionRepo repository.QuestionRepository,
	questionOptionRepo repository.QuestionOptionRepository,
) GameProgressService {
	return &gameProgressService{
		quizRepo:           quizRepo,
		questionRepo:       questionRepo,
		questionOptionRepo: questionOptionRepo,
	}
}

func (s *gameProgressService) LoadQuizAndPrepareQuestions(quizID uint) ([]*game.QuestionGameDTO, error) {
	questions, err := s.questionRepo.FindByQuizID(quizID)
	if err != nil {
		return nil, err
	}

	var questionDTOs []*game.QuestionGameDTO
	for _, q := range questions {
		// Get options for this question
		options, err := s.questionOptionRepo.FindByQuestionID(q.ID)
		if err != nil {
			return nil, err
		}

		var optionDTOs []game.OptionGameDTO
		for _, opt := range options {
			optionDTOs = append(optionDTOs, game.OptionGameDTO{
				OptionID: opt.ID,
				Content:  opt.Content,
			})
		}

		questionDTOs = append(questionDTOs, &game.QuestionGameDTO{
			QuestionID:   q.ID,
			Content:      q.Content,
			ImageURL:     q.ImageURL,
			QuestionType: string(q.Type),
			TimeLimit:    q.TimeLimit,
			Points:       q.Points,
			Options:      optionDTOs,
		})
	}

	return questionDTOs, nil
}

func (s *gameProgressService) CalculateResults(session *GameSession, questionIndex int) *game.QuestionResultDTO {
	if questionIndex >= len(session.Questions) {
		return nil
	}

	question := session.Questions[questionIndex]

	// Get correct options for this question
	options, _ := s.questionOptionRepo.FindByQuestionID(question.QuestionID)
	var correctOptions []uint
	for _, opt := range options {
		if opt.IsCorrect {
			correctOptions = append(correctOptions, opt.ID)
		}
	}

	// Calculate statistics
	var totalAnswers, correctCount, incorrectCount int
	var totalTime float64
	optionCounts := make(map[uint]int)
	var playerResults []game.PlayerResultDTO

	for _, participant := range session.Participants {
		answer, hasAnswer := participant.Answers[question.QuestionID]
		if !hasAnswer {
			continue
		}

		totalAnswers++
		totalTime += answer.AnswerTime

		// Check if answer is correct
		isCorrect := s.isAnswerCorrect(answer.SelectedOptions, correctOptions)
		answer.IsCorrect = isCorrect

		if isCorrect {
			correctCount++
			// Calculate score based on time and points
			answer.Score = s.calculateScore(question.Points, answer.AnswerTime, question.TimeLimit)
			participant.Score += answer.Score
		} else {
			incorrectCount++
			answer.Score = 0
		}

		// Count option selections
		for _, optID := range answer.SelectedOptions {
			optionCounts[optID]++
		}

		playerResults = append(playerResults, game.PlayerResultDTO{
			UserID:          participant.UserID,
			Username:        participant.Username,
			IsCorrect:       isCorrect,
			Score:           answer.Score,
			AnswerTime:      answer.AnswerTime,
			SelectedOptions: answer.SelectedOptions,
		})
	}

	// Calculate option statistics
	var optionStats []game.OptionStatsDTO
	for _, opt := range options {
		count := optionCounts[opt.ID]
		percentage := 0.0
		if totalAnswers > 0 {
			percentage = float64(count) / float64(totalAnswers) * 100
		}
		optionStats = append(optionStats, game.OptionStatsDTO{
			OptionID:   opt.ID,
			Count:      count,
			Percentage: percentage,
		})
	}

	avgTime := 0.0
	if totalAnswers > 0 {
		avgTime = totalTime / float64(totalAnswers)
	}

	return &game.QuestionResultDTO{
		QuestionID:     question.QuestionID,
		CorrectOptions: correctOptions,
		Statistics: game.QuestionStatsDTO{
			TotalAnswers:   totalAnswers,
			CorrectCount:   correctCount,
			IncorrectCount: incorrectCount,
			AvgTime:        avgTime,
			OptionStats:    optionStats,
		},
		PlayerResults: playerResults,
	}
}

func (s *gameProgressService) GenerateLeaderboard(session *GameSession) *game.LeaderboardDTO {
	// Collect and sort participants by score
	type participantScore struct {
		participant *game.ParticipantSession
		score       int
	}

	var scores []participantScore
	for _, p := range session.Participants {
		scores = append(scores, participantScore{participant: p, score: p.Score})
	}

	// Sort by score descending
	for i := 0; i < len(scores)-1; i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[j].score > scores[i].score {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}

	// Build rankings
	var rankings []game.PlayerRankingDTO
	for i, ps := range scores {
		ps.participant.Rank = i + 1
		rankings = append(rankings, game.PlayerRankingDTO{
			UserID:   ps.participant.UserID,
			Username: ps.participant.Username,
			Score:    ps.participant.Score,
			Rank:     i + 1,
		})
	}

	return &game.LeaderboardDTO{Rankings: rankings}
}

func (s *gameProgressService) GenerateFinalRankings(session *GameSession) []game.FinalPlayerRankingDTO {
	leaderboard := s.GenerateLeaderboard(session)

	var finalRankings []game.FinalPlayerRankingDTO
	for _, ranking := range leaderboard.Rankings {
		// Find participant
		var correctAnswers int
		var totalTime float64
		var answerCount int

		for _, p := range session.Participants {
			if (p.UserID != nil && ranking.UserID != nil && *p.UserID == *ranking.UserID) ||
				(p.UserID == nil && ranking.UserID == nil && p.Username == ranking.Username) {
				for _, answer := range p.Answers {
					if answer.IsCorrect {
						correctAnswers++
					}
					totalTime += answer.AnswerTime
					answerCount++
				}
				break
			}
		}

		avgTime := 0.0
		if answerCount > 0 {
			avgTime = totalTime / float64(answerCount)
		}

		finalRankings = append(finalRankings, game.FinalPlayerRankingDTO{
			UserID:         ranking.UserID,
			Username:       ranking.Username,
			TotalScore:     ranking.Score,
			CorrectAnswers: correctAnswers,
			AvgAnswerTime:  avgTime,
			Rank:           ranking.Rank,
		})
	}

	return finalRankings
}

func (s *gameProgressService) isAnswerCorrect(selected, correct []uint) bool {
	if len(selected) != len(correct) {
		return false
	}

	selectedMap := make(map[uint]bool)
	for _, id := range selected {
		selectedMap[id] = true
	}

	for _, id := range correct {
		if !selectedMap[id] {
			return false
		}
	}

	return true
}

func (s *gameProgressService) calculateScore(basePoints int, answerTime float64, timeLimit int) int {
	// Score = basePoints * (1 - answerTime/timeLimit * 0.5)
	// Minimum 50% of base points for correct answer
	timeFactor := 1 - (answerTime/float64(timeLimit))*0.5
	if timeFactor < 0.5 {
		timeFactor = 0.5
	}
	return int(float64(basePoints) * timeFactor)
}
