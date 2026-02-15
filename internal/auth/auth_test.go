package internal

import (
	"testing"
	"time"

	// adjust based on your package name

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	secret := "supersecret"
	userID := uuid.New()
	expiresIn := time.Minute * 5

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	// ✅ Valid token
	parsedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}
	if parsedID != userID {
		t.Fatalf("expected %v, got %v", userID, parsedID)
	}

	// ❌ Expired token
	expiredToken, _ := MakeJWT(userID, secret, -time.Minute)
	if _, err := ValidateJWT(expiredToken, secret); err == nil {
		t.Fatalf("expected expired token to fail")
	}

	// ❌ Wrong secret
	if _, err := ValidateJWT(token, "wrongsecret"); err == nil {
		t.Fatalf("expected token with wrong secret to fail")
	}
}
