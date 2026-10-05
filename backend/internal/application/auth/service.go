package auth

import (
	"time"

	"github.com/Maaku050/elabtrack-v2/backend/internal/application"
	domainauth "github.com/Maaku050/elabtrack-v2/backend/internal/domain/auth"
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
)

// Service is the auth application service. It coordinates user creation,
// credential verification, and token issuance/rotation through ports.
//
// Use-case methods are split across register.go, login.go, and refresh.go
// to keep each file focused on a single use case.
type Service struct {
	users       domainuser.Repository
	tokens      domainauth.Repository
	hasher      application.PasswordHasher
	issuer      application.TokenIssuer
	tokenHasher application.RefreshTokenHasher
	tx          application.Transactions
	accounts    application.RefreshAccounts
	now         func() time.Time
	accessTTL   time.Duration
	refreshTTL  time.Duration
}

// NewService constructs an auth application service.
func NewService(
	users domainuser.Repository,
	tokens domainauth.Repository,
	hasher application.PasswordHasher,
	issuer application.TokenIssuer,
	tokenHasher application.RefreshTokenHasher,
	tx application.Transactions,
	accounts application.RefreshAccounts,
	accessTTL, refreshTTL time.Duration,
) *Service {
	return &Service{
		users:       users,
		tokens:      tokens,
		hasher:      hasher,
		issuer:      issuer,
		tokenHasher: tokenHasher,
		tx:          tx,
		accounts:    accounts,
		now:         func() time.Time { return time.Now().UTC() },
		accessTTL:   accessTTL,
		refreshTTL:  refreshTTL,
	}
}
