package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hashed, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %v", err)
	}

	return hashed, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	isValid, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("error comparing passwords: %v", err)
	}
	return isValid, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHead := headers.Get("Authorization")

	if authHead == "" {
		return "", fmt.Errorf("no authorization header")
	}

	after, found := strings.CutPrefix(authHead, "Bearer ")
	if !found {
		return "", fmt.Errorf("unknown authorization format")
	}

	trimmed := strings.TrimSpace(after)
	if trimmed == "" {
		return "", fmt.Errorf("unknown authorization format")
	}

	return trimmed, nil
}

func MakeRefreshToken() string {
	rands := make([]byte, 32)
	_, err := rand.Read(rands)
	if err != nil {
		return ""
	}

	tk := hex.EncodeToString(rands)
	return tk
}

func GetAPIKey(headers http.Header) (string, error) {
	authHead := headers.Get("Authorization")

	if authHead == "" {
		return "", fmt.Errorf("no authorization header")
	}

	after, found := strings.CutPrefix(authHead, "ApiKey ")
	if !found {
		return "", fmt.Errorf("unknown authorization format")
	}

	trimmed := strings.TrimSpace(after)
	if trimmed == "" {
		return "", fmt.Errorf("unknown authorization format")
	}

	return trimmed, nil
}
