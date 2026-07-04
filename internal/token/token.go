package token

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var jwtSecret = []byte("your-secret-key-change-in-production")

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenStore struct {
	rdb *redis.Client
}

func NewTokenStore(rdb *redis.Client) *TokenStore {
	return &TokenStore{rdb: rdb}
}

func (s *TokenStore) GenerateRefreshToken(ctx context.Context, userID int) (string, error) {
	token := uuid.New().String()
	key := fmt.Sprintf("refresh:token:%s", token)
	err := s.rdb.Set(ctx, key, userID, 7*24*time.Hour).Err()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *TokenStore) GetRefreshTokenUserID(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("refresh:token:%s", token)
	return s.rdb.Get(ctx, key).Result()
}

func (s *TokenStore) DeleteRefreshToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("refresh:token:%s", token)
	return s.rdb.Del(ctx, key).Err()
}

func (s *TokenStore) GenerateVerificationToken(ctx context.Context, userID int) (string, error) {
	token := uuid.New().String()
	key := fmt.Sprintf("verify:token:%s", token)
	err := s.rdb.Set(ctx, key, userID, 24*time.Hour).Err()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *TokenStore) GetVerificationTokenUserID(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("verify:token:%s", token)
	return s.rdb.Get(ctx, key).Result()
}

func (s *TokenStore) DeleteVerificationToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("verify:token:%s", token)
	return s.rdb.Del(ctx, key).Err()
}

func GenerateAccessToken(userID int) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(jwtSecret)
}

func ParseAccessToken(tokenStr string) (*Claims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
