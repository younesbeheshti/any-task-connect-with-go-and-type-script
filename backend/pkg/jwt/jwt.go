package jwt

import (
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/configs"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// Claims holds JWT payload data.
type Claims struct {
	UserID    string      `json:"sub"`
	Role      common.Role `json:"role"`
	TokenType string      `json:"typ"`
	JTI       string      `json:"jti"`
	jwtlib.RegisteredClaims
}

// Service handles token generation and validation.
type Service struct {
	secret          []byte
	accessDuration  time.Duration
	refreshDuration time.Duration
}

// NewService creates a JWT service from configuration.
func NewService(cfg configs.JWTConfig) *Service {
	return &Service{
		secret:          []byte(cfg.Secret),
		accessDuration:  cfg.AccessTokenDuration(),
		refreshDuration: cfg.RefreshTokenDuration(),
	}
}

// GenerateAccessToken creates a short-lived access token.
func (s *Service) GenerateAccessToken(userID string, role common.Role) (string, error) {
	return s.generateToken(userID, role, TokenTypeAccess, s.accessDuration)
}

// GenerateRefreshToken creates a long-lived refresh token.
func (s *Service) GenerateRefreshToken(userID string, role common.Role) (string, error) {
	return s.generateToken(userID, role, TokenTypeRefresh, s.refreshDuration)
}

func (s *Service) generateToken(userID string, role common.Role, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		UserID:    userID,
		Role:      role,
		TokenType: tokenType,
		JTI:       uuid.NewString(),
		RegisteredClaims: jwtlib.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwtlib.NewNumericDate(now),
			ExpiresAt: jwtlib.NewNumericDate(now.Add(ttl)),
			ID:        uuid.NewString(),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}

// ValidateToken parses and validates a token string.
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	return s.ParseClaims(tokenString)
}

// ParseClaims extracts claims from a token without type checking.
func (s *Service) ParseClaims(tokenString string) (*Claims, error) {
	token, err := jwtlib.ParseWithClaims(tokenString, &Claims{}, func(t *jwtlib.Token) (any, error) {
		if _, ok := t.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

// RefreshToken validates a refresh token and issues a new access token.
func (s *Service) RefreshToken(refreshToken string) (accessToken string, newRefreshToken string, err error) {
	claims, err := s.ParseClaims(refreshToken)
	if err != nil {
		return "", "", err
	}
	if claims.TokenType != TokenTypeRefresh {
		return "", "", fmt.Errorf("token is not a refresh token")
	}

	accessToken, err = s.GenerateAccessToken(claims.UserID, claims.Role)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err = s.GenerateRefreshToken(claims.UserID, claims.Role)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

// ValidateAccessToken ensures the token is a valid access token.
func (s *Service) ValidateAccessToken(tokenString string) (*Claims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != TokenTypeAccess {
		return nil, fmt.Errorf("token is not an access token")
	}
	return claims, nil
}
