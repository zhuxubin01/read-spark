package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/readspark/backend/internal/domain"
	"github.com/readspark/backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": err.Error()})
		return
	}

	tokens, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		if err == domain.ErrAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"code": "USER_EXISTS", "message": "user already exists"})
			return
		}
		if err == domain.ErrInvalidCode {
			c.JSON(http.StatusUnauthorized, gin.H{"code": "INVALID_CODE", "message": "invalid verification code"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_ERROR", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": err.Error()})
		return
	}

	tokens, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if err == domain.ErrInvalidCode {
			c.JSON(http.StatusUnauthorized, gin.H{"code": "INVALID_CODE", "message": "invalid verification code"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_ERROR", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": err.Error()})
		return
	}

	tokens, err := h.authService.RefreshTokens(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "INVALID_TOKEN", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}
