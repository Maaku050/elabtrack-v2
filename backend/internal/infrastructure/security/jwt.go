package security

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/Maaku050/elabtrack-v2/backend/internal/application"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTIssuer implements application.TokenIssuer using HS256 JWTs for access
// tokens and cryptographically-random opaque strings for refresh tokens.
type JWTIssuer struct {
	secret    []byte
	accessTTL time.Duration
	issuer    string
}

// NewJWTIssuer constructs a JWT issuer from config.
func NewJWTIssuer(cfg config.JWTConfig) *JWTIssuer {
	return &JWTIssuer{
		secret:    []byte(cfg.Secret),
		accessTTL: cfg.AccessTTL,
		issuer:    cfg.Issuer,
	}
}

// accessClaims is the internal JWT claims structure.
type accessClaims struct {
	UserID uuid.UUID `json:"uid"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

// IssueAccessToken signs a short-lived JWT for the given user.
func (j *JWTIssuer) IssueAccessToken(ctx context.Context, u *domainuser.User) (string, error) {
	now := time.Now().UTC()
	claims := accessClaims{
		UserID: u.ID,
		Email:  u.Email,
		Role:   string(u.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   u.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessTTL)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}
	return signed, nil
}

// VerifyAccessToken validates the signature and expiry of a JWT and returns
// the application-level claims.
func (j *JWTIssuer) VerifyAccessToken(ctx context.Context, tokenStr string) (application.Claims, error) {
	claims := &accessClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return j.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return application.Claims{}, jwt.ErrTokenExpired
		}
		return application.Claims{}, fmt.Errorf("invalid token: %w", err)
	}
	if !token.Valid {
		return application.Claims{}, errors.New("invalid token")
	}
	return application.Claims{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
	}, nil
}

// GenerateRefreshToken returns a 32-byte cryptographically-random opaque token
// encoded as 64 hex characters.
func (j *JWTIssuer) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}
	return hex.EncodeToString(b), nil
}
