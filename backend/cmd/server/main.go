package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/takumi/fastchem/internal/database"
	"github.com/takumi/fastchem/internal/handlers"
	"github.com/takumi/fastchem/internal/middleware"
	"github.com/takumi/fastchem/internal/services"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize database
	database.Init()
	defer database.Close()

	// Periodically clean up expired questions from the in-memory store
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			services.GlobalQuestionStore.CleanupOlderThan(10 * time.Minute)
		}
	}()

	// Periodically clean up stale active matches (no activity for 15 minutes)
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			services.GlobalMatchStore.CleanupOlderThan(15 * time.Minute)
		}
	}()

	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// Rate limiters
	// Answer endpoint: max 2 answers per second burst, refill 1/s
	answerLimiter := middleware.NewRateLimiter(1, 3, time.Second)
	// General API: 30 requests per second for unauthenticated endpoints
	ipLimiter := middleware.NewIPRateLimiter(10, 30, time.Second)

	// Initialize services and handlers
	generator := services.NewQuestionGenerator()
	questionHandler := handlers.NewQuestionHandler(generator)
	authHandler := handlers.NewAuthHandler()
	leaderboardHandler := handlers.NewLeaderboardHandler()
	matchHandler := handlers.NewMatchHandler(generator)

	// Initialize ranked match services
	rankedMatchService := services.NewRankedMatchService()
	matchmakingQueue := services.NewMatchmakingQueue(rankedMatchService)
	rankedHandler := handlers.NewRankedHandler(rankedMatchService, matchmakingQueue)

	// Periodically clean up stale ranked matches (no activity for 30 minutes)
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			rankedMatchService.CleanupStaleMatches(30 * time.Minute)
		}
	}()

	// Routes
	api := router.Group("/api")
	api.Use(ipLimiter.Middleware())
	{
		api.GET("/question", questionHandler.GetQuestion)
		api.POST("/validate", questionHandler.ValidateAnswer)
		api.POST("/answer", questionHandler.SubmitAnswer)
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthRequired(), authHandler.Me)
		}

		// Leaderboard (public)
		api.GET("/leaderboard", leaderboardHandler.GetLeaderboard)

		// Profile (public)
		api.GET("/profile/:username", leaderboardHandler.GetProfile)

		// Score routes (authenticated)
		scores := api.Group("/scores")
		scores.Use(middleware.AuthRequired())
		{
			scores.POST("", leaderboardHandler.SubmitScore)
			scores.GET("/me", leaderboardHandler.GetUserScores)
		}

		// Match routes (authenticated, rate-limited)
		match := api.Group("/match")
		match.Use(middleware.AuthRequired())
		{
			match.POST("/start", matchHandler.StartMatch)
			match.POST("/answer", answerLimiter.Middleware(), matchHandler.AnswerMatch)
			match.POST("/end", matchHandler.EndMatch)
		}

		// Ranked match routes
		ranked := api.Group("/ranked")
		{
			// WebSocket endpoint (auth via query param)
			ranked.GET("/ws", rankedHandler.HandleWebSocket)

			// REST endpoints (authenticated)
			ranked.GET("/stats", middleware.AuthRequired(), rankedHandler.GetRankedStats)
			ranked.GET("/history", middleware.AuthRequired(), rankedHandler.GetRankedHistory)
			ranked.GET("/leaderboard", rankedHandler.GetRankedLeaderboard)
		}
	}

	log.Printf("FastChem backend starting on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
