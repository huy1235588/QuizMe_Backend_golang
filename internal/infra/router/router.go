package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authhandler "github.com/huy/quizme-backend/internal/features/auth/handler"
	categoryhandler "github.com/huy/quizme-backend/internal/features/category/handler"
	gamehandler "github.com/huy/quizme-backend/internal/features/game/handler"
	quizhandler "github.com/huy/quizme-backend/internal/features/quiz/handler"
	roomhandler "github.com/huy/quizme-backend/internal/features/room/handler"
	userhandler "github.com/huy/quizme-backend/internal/features/user/handler"
	infraMiddleware "github.com/huy/quizme-backend/internal/infra/middleware"
)

// Handlers holds all feature handlers
type Handlers struct {
	Auth       *authhandler.AuthHandler
	User       *userhandler.UserHandler
	Category   *categoryhandler.CategoryHandler
	Quiz       *quizhandler.QuizHandler
	Question   *quizhandler.QuestionHandler
	Room       *roomhandler.RoomHandler
	Chat       *roomhandler.ChatHandler
	Game       *gamehandler.GameHandler
	WebSocket  *gamehandler.WebSocketHandler
	Auth0      *infraMiddleware.AuthMiddleware
}

// SetupRoutes configures all application routes
func SetupRoutes(engine *gin.Engine, handlers *Handlers, allowedOrigins []string) {
	// CORS middleware
	engine.Use(infraMiddleware.CORSMiddleware(allowedOrigins))

	// WebSocket route
	engine.GET("/ws", handlers.WebSocket.HandleConnection)

	// API routes group
	api := engine.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Auth.Login)
			auth.POST("/register", handlers.Auth.Register)
			auth.POST("/logout", handlers.Auth.Logout)
			auth.POST("/refresh-token", handlers.Auth.RefreshToken)
		}

		// User routes
		users := api.Group("/users")
		{
			users.GET("/:id", handlers.User.GetUserByID)
			users.GET("/top", handlers.User.GetTopUsers)
			users.GET("/count", handlers.User.GetUserCount)
			users.GET("/paged", handlers.User.GetPagedUsers)
			users.GET("/profile/:id", handlers.User.GetUserProfile)

			// Protected routes
			users.GET("/profile", handlers.Auth0.RequireAuth(), handlers.User.GetCurrentUserProfile)
			users.POST("/avatar/upload", handlers.Auth0.RequireAuth(), handlers.User.UploadAvatar)
			users.DELETE("/avatar", handlers.Auth0.RequireAuth(), handlers.User.RemoveAvatar)

			// Admin routes
			users.POST("/create", handlers.Auth0.RequireAuth(), handlers.Auth0.RequireAdmin(), handlers.User.CreateUser)
			users.PUT("/:id", handlers.Auth0.RequireAuth(), handlers.Auth0.RequireAdmin(), handlers.User.UpdateUser)
			users.DELETE("/:id", handlers.Auth0.RequireAuth(), handlers.Auth0.RequireAdmin(), handlers.User.DeleteUser)
			users.PUT("/:id/lock", handlers.Auth0.RequireAuth(), handlers.Auth0.RequireAdmin(), handlers.User.ToggleUserActiveStatus)
		}

		// Category routes
		categories := api.Group("/categories")
		{
			categories.GET("", handlers.Category.GetAllCategories)
			categories.GET("/:id", handlers.Category.GetCategoryByID)
			categories.GET("/active", handlers.Category.GetActiveCategories)

			// Admin routes
			categories.POST("", handlers.Auth0.RequireAuth(), handlers.Auth0.RequireAdmin(), handlers.Category.CreateCategory)
			categories.PUT("/:id", handlers.Auth0.RequireAuth(), handlers.Auth0.RequireAdmin(), handlers.Category.UpdateCategory)
			categories.DELETE("/:id", handlers.Auth0.RequireAuth(), handlers.Auth0.RequireAdmin(), handlers.Category.DeleteCategory)
		}

		// Quiz routes
		quizzes := api.Group("/quizzes")
		{
			quizzes.GET("", handlers.Quiz.GetAllQuizzes)
			quizzes.GET("/:id", handlers.Quiz.GetQuizByID)
			quizzes.GET("/public", handlers.Quiz.GetPublicQuizzes)
			quizzes.GET("/difficulty/:difficulty", handlers.Quiz.GetQuizzesByDifficulty)
			quizzes.GET("/paged", handlers.Quiz.GetPagedQuizzes)

			// Protected routes
			quizzes.POST("", handlers.Auth0.RequireAuth(), handlers.Quiz.CreateQuiz)
			quizzes.PUT("/:id", handlers.Auth0.RequireAuth(), handlers.Quiz.UpdateQuiz)
			quizzes.DELETE("/:id", handlers.Auth0.RequireAuth(), handlers.Quiz.DeleteQuiz)
		}

		// Question routes
		questions := api.Group("/questions")
		{
			questions.GET("/:id", handlers.Question.GetQuestionByID)
			questions.GET("/quiz/:quizId", handlers.Question.GetQuestionsByQuizID)
			questions.POST("", handlers.Auth0.RequireAuth(), handlers.Question.CreateQuestion)
			questions.POST("/batch", handlers.Auth0.RequireAuth(), handlers.Question.CreateBatchQuestions)
			questions.PUT("/:id", handlers.Auth0.RequireAuth(), handlers.Question.UpdateQuestion)
			questions.DELETE("/:id", handlers.Auth0.RequireAuth(), handlers.Question.DeleteQuestion)
		}

		// Room routes
		rooms := api.Group("/rooms")
		{
			rooms.GET("/:code", handlers.Room.GetRoomByCode)
			rooms.GET("/waiting", handlers.Room.GetWaitingRooms)
			rooms.GET("/available", handlers.Room.GetAvailableRooms)

			// Protected/Optional auth routes
			rooms.POST("", handlers.Auth0.RequireAuth(), handlers.Room.CreateRoom)
			rooms.POST("/join", handlers.Auth0.OptionalAuth(), handlers.Room.JoinRoomByCode)
			rooms.POST("/join/:roomId", handlers.Auth0.OptionalAuth(), handlers.Room.JoinRoomByID)
			rooms.DELETE("/leave/:roomId", handlers.Auth0.OptionalAuth(), handlers.Room.LeaveRoom)
			rooms.PATCH("/close/:roomId", handlers.Auth0.RequireAuth(), handlers.Room.CloseRoom)
			rooms.PATCH("/:roomId", handlers.Auth0.RequireAuth(), handlers.Room.UpdateRoom)
			rooms.POST("/start/:roomId", handlers.Auth0.RequireAuth(), handlers.Room.StartGame)
		}

		// Chat routes
		chat := api.Group("/chat")
		{
			chat.GET("/room/:roomId", handlers.Chat.GetChatHistory)
			chat.POST("/send", handlers.Auth0.OptionalAuth(), handlers.Chat.SendMessage)
		}

		// Game routes
		game := api.Group("/game")
		{
			game.GET("/state/:roomId", handlers.Auth0.OptionalAuth(), handlers.Game.GetGameState)
			game.POST("/init/:roomId", handlers.Auth0.RequireAuth(), handlers.Game.InitGame)
			game.POST("/start/:roomId", handlers.Auth0.RequireAuth(), handlers.Game.StartGame)
		}
	}

	// Health check
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
