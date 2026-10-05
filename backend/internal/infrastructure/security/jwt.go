package security

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
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
	now       func() time.Time
}

// This single API has one fixed audience; it is not an organization scope.
const accessAudience = "elabtrack-v2-api"
const accessPurpose = "access"

var errInvalidAccessToken = errors.New("invalid access token")

// NewJWTIssuer constructs a JWT issuer from config.
func NewJWTIssuer(cfg config.JWTConfig) *JWTIssuer {
	return &JWTIssuer{
		secret:    []byte(cfg.Secret),
		accessTTL: cfg.AccessTTL,
		issuer:    cfg.Issuer,
		now:       func() time.Time { return time.Now().UTC() },
	}
}

// accessClaims is the internal JWT claims structure.
type accessClaims struct {
	UserID  uuid.UUID `json:"uid"`
	Email   string    `json:"email"`
	Role    string    `json:"role"`
	Purpose string    `json:"purpose"`
	jwt.RegisteredClaims
}

// IssueAccessToken signs a short-lived JWT for the given user.
func (j *JWTIssuer) IssueAccessToken(ctx context.Context, u *domainuser.User) (string, error) {
	if u == nil || u.ID == uuid.Nil {
		return "", errInvalidAccessToken
	}
	now := j.now()
	claims := accessClaims{
		UserID:  u.ID,
		Email:   u.Email,
		Role:    string(u.Role),
		Purpose: accessPurpose,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Audience:  jwt.ClaimStrings{accessAudience},
			Subject:   u.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessTTL)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(j.secret)
	if err != nil {
		return "", errInvalidAccessToken
	}
	return signed, nil
}

// VerifyAccessToken requires the complete access-token contract. Claim role
// and email remain hints; current account state authorizes HTTP operations.
func (j *JWTIssuer) VerifyAccessToken(ctx context.Context, tokenStr string) (application.Claims, error) {
	claims := &accessClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errInvalidAccessToken
		}
		return j.secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer(j.issuer),
		jwt.WithAudience(accessAudience), jwt.WithExpirationRequired(),
		jwt.WithNotBeforeRequired(), jwt.WithIssuedAt(), jwt.WithTimeFunc(j.now), jwt.WithStrictDecoding())
	if err != nil || token == nil || !token.Valid {
		return application.Claims{}, errInvalidAccessToken
	}
	if claims.UserID == uuid.Nil || claims.Subject != claims.UserID.String() ||
		claims.Purpose != accessPurpose || claims.IssuedAt == nil ||
		!claims.ExpiresAt.After(claims.IssuedAt.Time) ||
		claims.NotBefore.Before(claims.IssuedAt.Time) || !claims.NotBefore.Before(claims.ExpiresAt.Time) {
		return application.Claims{}, errInvalidAccessToken
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
		return "", errors.New("refresh credential generation failed")
	}
	return hex.EncodeToString(b), nil
}
