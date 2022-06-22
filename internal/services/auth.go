package services

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/persistence/postgres"
	"time"
)

const (
	salt       = "skdgndkxboi42e143okbd"
	signingKey = "lskd4231kfsd"
	tokenTTL   = 12 * time.Hour
)

type TokenClaims struct {
	jwt.StandardClaims
	StaffID uuid.UUID `json:"staff_id"`
}

type AuthService struct {
	rep postgres.StaffAuth
	ctx context.Context
}

func NewAuthService(ctx context.Context, rep postgres.StaffAuth) *AuthService {
	return &AuthService{rep: rep, ctx: ctx}
}

func (s *AuthService) GenerateToken(email, password string) (string, error) {
	staff, err := s.rep.GetStaffAuth(s.ctx, email, generatePasswordHash(password))
	if err != nil {
		return "", err
	}
	if staff == nil {
		return "", fmt.Errorf("no such user")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		staff.ID,
	})
	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return uuid.UUID{}, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return uuid.UUID{}, errors.New("token claims are not of type TokenClaims")
	}

	return claims.StaffID, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
