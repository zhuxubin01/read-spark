package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/readspark/backend/internal/domain"
	"github.com/readspark/backend/internal/middleware"
	"github.com/readspark/backend/internal/repository"
	"github.com/readspark/backend/internal/service"
)

func TestProgressSync_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewProgressHandler(nil)
	r.POST("/progress", h.Sync)

	req := httptest.NewRequest(http.MethodPost, "/progress", strings.NewReader(`{"position":1}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestProgressList_MissingUserContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewProgressHandler(nil)
	r.GET("/progress", h.List)

	req := httptest.NewRequest(http.MethodGet, "/progress", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestProgressList_WithJWT_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open("file:"+uuid.NewString()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&domain.ReadingProgress{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	userID := uuid.New()
	_ = db.Create(&domain.ReadingProgress{ID: uuid.New(), UserID: userID, ArticleID: uuid.New(), Position: 10, Percentage: 20, LastReadAt: time.Now()}).Error

	repo := repository.NewProgressRepository(db)
	svc := service.NewProgressService(repo)
	h := NewProgressHandler(svc)

	secret := "test-secret"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": userID.String(), "exp": time.Now().Add(time.Hour).Unix()})
	tok, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	r := gin.New()
	g := r.Group("/")
	g.Use(middleware.JWTAuth(secret))
	g.GET("/progress", h.List)

	req := httptest.NewRequest(http.MethodGet, "/progress", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}
