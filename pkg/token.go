package pkg

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = "secret-key"
func GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Генерация подписи токена
	secret := jwtSecret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}