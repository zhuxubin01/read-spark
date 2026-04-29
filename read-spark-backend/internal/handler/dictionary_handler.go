package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/readspark/backend/internal/service"
)

type DictionaryHandler struct {
	dictionaryService *service.DictionaryService
}

func NewDictionaryHandler(dictionaryService *service.DictionaryService) *DictionaryHandler {
	return &DictionaryHandler{dictionaryService: dictionaryService}
}

func (h *DictionaryHandler) Lookup(c *gin.Context) {
	word := c.Param("word")
	result, err := h.dictionaryService.Lookup(c.Request.Context(), word)
	if err != nil {
		switch err.Error() {
		case "word is required":
			c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": "word is required"})
		case "word not found":
			c.JSON(http.StatusNotFound, gin.H{"code": "WORD_NOT_FOUND", "message": "word not found"})
		default:
			c.JSON(http.StatusBadGateway, gin.H{"code": "DICTIONARY_UPSTREAM_ERROR", "message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, result)
}
