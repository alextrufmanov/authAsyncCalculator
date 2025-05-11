package orchestrator

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var privateKey = []byte("7444071A-63F5-4829-9344-D4070E5C8DF7")

func GetJWT(login string, userId int) (string, bool) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login":   login,
		"user_id": userId,
		"nbf":     time.Now().Unix(),
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	})
	jwt, err := token.SignedString(privateKey)
	if err != nil {
		return "", false
	}
	return jwt, true
}

func JWTToUserId(tokenStr string) (int32, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return privateKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return int32(claims["user_id"].(float64)), nil
	}
	return 0, err
}
