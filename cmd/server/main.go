package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huy/quizme-backend/internal/config"
	"github.com/huy/quizme-backend/internal/handler"
	"github.com/huy/quizme-backend/internal/middleware"
	"github.com/huy/quizme-backend/internal/pkg/jwt"
	"github.com/huy/quizme-backend/internal/repository"
	"github.com/huy/quizme-backend/internal/service"
	"github.com/huy/quizme-backend/internal/service/storage"
	"github.com/huy/quizme-backend/internal/websocket"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize JWT provider
	jwtProvider := jwt.NewJWTProvider(cfg.JWT.Secret, cfg.JWT.ExpirationMs, cfg.JWT.RefreshExpirationMs)

	// Initialize Cloudinary service
	cloudinaryService, err := storage.NewCloudinaryService(cfg)
	if err != nil {
		log.Printf("Warning: Failed to initialize Cloudinary service: %v", err)
		cloudinaryService = nil // Continue without Cloudinary
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	userProfileRepo := repository.NewUserProfileRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	quizRepo := repository.NewQuizRepository(db)
	questionRepo := repository.NewQuestionRepository(db)
	questionOptionRepo := repository.NewQuestionOptionRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	roomParticipantRepo := repository.NewRoomParticipantRepository(db)
	roomChatRepo := repository.NewRoomChatRepository(db)
	gameResultRepo := repository.NewGameResultRepository(db)
	gamePlayerAnswerRepo := repository.NewGamePlayerAnswerRepository(db)

	// Initialize WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// Initialize services
	authService := service.NewAuthService(userRepo, userProfileRepo, refreshTokenRepo, jwtProvider)
	userService := service.NewUserService(userRepo, userProfileRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	quizService := service.NewQuizService(quizRepo, questionRepo, questionOptionRepo, categoryRepo)
	questionService := service.NewQuestionService(questionRepo, questionOptionRepo, quizRepo)
	roomService := service.NewRoomService(roomRepo, roomParticipantRepo, quizRepo, userRepo)
	chatService := service.NewChatService(roomChatRepo, roomRepo)

	// Initialize game services
	gameProgressService := service.NewGameProgressService(quizRepo, questionRepo, questionOptionRepo)
	gameResultService := service.NewGameResultService(gameResultRepo, gamePlayerAnswerRepo, roomParticipantRepo)
	gameSessionService := service.NewGameSessionService(
		wsHub,
		gameProgressService,
		gameResultService,
		roomRepo,
		roomParticipantRepo,
		quizRepo,
	)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtProvider, userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService, cloudinaryService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	quizHandler := handler.NewQuizHandler(quizService)
	questionHandler := handler.NewQuestionHandler(questionService)
	roomHandler := handler.NewRoomHandler(roomService)
	chatHandler := handler.NewChatHandler(chatService)
	gameHandler := handler.NewGameHandler(gameSessionService, roomService)
	wsHandler := handler.NewWebSocketHandler(
		wsHub,
		authMiddleware,
		gameSessionService,
		roomService,
		chatService,
		roomParticipantRepo,
	)

	// Setup Gin router
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// CORS middleware
	router.Use(middleware.CORSMiddleware(cfg.CORS.AllowedOrigins))

	// WebSocket route
	router.GET("/ws", wsHandler.HandleConnection)

	// API routes
	api := router.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/refresh-token", authHandler.RefreshToken)
		}

		// User routes
		users := api.Group("/users")
		{
			users.GET("/:id", userHandler.GetUserByID)
			users.GET("/top", userHandler.GetTopUsers)
			users.GET("/count", userHandler.GetUserCount)
			users.GET("/paged", userHandler.GetPagedUsers)
			users.GET("/profile/:id", userHandler.GetUserProfile)

			// Protected routes
			users.GET("/profile", authMiddleware.RequireAuth(), userHandler.GetCurrentUserProfile)
			users.POST("/avatar/upload", authMiddleware.RequireAuth(), userHandler.UploadAvatar)
			users.DELETE("/avatar", authMiddleware.RequireAuth(), userHandler.RemoveAvatar)

			// Admin routes
			users.POST("/create", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin(), userHandler.CreateUser)
			users.PUT("/:id", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin(), userHandler.UpdateUser)
			users.DELETE("/:id", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin(), userHandler.DeleteUser)
			users.PUT("/:id/lock", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin(), userHandler.ToggleUserActiveStatus)
		}

		// Category routes
		categories := api.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAllCategories)
			categories.GET("/:id", categoryHandler.GetCategoryByID)
			categories.GET("/active", categoryHandler.GetActiveCategories)

			// Admin routes
			categories.POST("", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin(), categoryHandler.CreateCategory)
			categories.PUT("/:id", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin(), categoryHandler.UpdateCategory)
			categories.DELETE("/:id", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin(), categoryHandler.DeleteCategory)
		}

		// Quiz routes
		quizzes := api.Group("/quizzes")
		{
			quizzes.GET("", quizHandler.GetAllQuizzes)
			quizzes.GET("/:id", quizHandler.GetQuizByID)
			quizzes.GET("/public", quizHandler.GetPublicQuizzes)
			quizzes.GET("/difficulty/:difficulty", quizHandler.GetQuizzesByDifficulty)
			quizzes.GET("/paged", quizHandler.GetPagedQuizzes)

			// Protected routes
			quizzes.POST("", authMiddleware.RequireAuth(), quizHandler.CreateQuiz)
			quizzes.PUT("/:id", authMiddleware.RequireAuth(), quizHandler.UpdateQuiz)
			quizzes.DELETE("/:id", authMiddleware.RequireAuth(), quizHandler.DeleteQuiz)
		}

		// Question routes
		questions := api.Group("/questions")
		{
			questions.GET("/:id", questionHandler.GetQuestionByID)
			questions.GET("/quiz/:quizId", questionHandler.GetQuestionsByQuizID)
			questions.POST("", authMiddleware.RequireAuth(), questionHandler.CreateQuestion)
			questions.POST("/batch", authMiddleware.RequireAuth(), questionHandler.CreateBatchQuestions)
			questions.PUT("/:id", authMiddleware.RequireAuth(), questionHandler.UpdateQuestion)
			questions.DELETE("/:id", authMiddleware.RequireAuth(), questionHandler.DeleteQuestion)
		}

		// Room routes
		rooms := api.Group("/rooms")
		{
			rooms.GET("/:code", roomHandler.GetRoomByCode)
			rooms.GET("/waiting", roomHandler.GetWaitingRooms)
			rooms.GET("/available", roomHandler.GetAvailableRooms)

			// Protected/Optional auth routes
			rooms.POST("", authMiddleware.RequireAuth(), roomHandler.CreateRoom)
			rooms.POST("/join", authMiddleware.OptionalAuth(), roomHandler.JoinRoomByCode)
			rooms.POST("/join/:roomId", authMiddleware.OptionalAuth(), roomHandler.JoinRoomByID)
			rooms.DELETE("/leave/:roomId", authMiddleware.OptionalAuth(), roomHandler.LeaveRoom)
			rooms.PATCH("/close/:roomId", authMiddleware.RequireAuth(), roomHandler.CloseRoom)
			rooms.PATCH("/:roomId", authMiddleware.RequireAuth(), roomHandler.UpdateRoom)
			rooms.POST("/start/:roomId", authMiddleware.RequireAuth(), roomHandler.StartGame)
		}

		// Chat routes
		chat := api.Group("/chat")
		{
			chat.GET("/room/:roomId", chatHandler.GetChatHistory)
			chat.POST("/send", authMiddleware.OptionalAuth(), chatHandler.SendMessage)
		}

		// Game routes
		game := api.Group("/game")
		{
			game.GET("/state/:roomId", authMiddleware.OptionalAuth(), gameHandler.GetGameState)
			game.POST("/init/:roomId", authMiddleware.RequireAuth(), gameHandler.InitGame)
			game.POST("/start/:roomId", authMiddleware.RequireAuth(), gameHandler.StartGame)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
