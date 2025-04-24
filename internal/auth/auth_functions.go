package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrHashMismatch = errors.New("hash and password do not match")
var ErrTokenSigning = errors.New("error signing token")
var ErrNoAuthHeader = errors.New("no authorization header provided")
var ErrUnauthorized = errors.New("user not authorized")

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("unexpected error encountered when generating password hash: %w", err)
	}

	return string(hash), err
}

func CheckPasswordHash(hash, password string) error {
	// hashedPassword, err := HashPassword(password)
	// if err != nil {
	// 	return fmt.Errorf("error hashing password")
	// }

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return ErrHashMismatch
	}

	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	currTime := time.Now()
	currTimeJWT := jwt.NewNumericDate(currTime)

	expiresAt := currTime.Add(expiresIn)
	expiresAtJWT := jwt.NewNumericDate(expiresAt)

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  currTimeJWT,
		ExpiresAt: expiresAtJWT,
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", ErrTokenSigning
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if errors.Is(err, jwt.ErrTokenExpired) {
		return uuid.Nil, jwt.ErrTokenExpired
	} else if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	id, err := claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	idUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, err
	}

	return idUUID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeader
	}

	splitString := strings.Split(authHeader, "Bearer ")
	TOKEN_STRING := splitString

	return TOKEN_STRING[1], nil
}
