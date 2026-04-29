package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/readspark/backend/internal/config"
	"github.com/readspark/backend/internal/domain"
	"github.com/readspark/backend/internal/repository"
)

type AuthService struct {
	userRepo         *repository.UserRepository
	jwtCfg           config.JWTConfig
	verificationCode string
}

func NewAuthService(userRepo *repository.UserRepository, jwtCfg config.JWTConfig, authCfg config.AuthConfig) *AuthService {
	code := authCfg.VerificationCode
	if code == "" {
		code = "123456"
	}
	return &AuthService{
		userRepo:         userRepo,
		jwtCfg:           jwtCfg,
		verificationCode: code,
	}
}

func (s *AuthService) VerifyCode(ctx context.Context, phone, code string) error {
	_ = ctx
	_ = phone
	if code != s.verificationCode {
		return domain.ErrInvalidCode
	}
	return nil
}

func (s *AuthService) Register(ctx context.Context, req domain.UserRegisterRequest) (*domain.TokenPair, error) {
	if err := s.VerifyCode(ctx, req.Phone, req.Code); err != nil {
		return nil, err
	}

	existing, err := s.userRepo.FindByPhone(ctx, req.Phone)
	if err == nil && existing != nil {
		return nil, domain.ErrAlreadyExists
	}

	user := &domain.User{
		ID:    uuid.New(),
		Phone: req.Phone,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.generateTokens(user.ID)
}

func (s *AuthService) Login(ctx context.Context, req domain.UserLoginRequest) (*domain.TokenPair, error) {
	if err := s.VerifyCode(ctx, req.Phone, req.Code); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByPhone(ctx, req.Phone)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return s.Register(ctx, domain.UserRegisterRequest(req))
		}
		return nil, err
	}

	return s.generateTokens(user.ID)
}

func (s *AuthService) generateTokens(userID uuid.UUID) (*domain.TokenPair, error) {
	accessTTL, _ := time.ParseDuration("15m")
	refreshTTL, _ := time.ParseDuration("168h")

	accessClaims := jwt.MapClaims{
		"sub": userID.String(),
		"typ": "access",
		"exp": time.Now().Add(accessTTL).Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(s.jwtCfg.Secret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshClaims := jwt.MapClaims{
		"sub": userID.String(),
		"typ": "refresh",
		"exp": time.Now().Add(refreshTTL).Unix(),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(s.jwtCfg.Secret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(accessTTL.Seconds()),
	}, nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtCfg.Secret), nil
	})
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	if typ, ok := claims["typ"].(string); !ok || typ != "refresh" {
		return nil, domain.ErrInvalidToken
	}

	userIDStr, ok := claims["sub"].(string)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	if _, err := s.userRepo.FindByID(ctx, userID); err != nil {
		return nil, err
	}

	return s.generateTokens(userID)
}
