package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tristnaja/bloc_test_backend/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

type AuthController struct {
	DB *gorm.DB
}

func (ac *AuthController) Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user := model.User{Username: input.Username, Email: input.Email, Password: string(hashedPassword)}

	if err := ac.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func (ac *AuthController) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	c.ShouldBindJSON(&input)

	var user model.User
	if err := ac.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, _ := GenerateToken(user.Email)
	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

func (ac *AuthController) GetProfile(c *gin.Context) {
	email, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Context lost"})
		return
	}

	var user model.User
	if err := ac.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"email":    user.Email,
		"joined":   user.CreatedAt,
	})
}
