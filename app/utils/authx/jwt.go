package authx

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
	ErrInvalidTokenFormat = errors.New("invalid token format")
	ErrInvalidTokenSign   = errors.New("invalid token signature")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidTokenClaims = errors.New("invalid token claims")
	ErrInvalidTokenType   = errors.New("invalid token type")
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type MemberTokenClaims struct {
	MemberID  string `json:"member_id"`
	Username  string `json:"username"`
	TokenType string `json:"token_type,omitempty"`
	Exp       int64  `json:"exp"`
	Iat       int64  `json:"iat"`
}

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

func GenerateMemberToken(secret string, memberID string, username string, ttl time.Duration) (string, int64, error) {
	return generateMemberToken(secret, memberID, username, TokenTypeAccess, ttl)
}

func GenerateMemberRefreshToken(secret string, memberID string, username string, ttl time.Duration) (string, int64, error) {
	return generateMemberToken(secret, memberID, username, TokenTypeRefresh, ttl)
}

func generateMemberToken(secret string, memberID string, username string, tokenType string, ttl time.Duration) (string, int64, error) {
	now := time.Now().Unix()
	exp := time.Now().Add(ttl).Unix()

	headerJSON, err := json.Marshal(jwtHeader{Alg: "HS256", Typ: "JWT"})
	if err != nil {
		return "", 0, err
	}
	payloadJSON, err := json.Marshal(MemberTokenClaims{MemberID: memberID, Username: username, TokenType: tokenType, Exp: exp, Iat: now})
	if err != nil {
		return "", 0, err
	}

	headerPart := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadPart := base64.RawURLEncoding.EncodeToString(payloadJSON)
	signingInput := fmt.Sprintf("%s.%s", headerPart, payloadPart)

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(signingInput))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("%s.%s", signingInput, signature), exp, nil
}

func ParseMemberToken(secret string, token string) (*MemberTokenClaims, error) {
	claims, err := parseMemberToken(secret, token)
	if err != nil {
		return nil, err
	}
	if claims.TokenType == "" {
		claims.TokenType = TokenTypeAccess
	}
	if claims.TokenType != TokenTypeAccess {
		return nil, ErrInvalidTokenType
	}
	return claims, nil
}

func ParseMemberRefreshToken(secret string, token string) (*MemberTokenClaims, error) {
	claims, err := parseMemberToken(secret, token)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != TokenTypeRefresh {
		return nil, ErrInvalidTokenType
	}
	return claims, nil
}

func parseMemberToken(secret string, token string) (*MemberTokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidTokenFormat
	}

	signingInput := fmt.Sprintf("%s.%s", parts[0], parts[1])
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(signingInput))
	expected := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		return nil, ErrInvalidTokenSign
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidTokenClaims
	}

	claims := &MemberTokenClaims{}
	if err := json.Unmarshal(payload, claims); err != nil {
		return nil, ErrInvalidTokenClaims
	}
	if strings.TrimSpace(claims.MemberID) == "" || claims.Exp == 0 {
		return nil, ErrInvalidTokenClaims
	}
	if time.Now().Unix() >= claims.Exp {
		return nil, ErrTokenExpired
	}

	return claims, nil
}
