package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

var secret = "supersecrethahaha"

// JWT parts
type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Exp      int64  `json:"exp"`
}

// base64 URL encode (without padding)
func base64URLEncode(data []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}

// base64 URL decode (with padding fix)
func base64URLDecode(s string) ([]byte, error) {
	if m := len(s) % 4; m != 0 {
		s += strings.Repeat("=", 4-m)
	}
	return base64.URLEncoding.DecodeString(s)
}

// create JWT
func CreateJWT(username, avatar, secret string, userID int, duration time.Duration) (string, error) {
	header := Header{"HS256", "JWT"}
	claims := Claims{
		UserID:   userID,
		Username: username,
		Avatar:   avatar,
		Exp:      time.Now().Add(duration).Unix(),
	}

	// JSON encode
	headerJSON, _ := json.Marshal(header)
	claimsJSON, _ := json.Marshal(claims)

	// Base64 encode
	hEnc := base64URLEncode(headerJSON)
	cEnc := base64URLEncode(claimsJSON)

	// Sign
	unsigned := hEnc + "." + cEnc
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(unsigned))
	signature := base64URLEncode(h.Sum(nil))

	return unsigned + "." + signature, nil
}

// verify JWT
func VerifyJWT(token, secret string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	unsigned := parts[0] + "." + parts[1]

	// verify signature
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(unsigned))
	expectedSig := base64URLEncode(h.Sum(nil))
	if parts[2] != expectedSig {
		return nil, fmt.Errorf("invalid signature")
	}

	// decode claims
	claimsJSON, err := base64URLDecode(parts[1])
	if err != nil {
		return nil, err
	}

	var claims Claims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, err
	}

	// check expiration
	if time.Now().Unix() > claims.Exp {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}
