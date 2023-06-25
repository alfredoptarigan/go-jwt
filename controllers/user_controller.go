package controllers

import (
	"errors"
	"github.com/alfredoptarigan/go-jwt/initializers"
	"github.com/alfredoptarigan/go-jwt/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

func Signup(c *gin.Context) {
	// Get fullname, email and password from the body
	var body struct {
		Fullname string `json:"fullname" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]models.APIError, len(ve))
			for i, err := range ve {
				out[i] = models.APIError{
					Field:   err.Field(),
					Message: err.Tag(),
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": out})
		}
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Create user
	user := models.User{
		Name:     body.Fullname,
		Email:    body.Email,
		Password: string(hash),
	}

	// Save user
	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "User created successfully!",
	})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]models.APIError, len(ve))
			for i, err := range ve {
				out[i] = models.APIError{
					Field:   err.Field(),
					Message: err.Tag(),
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": out})
		}
		return
	}

	// Look up the requested user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare sent in password with the hashed password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_TOKEN")))

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create Token"})
		return
	}

	// Set Some Cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	// Send it to API
	c.JSON(200, gin.H{
		"token": tokenString,
	})

}

func Validate(c *gin.Context) {
	user, _ := c.MustGet("user").(models.User)

	c.JSON(http.StatusOK, gin.H{"message": user})
}
