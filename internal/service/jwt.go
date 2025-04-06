package service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// генерация JWT-токена
func (s *Service) TokenGenerate(phone string, ttl time.Duration) (string, time.Time, error) {
	expiredAt := time.Now().Add(ttl)

	claims := jwt.MapClaims{
		"phone": phone,
		"exp":   expiredAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenValue, err := token.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		return "", expiredAt, fmt.Errorf("unable to create token for user with phone %s: %w", phone, err)
	}

	return tokenValue, expiredAt, nil
}

// проверка JWT-токена
func (s *Service) TokenVerify(accessToken string) (jwt.MapClaims, error) {
	// parse token
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		// check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(s.config.SecretKey), nil
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
