package controllers

import (
	"errors"
	"net/http"
	"strings"

	"backend/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (controller *AuthController) Register(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "name, email and password are required",
		})
		return
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.TrimSpace(strings.ToLower(request.Email))
	request.Password = strings.TrimSpace(request.Password)

	if request.Name == "" || request.Email == "" || request.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "name, email and password are required",
		})
		return
	}

	user, err := controller.authService.Register(services.RegisterInput{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"error":   "email already registered",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "could not register user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "user registered successfully",
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func (controller *AuthController) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "email and password are required",
		})
		return
	}

	request.Email = strings.TrimSpace(strings.ToLower(request.Email))
	request.Password = strings.TrimSpace(request.Password)

	if request.Email == "" || request.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "email and password are required",
		})
		return
	}

	login, err := controller.authService.Login(services.LoginInput{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "invalid credentials",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "could not login user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "login successful",
		"data": gin.H{
			"token": login.Token,
			"id":    login.User.ID,
			"name":  login.User.Name,
			"email": login.User.Email,
		},
	})
}
