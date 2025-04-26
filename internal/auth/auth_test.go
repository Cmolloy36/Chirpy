package auth

import (
	"errors"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHashSuccess(t *testing.T) {
	inputString := "hello"
	hash, err := HashPassword(inputString)
	if err != nil {
		t.Fatalf("error hashing password: %v", err)
	}

	if err = CheckPasswordHash(hash, inputString); err != nil {
		t.Errorf("error comparing hash and password: %v", err)
	}
}

func TestHashFailure(t *testing.T) {
	inputString := "hello"
	incorrectInputString := "olleh"
	hash, err := HashPassword(inputString)
	if err != nil {
		t.Fatalf("error hashing password: %v", err)
	}

	if err = CheckPasswordHash(hash, incorrectInputString); !errors.Is(err, ErrHashMismatch) {
		t.Errorf("password was not rejected appropriately: %v", err)
	}
}

func TestMakeJWT(t *testing.T) {
	userId := uuid.New()
	tokenSecret := os.Getenv("TOKEN_SECRET")

	signedToken, err := MakeJWT(userId, tokenSecret)
	if err != nil {
		t.Fatalf("error signing token: %v", err)
	}

	userIdValidated, err := ValidateJWT(signedToken, tokenSecret)
	if err != nil {
		t.Fatalf("error validating token: %v", err)
	}

	assert.Equal(t, userId, userIdValidated)

}

func TestExpiredToken(t *testing.T) {
	userId := uuid.New()
	tokenSecret := os.Getenv("TOKEN_SECRET")

	signedToken, err := MakeJWT(userId, tokenSecret)
	if err != nil {
		t.Fatalf("error signing token: %v", err)
	}

	userIdValidated, err := ValidateJWT(signedToken, tokenSecret)
	if !errors.Is(err, jwt.ErrTokenExpired) {
		t.Fatalf("token should be expired, but isn't")
	}

	assert.Equal(t, userIdValidated, uuid.Nil)
}

func TestInvalidSecret(t *testing.T) {
	// Create a token with a specific secret
	userId := uuid.New()
	correctSecret := "correct-secret"

	signedToken, err := MakeJWT(userId, correctSecret)
	if err != nil {
		t.Fatalf("error signing token: %v", err)
	}

	// Try to validate with a different secret
	invalidTokenSecret := os.Getenv("INVALID_TOKEN_SECRET")
	_, err = ValidateJWT(signedToken, invalidTokenSecret)

	// Assert that there was an error
	assert.Error(t, err)
}

func TestInvalidSecret2(t *testing.T) {
	userId := uuid.New()
	tokenSecret := "right_secret"

	signedToken, err := MakeJWT(userId, tokenSecret)
	if err != nil {
		t.Fatalf("error signing token: %v", err)
	}

	invalidTokenSecret := "wrong_secret"
	_, err = ValidateJWT(signedToken, invalidTokenSecret)

	assert.Error(t, err)
}

func TestGetAuthHeader(t *testing.T) {

}
