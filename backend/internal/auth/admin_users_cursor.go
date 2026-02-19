package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AdminUsersCursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
}

func EncodeAdminUsersCursor(createdAt time.Time, id uuid.UUID) (string, error) {
	if createdAt.IsZero() || id == uuid.Nil {
		return "", fmt.Errorf("%w: missing created_at or id", ErrInvalidCursor)
	}

	payload, err := json.Marshal(AdminUsersCursor{
		CreatedAt: createdAt.UTC(),
		ID:        id,
	})
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(payload), nil
}

func DecodeAdminUsersCursor(value string) (AdminUsersCursor, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return AdminUsersCursor{}, fmt.Errorf("%w: base64 decode failed", ErrInvalidCursor)
	}

	var cursor AdminUsersCursor
	if err := json.Unmarshal(decoded, &cursor); err != nil {
		return AdminUsersCursor{}, fmt.Errorf("%w: payload decode failed", ErrInvalidCursor)
	}

	if cursor.CreatedAt.IsZero() || cursor.ID == uuid.Nil {
		return AdminUsersCursor{}, fmt.Errorf("%w: incomplete payload", ErrInvalidCursor)
	}

	return cursor, nil
}
