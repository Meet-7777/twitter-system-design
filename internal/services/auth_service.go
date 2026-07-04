package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"twitter-system-design/internal/models"
	"twitter-system-design/internal/repository"
	"twitter-system-design/internal/token"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo   *repository.UserRepository
	tokenStore *token.TokenStore
}

func NewAuthService(userRepo *repository.UserRepository, tokenStore *token.TokenStore) *AuthService {
	return &AuthService{userRepo: userRepo, tokenStore: tokenStore}
}

func (s *AuthService) Signup(ctx context.Context, username, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user, err := s.userRepo.Create(username, email, string(hash))
	if err != nil {
		return err
	}
	verificationToken, err := s.tokenStore.GenerateVerificationToken(ctx, user.ID)
	if err != nil {
		return err
	}
	verificationLink := fmt.Sprintf("http://localhost:8080/verify?token=%s", verificationToken)
	log.Printf("VERIFY EMAIL for %s -> %s", email, verificationLink)
	return nil
}

func (s *AuthService) VerifyEmail(ctx context.Context, verificationToken string) error {
	userIDStr, err := s.tokenStore.GetVerificationTokenUserID(ctx, verificationToken)
	if err != nil {
		return err
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return err
	}
	if err := s.userRepo.UpdateVerified(userID); err != nil {
		return err
	}
	return s.tokenStore.DeleteVerificationToken(ctx, verificationToken)

}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *models.User `json:"user"`
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")

	}
	if !user.IsVerified {
		return nil, errors.New("user id is not verified")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	accessToken, err := token.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.tokenStore.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil

}
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (string, error) {
	userIDStr, err := s.tokenStore.GetRefreshTokenUserID(ctx, refreshToken)
	if err != nil {
		return "", errors.New("invalid or expired refresh token")
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return "", err
	}

	return token.GenerateAccessToken(userID)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.tokenStore.DeleteRefreshToken(ctx, refreshToken)
}
