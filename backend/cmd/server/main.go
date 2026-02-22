package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/takumi/fastchem/internal/config"
	"github.com/takumi/fastchem/internal/database"
	"github.com/takumi/fastchem/internal/handlers"
	"github.com/takumi/fastchem/internal/middleware"
	"github.com/takumi/fastchem/internal/services"
)

func main() {
	cfg := config.Load()

	// Initialize JWT (must happen before any handler runs)
	middleware.InitJWT(cfg.JWTSecret, cfg.IsProduction())

	// Initialize database
	database.Init(cfg.DBPath)
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

	// Request ID middleware (added before CORS so every response gets it)
	router.Use(middleware.RequestID())

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
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
	rankedHandler := handlers.NewRankedHandler(rankedMatchService, matchmakingQueue, cfg.AllowedOriginsSet())

	// Initialize custom room service
	roomService := services.NewRoomService(rankedMatchService)
	roomHandler := handlers.NewRoomHandler(roomService, rankedMatchService, cfg.AllowedOriginsSet())

	// Periodically clean up stale ranked matches (no activity for 30 minutes)
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			rankedMatchService.CleanupStaleMatches(30 * time.Minute)
			roomService.CleanupStaleRooms(15 * time.Minute)
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

		// Custom room routes
		room := api.Group("/room")
		{
			// WebSocket endpoint (auth via query param)
			room.GET("/ws", roomHandler.HandleWebSocket)
		}
	}

	// ── Serve the Next.js static export (frontend) ──────────────
	// Resolve the frontend dist directory relative to the binary or CWD.
	staticDir := cfg.FrontendDir
	if staticDir == "" {
		// Default: assume project root layout  <root>/backend  and  <root>/frontend/out
		exe, _ := os.Executable()
		staticDir = filepath.Join(filepath.Dir(exe), "..", "frontend", "out")
	}
	// Check if static dir exists; if not, try relative to cwd
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		staticDir = filepath.Join(".", "..", "frontend", "out")
	}

	slog.Info("serving frontend", "dir", staticDir)

	// Serve static assets (JS, CSS, images, etc.)
	router.Use(func(c *gin.Context) {
		// Skip API routes
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Next()
			return
		}

		// Try to serve the file directly
		requestPath := c.Request.URL.Path
		filePath := filepath.Join(staticDir, requestPath)

		// If the path points to a file that exists, serve it
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			http.ServeFile(c.Writer, c.Request, filePath)
			c.Abort()
			return
		}

		// For directory paths or clean URLs, try serving index.html inside that dir
		// (Next.js static export with trailingSlash creates  /about/index.html)
		indexPath := filepath.Join(filePath, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			http.ServeFile(c.Writer, c.Request, indexPath)
			c.Abort()
			return
		}

		// For dynamic routes like /profile/username, try parent directories' index.html
		// Next.js static export creates /profile/index.html but not /profile/username/index.html
		parts := strings.Split(strings.Trim(requestPath, "/"), "/")
		for i := len(parts) - 1; i >= 1; i-- {
			parentIndex := filepath.Join(staticDir, strings.Join(parts[:i], "/"), "index.html")
			if _, err := os.Stat(parentIndex); err == nil {
				http.ServeFile(c.Writer, c.Request, parentIndex)
				c.Abort()
				return
			}
		}

		// Fallback: serve the root index.html (SPA-like behaviour)
		rootIndex := filepath.Join(staticDir, "index.html")
		if _, err := os.Stat(rootIndex); err == nil {
			http.ServeFile(c.Writer, c.Request, rootIndex)
			c.Abort()
			return
		}

		// Nothing found — let Gin return 404
		c.Next()
	})

	slog.Info("FastChem backend starting", "port", cfg.Port)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}
	slog.Info("server exited cleanly")
}
