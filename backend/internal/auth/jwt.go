package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
)

// Claims are the JWT claims for access tokens.
type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"uid"`
}

// JWTConfig holds JWT signing configuration.
type JWTConfig struct {
	SigningKey     []byte        // HMAC-SHA256 secret (>= 32 bytes)
	Issuer        string        // e.g. "borderform"
	Audience      string        // e.g. "borderform-api"
	AccessTTL     time.Duration // e.g. 15 * time.Minute
	RefreshTTL    time.Duration // e.g. 7 * 24 * time.Hour
	ClockSkew     time.Duration // e.g. 30 * time.Second
}

// IssueAccessToken creates a signed JWT for the given user.
func IssueAccessToken(cfg JWTConfig, userID uuid.UUID) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			Audience:  jwt.ClaimStrings{cfg.Audience},
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.AccessTTL)),
			ID:        uuid.New().String(), // jti for uniqueness
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(cfg.SigningKey)
}

// VerifyAccessToken parses and validates a JWT, returning the claims.
func VerifyAccessToken(cfg JWTConfig, tokenString string) (*Claims, error) {
	parser := jwt.NewParser(
		jwt.WithIssuer(cfg.Issuer),
		jwt.WithAudience(cfg.Audience),
		jwt.WithLeeway(cfg.ClockSkew),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	token, err := parser.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		return cfg.SigningKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GenerateOpaqueToken generates a cryptographically random opaque token
// of the given byte length, returned as a hex string.
func GenerateOpaqueToken(byteLen int) (string, error) {
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return hex.EncodeToString(b), nil
}
