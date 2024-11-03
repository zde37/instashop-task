package pkg

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zde37/instashop-task/internal/models"
)

// TokenType defines the type of token
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// CustomClaims extends jwt.RegisteredClaims to include custom fields
type CustomClaims struct {
	UserID    string          `json:"user_id"`
	UserEmail string          `json:"user_email"`
	UserRole  models.UserRole `json:"user_role"`
	TokenType TokenType       `json:"token_type"`
	jwt.RegisteredClaims
}

// JWTMaker is a struct that handles JWT operations
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker instance
func NewJWTMaker(secretKey string) (*JWTMaker, error) {
	if len(secretKey) < 32 {
		return nil, errors.New("secret key must be at least 32 characters long")
	}
	return &JWTMaker{secretKey: secretKey}, nil
}

// CreateToken creates a new token for a specific user and duration
func (maker *JWTMaker) CreateToken(
	userID string,
	userEmail string,
	userRole models.UserRole,
	tokenType TokenType,
	duration time.Duration,
) (string, error) {
	claims := &CustomClaims{
		UserID:    userID,
		UserEmail: userEmail,
		UserRole:  userRole,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "InstaShop",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(maker.secretKey))
}

// VerifyToken checks if the token is valid and returns the claims
func (maker *JWTMaker) VerifyToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(maker.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// RefreshToken creates a new access token if the refresh token is valid
func (maker *JWTMaker) RefreshToken(refreshToken string, accessTokenDuration time.Duration) (string, error) {
	claims, err := maker.VerifyToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("failed to verify refresh token: %w", err)
	}

	if claims.TokenType != RefreshToken {
		return "", errors.New("provided token is not a refresh token")
	}

	return maker.CreateToken(
		claims.UserID,
		claims.UserEmail,
		claims.UserRole,
		AccessToken,
		accessTokenDuration,
	)
}

// GetTokenExpiry extracts the expiry time from a token
func (maker *JWTMaker) GetTokenExpiry(tokenString string) (time.Time, error) {
	claims, err := maker.VerifyToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}

	if claims.ExpiresAt == nil {
		return time.Time{}, errors.New("token has no expiry time")
	}

	return claims.ExpiresAt.Time, nil
}
