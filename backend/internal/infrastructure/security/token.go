package security

// token.go documents the token design and re-exports the issuer for wiring.
//
// Token design:
//   - Access token:  JWT (HS256), short TTL (default 15m).
//                    Sent by clients via "Authorization: Bearer <token>".
//   - Refresh token: opaque 32-byte random string, persisted server-side
//                    with rotation. Long TTL (default 7d). Designed to be
//                    delivered via HttpOnly + Secure + SameSite cookies in
//                    production.
//
// Rotation strategy (see application/auth/refresh.go):
//   1. Client presents refresh token.
//   2. Server validates + revokes it.
//   3. Server issues a new access + refresh token pair.
//
// This file is intentionally minimal; the implementation lives in jwt.go.
