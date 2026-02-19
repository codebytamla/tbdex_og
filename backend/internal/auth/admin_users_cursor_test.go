package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAdminUsersCursorRoundTrip(t *testing.T) {
	now := time.Date(2026, 2, 13, 12, 30, 0, 0, time.UTC)
	id := uuid.New()

	token, err := EncodeAdminUsersCursor(now, id)
	if err != nil {
		t.Fatalf("encode cursor: %v", err)
	}

	decoded, err := DecodeAdminUsersCursor(token)
	if err != nil {
		t.Fatalf("decode cursor: %v", err)
	}

	if !decoded.CreatedAt.Equal(now) {
		t.Fatalf("created_at mismatch: got %v want %v", decoded.CreatedAt, now)
	}
	if decoded.ID != id {
		t.Fatalf("id mismatch: got %v want %v", decoded.ID, id)
	}
}

func TestDecodeAdminUsersCursorInvalid(t *testing.T) {
	_, err := DecodeAdminUsersCursor("not-a-valid-cursor")
	if err == nil {
		t.Fatal("expected error for invalid cursor")
	}
	if !errors.Is(err, ErrInvalidCursor) {
		t.Fatalf("expected ErrInvalidCursor, got %v", err)
	}
}
