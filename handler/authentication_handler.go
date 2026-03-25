package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snck/book-keeper-api/service"
)

type AuthenticationHandler struct {
	service *service.AuthenticationService
}

func NewAuthenticationHandler(service *service.AuthenticationService) *AuthenticationHandler {
	return &AuthenticationHandler{service: service}
}

func (h *AuthenticationHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	newUser, err := h.service.AddNewUser(req.UserName, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error with database"})
		return
	}

	if newUser == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "user name already exist"})
		return
	}

	res := SignupResponse{
		ID:       newUser.ID,
		UserName: newUser.UserName,
	}

	c.JSON(http.StatusCreated, res)
}

func (h *AuthenticationHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	user, err := h.service.GetUser(req.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error with database"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "invalid user name or password"})
		return
	}

	if !h.service.IsPasswordValid(user.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid user name or password"})
		return
	}

	token, err := h.service.GenerateToken(*user)
	if err != nil {
		fmt.Println("error generating token", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error generating token"})
		return
	}

	res := LoginResponse{
		UserID:   user.ID,
		UserName: user.UserName,
		Token:    token,
	}

	c.JSON(http.StatusOK, res)
}
