package http

import (
	"go-hex/internal/core/port"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc port.AuthService
}

func NewAuthHandler(svc port.AuthService) *AuthHandler {
	return &AuthHandler{svc}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"alex@mail.com"`
	Password string `json:"password" binding:"required,min=8" example:"12345679" minLength:"8"`
}

func (ah *AuthHandler) Login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	token, err := ah.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
