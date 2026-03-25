package handler

import (
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
