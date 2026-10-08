package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Maaku050/elabtrack-v2/backend/internal/application"
	domainauth "github.com/Maaku050/elabtrack-v2/backend/internal/domain/auth"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/security"
	"github.com/google/uuid"
)

// This serial, copy-on-write test double models rollback/one-use semantics.
// Its mutex is only test infrastructure; live exclusion uses PostgreSQL locks.
type sessionStateKey struct{}
type sessionState map[string]*domainauth.RefreshToken
type sessionHarness struct {
	domainauth.Repository
	mu                                                    sync.Mutex
	sessions                                              sessionState
	account                                               *domainuser.Account
	now                                                   time.Time
	findErr, accountErr, createErr, consumeErr, commitErr error
	lookedUp, revoked                                     string
	operations                                            []string
	accountHook                                           func()
}

func cloneSessions(in sessionState) sessionState {
	out := sessionState{}
	for hash, session := range in {
		copy := *session
		out[hash] = &copy
	}
	return out
}
func (h *sessionHarness) Within(ctx context.Context, fn func(context.Context) error) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	state := cloneSessions(h.sessions)
	h.operations = append(h.operations, "begin")
	if err := fn(context.WithValue(ctx, sessionStateKey{}, state)); err != nil {
		h.operations = append(h.operations, "rollback")
		return err
	}
	if h.commitErr != nil {
		h.operations = append(h.operations, "rollback")
		return h.commitErr
	}
	h.sessions = state
	h.operations = append(h.operations, "commit")
	return nil
}
func (h *sessionHarness) FindByHashForUpdate(ctx context.Context, hash string) (*domainauth.RefreshToken, error) {
	state, ok := ctx.Value(sessionStateKey{}).(sessionState)
	if !ok {
		return nil, errors.New("transaction required")
	}
	h.operations = append(h.operations, "lock session")
	h.lookedUp = hash
	if h.findErr != nil {
		return nil, h.findErr
	}
	session := state[hash]
	if session == nil {
		return nil, domainauth.ErrTokenNotFound
	}
	return session, nil
}
func (h *sessionHarness) LockAccountByID(ctx context.Context, id uuid.UUID) (*domainuser.Account, error) {
	if ctx.Value(sessionStateKey{}) == nil {
		return nil, errors.New("transaction required")
	}
	h.operations = append(h.operations, "lock account")
	if h.accountHook != nil {
		h.accountHook()
	}
	if h.accountErr != nil {
		return nil, h.accountErr
	}
	if h.account == nil {
		return nil, domainuser.ErrUserNotFound
	}
	copy := *h.account
	return &copy, nil
}
func (h *sessionHarness) Create(ctx context.Context, session *domainauth.RefreshToken) error {
	state, ok := ctx.Value(sessionStateKey{}).(sessionState)
	if !ok {
		state = h.sessions
	}
	h.operations = append(h.operations, "insert")
	copy := *session
	state[session.TokenHash] = &copy
	return h.createErr
}
func (h *sessionHarness) Consume(ctx context.Context, hash string, replacement uuid.UUID, now time.Time) error {
	state, ok := ctx.Value(sessionStateKey{}).(sessionState)
	if !ok {
		return errors.New("transaction required")
	}
	h.operations = append(h.operations, "consume")
	session := state[hash]
	if session == nil || !session.ValidAt(now) {
		return domainauth.ErrTokenInvalid
	}
	session.RevokedAt = &now
	session.ReplacedBy = &replacement
	session.UpdatedAt = now
	return h.consumeErr
}
func (h *sessionHarness) Revoke(_ context.Context, hash string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.revoked = hash
	if session := h.sessions[hash]; session != nil && session.RevokedAt == nil {
		now := h.now
		session.RevokedAt = &now
	}
	return h.consumeErr
}

type pairIssuer struct {
	application.TokenIssuer
	next                 int
	issuedRole           domainuser.Role
	accessErr, randomErr error
}

func (i *pairIssuer) IssueAccessToken(_ context.Context, u *domainuser.User) (string, error) {
	i.issuedRole = u.Role
	return "synthetic-access", i.accessErr
}
func (i *pairIssuer) GenerateRefreshToken() (string, error) {
	i.next++
	return fmt.Sprintf("%064x", i.next), i.randomErr
}
func setupRefresh() (*Service, *sessionHarness, *pairIssuer, string) {
	now := time.Date(2026, 10, 6, 0, 0, 0, 0, time.UTC)
	id := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	raw := strings.Repeat("a", 64)
	hash, _ := (security.SHA256RefreshHasher{}).Hash(raw)
	h := &sessionHarness{sessions: sessionState{}, now: now, account: &domainuser.Account{ID: id, Email: "current@example.invalid", Role: domainuser.RoleBorrower, IsActive: true}}
	h.sessions[hash] = domainauth.NewRefreshToken(hash, id, now.Add(-time.Hour), now.Add(time.Hour))
	i := &pairIssuer{}
	s := NewService(nil, h, nil, i, security.SHA256RefreshHasher{}, h, h, time.Minute, time.Hour)
	s.now = func() time.Time { return h.now }
	return s, h, i, raw
}
func assertEmptyPair(t *testing.T, pair TokenPairDTO) {
	t.Helper()
	if pair != (TokenPairDTO{}) {
		t.Fatal("failed operation returned credentials")
	}
}
func TestRefreshRotationAndHashOnlyPersistence(t *testing.T) {
	s, h, i, raw := setupRefresh()
	pair, err := s.Refresh(context.Background(), RefreshRequest{RefreshToken: raw})
	if err != nil {
		t.Fatal(err)
	}
	hash, _ := (security.SHA256RefreshHasher{}).Hash(raw)
	if h.lookedUp != hash || h.lookedUp == raw {
		t.Fatal("lookup did not digest presented secret")
	}
	old := h.sessions[hash]
	replacementHash, _ := (security.SHA256RefreshHasher{}).Hash(pair.RefreshToken)
	replacement := h.sessions[replacementHash]
	if len(h.sessions) != 2 || old.RevokedAt == nil || old.ReplacedBy == nil || replacement == nil || *old.ReplacedBy != replacement.ID || replacement.TokenHash == pair.RefreshToken || replacement.UserID != h.account.ID {
		t.Fatal("rotation/persistence invariant failed")
	}
	if replacement.ExpiresAt != h.now.Add(time.Hour) || pair.ExpiresAt != h.now.Add(time.Minute) || pair.TokenType != "Bearer" || i.issuedRole != domainuser.RoleBorrower {
		t.Fatal("typed TTL/transport/current role contract")
	}
	if strings.Join(h.operations, ",") != "begin,lock session,lock account,insert,consume,commit" {
		t.Fatal("transaction did not encompass rotation")
	}
	denied, err := s.Refresh(context.Background(), RefreshRequest{RefreshToken: raw})
	assertEmptyPair(t, denied)
	if !errors.Is(err, domainauth.ErrTokenInvalid) || len(h.sessions) != 2 {
		t.Fatal("replay produced successor")
	}
	if _, err := s.Refresh(context.Background(), RefreshRequest{RefreshToken: pair.RefreshToken}); err != nil {
		t.Fatal("replacement not usable", err)
	}
	if len(h.sessions) != 3 {
		t.Fatal("replacement rotation missing")
	}
}
func TestRefreshRejectsInvalidSessionOrCurrentAccount(t *testing.T) {
	for _, name := range []string{"expired", "exact expiry", "revoked", "unknown", "malformed", "deleted account", "inactive account", "unknown role", "mismatched account", "expires while waiting"} {
		t.Run(name, func(t *testing.T) {
			s, h, _, raw := setupRefresh()
			hash, _ := (security.SHA256RefreshHasher{}).Hash(raw)
			switch name {
			case "expired":
				h.sessions[hash].ExpiresAt = h.now.Add(-time.Second)
			case "exact expiry":
				h.sessions[hash].ExpiresAt = h.now
			case "revoked":
				h.sessions[hash].RevokedAt = &h.now
			case "unknown":
				raw = strings.Repeat("b", 64)
			case "malformed":
				raw = "credential-sentinel"
			case "deleted account":
				h.account = nil
			case "inactive account":
				h.account.IsActive = false
			case "unknown role":
				h.account.Role = "superadmin"
			case "mismatched account":
				h.account.ID = uuid.New()
			case "expires while waiting":
				h.accountHook = func() { h.now = h.now.Add(2 * time.Hour) }
			}
			before := len(h.sessions)
			pair, err := s.Refresh(context.Background(), RefreshRequest{RefreshToken: raw})
			assertEmptyPair(t, pair)
			if !errors.Is(err, domainauth.ErrTokenInvalid) || len(h.sessions) != before {
				t.Fatal("invalid refresh changed sessions or returned wrong error")
			}
		})
	}
}
func TestRefreshRollbackAndSafeErrors(t *testing.T) {
	for _, name := range []string{"lookup", "account", "access issuance", "random generation", "replacement creation", "consumption", "commit"} {
		t.Run(name, func(t *testing.T) {
			s, h, i, raw := setupRefresh()
			hash, _ := (security.SHA256RefreshHasher{}).Hash(raw)
			sensitive := errors.New(raw + hash + " password-sentinel signing-secret Authorization")
			switch name {
			case "lookup":
				h.findErr = sensitive
			case "account":
				h.accountErr = sensitive
			case "access issuance":
				i.accessErr = sensitive
			case "random generation":
				i.randomErr = sensitive
			case "replacement creation":
				h.createErr = sensitive
			case "consumption":
				h.consumeErr = sensitive
			case "commit":
				h.commitErr = sensitive
			}
			pair, err := s.Refresh(context.Background(), RefreshRequest{RefreshToken: raw})
			assertEmptyPair(t, pair)
			if !errors.Is(err, shared.ErrInternal) || strings.Contains(err.Error(), raw) || strings.Contains(err.Error(), hash) {
				t.Fatal("unsafe/non-internal error")
			}
			if len(h.sessions) != 1 || h.sessions[hash].RevokedAt != nil || h.sessions[hash].ReplacedBy != nil {
				t.Fatal("failed transaction left partial rotation")
			}
			h.findErr = nil
			h.accountErr = nil
			h.createErr = nil
			h.consumeErr = nil
			h.commitErr = nil
			i.accessErr = nil
			i.randomErr = nil
			if _, err = s.Refresh(context.Background(), RefreshRequest{RefreshToken: raw}); err != nil {
				t.Fatal("old credential stranded after rollback")
			}
		})
	}
}
func TestConcurrentRefreshAtMostOneSuccessWithTransactionDouble(t *testing.T) {
	s, h, _, raw := setupRefresh()
	start := make(chan struct{})
	var wg sync.WaitGroup
	var successes atomic.Int32
	var failures atomic.Int32
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			pair, err := s.Refresh(context.Background(), RefreshRequest{RefreshToken: raw})
			if err == nil {
				successes.Add(1)
			} else if errors.Is(err, domainauth.ErrTokenInvalid) && pair == (TokenPairDTO{}) {
				failures.Add(1)
			}
		}()
	}
	close(start)
	wg.Wait()
	if successes.Load() != 1 || failures.Load() != 15 || len(h.sessions) != 2 {
		t.Fatal("more than one successor or wrong losing behavior")
	}
}
func TestLogoutIdempotentAndHashed(t *testing.T) {
	s, h, _, raw := setupRefresh()
	hash, _ := (security.SHA256RefreshHasher{}).Hash(raw)
	for range 2 {
		if err := s.Logout(context.Background(), RefreshRequest{RefreshToken: raw}); err != nil {
			t.Fatal(err)
		}
	}
	if h.revoked != hash || h.sessions[hash].RevokedAt == nil {
		t.Fatal("logout did not hash/revoke")
	}
	pair, err := s.Refresh(context.Background(), RefreshRequest{RefreshToken: raw})
	assertEmptyPair(t, pair)
	if !errors.Is(err, domainauth.ErrTokenInvalid) {
		t.Fatal("revoked session refreshed")
	}
	for _, unknown := range []string{strings.Repeat("b", 64), "malformed-sentinel"} {
		if err := s.Logout(context.Background(), RefreshRequest{RefreshToken: unknown}); err != nil {
			t.Fatal("logout exposed token existence")
		}
	}
	h.consumeErr = errors.New(raw + hash)
	if err := s.Logout(context.Background(), RefreshRequest{RefreshToken: raw}); !errors.Is(err, shared.ErrInternal) {
		t.Fatal("logout leaked backend error")
	}
}
func TestInitialIssuancePersistsDigestOnly(t *testing.T) {
	s, h, _, _ := setupRefresh()
	u := &domainuser.User{ID: h.account.ID, Email: h.account.Email, Role: h.account.Role}
	pair, err := s.issueTokenPair(context.Background(), u)
	if err != nil {
		t.Fatal(err)
	}
	hash, _ := (security.SHA256RefreshHasher{}).Hash(pair.RefreshToken)
	session := h.sessions[hash]
	if session == nil || session.TokenHash == pair.RefreshToken {
		t.Fatal("initial issuance stored raw secret")
	}
	h.createErr = errors.New(pair.RefreshToken + hash)
	failed, err := s.issueTokenPair(context.Background(), u)
	assertEmptyPair(t, failed)
	if !errors.Is(err, shared.ErrInternal) {
		t.Fatal("initial session error exposed detail")
	}
}

func TestStoredDigestIsNotRefreshCredential(t *testing.T) {
	s, h, _, raw := setupRefresh()
	hash, _ := (security.SHA256RefreshHasher{}).Hash(raw)
	pair, err := s.Refresh(context.Background(), RefreshRequest{RefreshToken: hash})
	assertEmptyPair(t, pair)
	if !errors.Is(err, domainauth.ErrTokenInvalid) || len(h.sessions) != 1 {
		t.Fatal("stored digest accepted as bearer")
	}
}
