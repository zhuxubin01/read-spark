package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/readspark/backend/internal/domain"
	"github.com/readspark/backend/internal/service"
)

type AnnotationHandler struct {
	annotationService *service.AnnotationService
}

func NewAnnotationHandler(annotationService *service.AnnotationService) *AnnotationHandler {
	return &AnnotationHandler{annotationService: annotationService}
}

func (h *AnnotationHandler) Create(c *gin.Context) {
	var req domain.CreateAnnotationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": err.Error()})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "UNAUTHORIZED", "message": "missing user context"})
		return
	}

	row, err := h.annotationService.Create(c.Request.Context(), userID.(uuid.UUID), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, row)
}

func (h *AnnotationHandler) List(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "UNAUTHORIZED", "message": "missing user context"})
		return
	}

	var articleID *uuid.UUID
	if raw := c.Query("article_id"); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": "invalid article_id"})
			return
		}
		articleID = &id
	}

	rows, err := h.annotationService.List(c.Request.Context(), userID.(uuid.UUID), articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_ERROR", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"annotations": rows})
}
