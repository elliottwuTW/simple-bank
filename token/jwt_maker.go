package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < 32 {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil
}

type JWTClaims struct {
	// 要用 embed 的方式，不能用 composition => 轉換完的 payload 會不見
	*Payload
	// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-NewWithClaims-CustomClaimsType
	jwt.RegisteredClaims
}

func NewJWTClaims(username string, duration time.Duration) (*JWTClaims, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return nil, err
	}

	return &JWTClaims{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
			IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}, nil
}

func (maker *JWTMaker) SignToken(username string, duration time.Duration) (string, error) {
	claims, err := NewJWTClaims(username, duration)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSigningKey)
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		// 檢查是否為當初使用的簽章演算法
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSigningKey, nil
	})
	if err != nil || !jwtToken.Valid {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("JWT token claims failed")
	}
	return claims.Payload, nil
}
