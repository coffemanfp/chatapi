package auth

import (
	"github.com/golang-jwt/jwt"
)

func GenerateJWT(secretKey string, id int) (tokenS string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = id

	return token.SignedString([]byte(secretKey))
}
