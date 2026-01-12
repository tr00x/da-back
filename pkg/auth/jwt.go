package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateToken(
	id int, expiration time.Duration, secret_key string, role_id int,
) string {
	unixTime := time.Now().Add(expiration).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      id,
		"role_id": role_id,
		"exp":     unixTime,
	})

	tokenString, _ := token.SignedString([]byte(secret_key))

	return tokenString
}

func CreateRefreshAccsessToken(id, role_id int) (string, string) {

	accessToken := CreateToken(id, ENV.REFRESH_TIME, ENV.ACCESS_KEY, role_id)
	refreshToken := CreateToken(id, ENV.REFRESH_TIME, ENV.REFRESH_KEY, role_id)

	return accessToken, refreshToken
}

func HashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}
