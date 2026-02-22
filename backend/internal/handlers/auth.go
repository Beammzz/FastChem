package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/takumi/fastchem/internal/database"
	"github.com/takumi/fastchem/internal/middleware"
	"github.com/takumi/fastchem/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct{}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// Register handles POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username (3-20 chars) and password (6+ chars) are required"})
		return
	}

	// Normalize username
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username cannot be empty"})
		return
	}

	// Check if username already exists
	var exists int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE LOWER(username) = LOWER(?)", req.Username).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Insert user
	result, err := database.DB.Exec(
		"INSERT INTO users (username, password_hash) VALUES (?, ?)",
		req.Username, string(hash),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	userID, _ := result.LastInsertId()

	// Generate token
	token, err := middleware.GenerateToken(userID, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		Token: token,
		User: models.User{
			ID:       userID,
			Username: req.Username,
		},
	})
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	var user models.User
	err := database.DB.QueryRow(
		"SELECT id, username, password_hash, total_points, created_at FROM users WHERE LOWER(username) = LOWER(?)",
		req.Username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.TotalPoints, &user.CreatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate token
	token, err := middleware.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// Me handles GET /api/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetInt64("user_id")
	username := c.GetString("username")

	// Fetch user stats
	var stats models.UserStats

	// Get total points from user
	database.DB.QueryRow("SELECT total_points FROM users WHERE id = ?", userID).Scan(&stats.TotalPoints)

	err := database.DB.QueryRow(`
		SELECT 
			COUNT(*) as total_games,
			COALESCE(MAX(score), 0) as high_score,
			COALESCE(SUM(correct_answers), 0) as total_correct,
			COALESCE(SUM(total_answered), 0) as total_answered,
			COALESCE(AVG(score), 0) as average_score
		FROM scores WHERE user_id = ?
	`, userID).Scan(&stats.TotalGames, &stats.HighScore, &stats.TotalCorrect, &stats.TotalAnswered, &stats.AverageScore)

	if err != nil {
		stats = models.UserStats{}
	}

	if stats.TotalAnswered > 0 {
		stats.AverageAccuracy = float64(stats.TotalCorrect) / float64(stats.TotalAnswered) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"user": models.UserPublic{
			ID:          userID,
			Username:    username,
			TotalPoints: stats.TotalPoints,
		},
		"stats": stats,
	})
}
