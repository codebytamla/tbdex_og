package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashToken returns the SHA-256 hex digest of an opaque token.
// We use SHA-256 (not argon2) because opaque tokens are high-entropy
// random values, not human-chosen passwords. SHA-256 is sufficient to
// prevent storage-leak attacks while being fast to compute.
func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
