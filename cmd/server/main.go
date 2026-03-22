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
	infraconfig "github.com/huy/quizme-backend/internal/infra/config"
	infraMiddleware "github.com/huy/quizme-backend/internal/infra/middleware"
	infrarouter "github.com/huy/quizme-backend/internal/infra/router"
	infrastorage "github.com/huy/quizme-backend/internal/infra/storage"
	"github.com/huy/quizme-backend/internal/pkg/jwt"

	authhandler "github.com/huy/quizme-backend/internal/features/auth/handler"
	authrepo "github.com/huy/quizme-backend/internal/features/auth/repository"
	authservice "github.com/huy/quizme-backend/internal/features/auth/service"

	categoryhandler "github.com/huy/quizme-backend/internal/features/category/handler"
	categoryrepo "github.com/huy/quizme-backend/internal/features/category/repository"
	categoryservice "github.com/huy/quizme-backend/internal/features/category/service"

	gamehandler "github.com/huy/quizme-backend/internal/features/game/handler"
	gamerepo "github.com/huy/quizme-backend/internal/features/game/repository"
	gameservice "github.com/huy/quizme-backend/internal/features/game/service"
	gamewebsocket "github.com/huy/quizme-backend/internal/features/game/websocket"

	quizhandler "github.com/huy/quizme-backend/internal/features/quiz/handler"
	quizrepo "github.com/huy/quizme-backend/internal/features/quiz/repository"
	quizservice "github.com/huy/quizme-backend/internal/features/quiz/service"

	roomhandler "github.com/huy/quizme-backend/internal/features/room/handler"
	roomrepo "github.com/huy/quizme-backend/internal/features/room/repository"
	roomservice "github.com/huy/quizme-backend/internal/features/room/service"

	userhandler "github.com/huy/quizme-backend/internal/features/user/handler"
	userrepo "github.com/huy/quizme-backend/internal/features/user/repository"
	userservice "github.com/huy/quizme-backend/internal/features/user/service"
)

func main() {
	// Load configuration
	cfg, err := infraconfig.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := infraconfig.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize JWT provider
	jwtProvider := jwt.NewJWTProvider(cfg.JWT.Secret, cfg.JWT.ExpirationMs, cfg.JWT.RefreshExpirationMs)

	// Initialize Cloudinary service
	cloudinaryService, err := infrastorage.NewCloudinaryService(cfg)
	if err != nil {
		log.Printf("Warning: Failed to initialize Cloudinary service: %v", err)
		cloudinaryService = nil // Continue without Cloudinary
	}

	// Initialize Repositories
	// Auth repositories
	authUserRepo := authrepo.NewUserRepository(db)
	authUserProfileRepo := authrepo.NewUserProfileRepository(db)
	authRefreshTokenRepo := authrepo.NewRefreshTokenRepository(db)

	// User repositories
	userUserRepo := userrepo.NewUserRepository(db)
	userProfileRepo := userrepo.NewUserProfileRepository(db)

	// Category repositories
	categoryRepo := categoryrepo.NewCategoryRepository(db)

	// Quiz repositories
	quizRepo := quizrepo.NewQuizRepository(db)
	questionRepo := quizrepo.NewQuestionRepository(db)
	questionOptionRepo := quizrepo.NewQuestionOptionRepository(db)

	// Room repositories
	roomRoomRepo := roomrepo.NewRoomRepository(db)
	roomParticipantRepo := roomrepo.NewRoomParticipantRepository(db)
	roomChatRepo := roomrepo.NewRoomChatRepository(db)

	// Game repositories
	gameResultRepo := gamerepo.NewGameResultRepository(db)
	gamePlayerAnswerRepo := gamerepo.NewGamePlayerAnswerRepository(db)

	// Initialize WebSocket hub
	wsHub := gamewebsocket.NewHub()
	go wsHub.Run()

	// Initialize Services
	// Auth service
	authSvc := authservice.NewAuthService(authUserRepo, authUserProfileRepo, authRefreshTokenRepo, jwtProvider)

	// User service
	userSvc := userservice.NewUserService(userUserRepo, userProfileRepo)

	// Category service
	categorySvc := categoryservice.NewCategoryService(categoryRepo)

	// Quiz services
	quizSvc := quizservice.NewQuizService(quizRepo, questionRepo, questionOptionRepo, categoryRepo)
	questionSvc := quizservice.NewQuestionService(questionRepo, questionOptionRepo, quizRepo)

	// Room services
	roomSvc := roomservice.NewRoomService(roomRoomRepo, roomParticipantRepo, quizRepo, userUserRepo)
	chatSvc := roomservice.NewChatService(roomChatRepo, roomRoomRepo)

	// Game services
	gameProgressSvc := gameservice.NewGameProgressService(quizRepo, questionRepo, questionOptionRepo)
	gameResultSvc := gameservice.NewGameResultService(gameResultRepo, gamePlayerAnswerRepo, roomParticipantRepo)
	gameSessionSvc := gameservice.NewGameSessionService(
		wsHub,
		gameProgressSvc,
		gameResultSvc,
		roomRoomRepo,
		roomParticipantRepo,
		quizRepo,
	)

	// Initialize middleware
	authMiddleware := infraMiddleware.NewAuthMiddleware(jwtProvider, userUserRepo)

	// Initialize Handlers
	authHandler := authhandler.NewAuthHandler(authSvc)
	userHandler := userhandler.NewUserHandler(userSvc, cloudinaryService)
	categoryHandler := categoryhandler.NewCategoryHandler(categorySvc)
	quizHandler := quizhandler.NewQuizHandler(quizSvc)
	questionHandler := quizhandler.NewQuestionHandler(questionSvc)
	roomHandler := roomhandler.NewRoomHandler(roomSvc)
	chatHandler := roomhandler.NewChatHandler(chatSvc)
	gameHandler := gamehandler.NewGameHandler(gameSessionSvc, roomSvc)
	wsHandler := gamehandler.NewWebSocketHandler(
		wsHub,
		authMiddleware,
		gameSessionSvc,
		roomSvc,
		chatSvc,
		roomParticipantRepo,
	)

	// Setup Gin router
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// Setup routes using router package
	handlers := &infrarouter.Handlers{
		Auth:      authHandler,
		User:      userHandler,
		Category:  categoryHandler,
		Quiz:      quizHandler,
		Question:  questionHandler,
		Room:      roomHandler,
		Chat:      chatHandler,
		Game:      gameHandler,
		WebSocket: wsHandler,
		Auth0:     authMiddleware,
	}

	infrarouter.SetupRoutes(router, handlers, cfg.CORS.AllowedOrigins)

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
