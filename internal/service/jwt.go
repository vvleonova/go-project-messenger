package service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// генерация JWT-токена
func (s *Service) TokenGenerate(phone string) (string, error) {
	claims := jwt.MapClaims{}
	claims["phone"] = phone
	claims["exp"] = time.Now().Add(time.Minute * 2).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenValue, err := token.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		return "", fmt.Errorf("unable to create token for user with phone %s: %s", phone, err)
	}

	return tokenValue, nil
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
