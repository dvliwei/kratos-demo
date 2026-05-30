package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

type Manager struct {
	secret []byte
	ttl    time.Duration
}

type Claims struct {
	UserID uint64 `json:"user_id"`
	Name   string `json:"name"`
	Issued int64  `json:"iat"`
	Expiry int64  `json:"exp"`
}

func NewManager(secret string, ttl time.Duration) (*Manager, error) {
	if secret == "" {
		return nil, errors.New("jwt password is required")
	}
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return &Manager{secret: []byte(secret), ttl: ttl}, nil
}

func (m *Manager) GenerateToken(userID uint64, name string) (string, error) {
	now := time.Now()
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	claims := Claims{
		UserID: userID,
		Name:   name,
		Issued: now.Unix(),
		Expiry: now.Add(m.ttl).Unix(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	unsigned := base64.RawURLEncoding.EncodeToString(headerJSON) + "." +
		base64.RawURLEncoding.EncodeToString(claimsJSON)
	return unsigned + "." + m.sign(unsigned), nil
}

func (m *Manager) ValidateToken(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	unsigned := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(parts[2]), []byte(m.sign(unsigned))) {
		return nil, ErrInvalidToken
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	if claims.Expiry > 0 && time.Now().Unix() > claims.Expiry {
		return nil, ErrExpiredToken
	}
	return &claims, nil
}

func (m *Manager) sign(unsigned string) string {
	mac := hmac.New(sha256.New, m.secret)
	_, _ = mac.Write([]byte(unsigned))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func DurationFromSeconds(seconds int64) time.Duration {
	if seconds <= 0 {
		return 24 * time.Hour
	}
	return time.Duration(seconds) * time.Second
}
