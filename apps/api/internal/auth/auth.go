package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Secret is read from env in production; kept here as a default for dev only.
const defaultSecret = "dev-secret-change-me"

// SecretKey returns the signing secret from env (JWT_SECRET) or a dev default.
func SecretKey() string {
	if v := os.Getenv("JWT_SECRET"); v != "" {
		return v
	}
	return defaultSecret
}

// HashPassword hashes a plaintext password using bcrypt.
func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// CheckPassword verifies a plaintext password against a bcrypt hash.
func CheckPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

// Claims is the JWT payload.
type Claims struct {
	UserID string `json:"sub"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken mints a signed JWT for the given user.
func GenerateToken(userID, email, secret string, ttl time.Duration) (string, error) {
	if secret == "" {
		secret = defaultSecret
	}
	if ttl == 0 {
		ttl = 24 * time.Hour
	}
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString([]byte(secret))
}

// ParseToken validates and parses a JWT, returning its claims.
func ParseToken(tokenStr, secret string) (*Claims, error) {
	if secret == "" {
		secret = defaultSecret
	}
	tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
