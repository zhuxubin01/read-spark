package domain

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidCode     = errors.New("invalid verification code")
	ErrInvalidToken    = errors.New("invalid token")
	ErrTokenExpired    = errors.New("token expired")
	ErrArticleNotFound = errors.New("article not found")
	ErrAlreadyExists   = errors.New("resource already exists")
)
