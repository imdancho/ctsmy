package main

import (
	"net/http"
	"time"

	"github.com/ZhangeldyB/CTSDerailmentAnalysis/backend/internal/authorization"
	"github.com/ZhangeldyB/CTSDerailmentAnalysis/backend/internal/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	db.InitDB()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/api/login", func(c *gin.Context) {
		var requestData struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var user db.User
		if err := db.DB.Where("email = ?", requestData.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		if err := authorization.VerifyPassword(user.Password, requestData.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		token, refreshToken, err := authorization.GenerateTokens(user.Id, user.Email, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  token,
			"refresh_token": refreshToken,
		})
	})

	r.POST("/api/signup", func(c *gin.Context) {
		var requestData struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		var existingUser db.User
		if err := db.DB.Where("email = ?", requestData.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestData.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user := db.User{
			Email:    requestData.Email,
			Password: string(hashedPassword),
		}

		if err := db.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	})
	r.GET("/api/protected", authorization.AuthMiddleware(), func(c *gin.Context) {
		userID := c.GetInt("user_id")
		email := c.GetString("email")
		role := c.GetString("role")

		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
			"email":   email,
			"role":    role,
		})
	})
	r.POST("/api/refresh", func(c *gin.Context) {
		var requestData struct {
			RefreshToken string `json:"refresh_token"`
		}

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		token, err := jwt.Parse(requestData.RefreshToken, func(token *jwt.Token) (interface{}, error) {
			return authorization.JwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims in refresh token"})
			return
		}

		userID, _ := claims["user_id"].(float64)
		email, _ := claims["email"].(string)
		role, _ := claims["role"].(string)

		accessToken, refreshToken, err := authorization.GenerateTokens(int(userID), email, role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	})

	r.Run("localhost:8080")
}
