package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Permissions int32 `json:"permissions"`
}

func GenerateJWT(userID uuid.UUID, permissions int32, secret string, expiry time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
		},
		Permissions: permissions,
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string, secret string) (uuid.UUID, int32, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("wrong signing method")
		}
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil || !token.Valid {
		return uuid.UUID{}, 0, fmt.Errorf("error parsing token")
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, 0, fmt.Errorf("error getting subject from token")
	}

	uuidFromSubject, err := uuid.Parse(subject)
	if err != nil {
		return uuid.UUID{}, 0, fmt.Errorf("error parsing into uuid")
	}

	return uuidFromSubject, claims.Permissions, nil
}
