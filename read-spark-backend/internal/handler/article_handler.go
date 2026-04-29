package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/readspark/backend/internal/domain"
	"github.com/readspark/backend/internal/service"
)

type ArticleHandler struct {
	articleService *service.ArticleService
}

func NewArticleHandler(articleService *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService: articleService}
}

func (h *ArticleHandler) GetDaily(c *gin.Context) {
	articles, err := h.articleService.GetDailyArticles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_ERROR", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"articles": articles})
}

func (h *ArticleHandler) GetArticle(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_ID", "message": "invalid article id"})
		return
	}

	userIDVal, exists := c.Get("userID")
	var userID uuid.UUID
	var isSubscribed bool
	if exists {
		userID = userIDVal.(uuid.UUID)
		isSubscribed = true
	}

	article, err := h.articleService.GetArticle(c.Request.Context(), id, userID, isSubscribed)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "ARTICLE_NOT_FOUND", "message": "article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_ERROR", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) ListArticles(c *gin.Context) {
	category := c.Query("category")
	difficulty := c.Query("difficulty")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var catPtr, diffPtr *string
	if category != "" {
		catPtr = &category
	}
	if difficulty != "" {
		diffPtr = &difficulty
	}

	articles, total, err := h.articleService.ListArticles(c.Request.Context(), catPtr, diffPtr, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_ERROR", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"articles":  articles,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
