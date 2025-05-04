package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// генерация JWT-токена
func TokenGenerate(id uuid.UUID, secretKey string, ttl time.Duration) (string, time.Time, error) {
	expiredAt := time.Now().Add(ttl)

	claims := jwt.MapClaims{
		"id":  id.String(),
		"exp": expiredAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenValue, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", expiredAt, fmt.Errorf("unable to create token for user with id %s: %w", id, err)
	}

	return tokenValue, expiredAt, nil
}

// проверка JWT-токена
func TokenVerify(accessToken, secretKey string) (jwt.MapClaims, error) {
	// parse token
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		// check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// validate token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid JWT-token")
}
