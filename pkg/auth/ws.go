package auth

import (
	"dubai-auto/internal/model"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateWSJWT(tokenString string) (*model.WSUser, error) {

	if tokenString == "" {
		return nil, fmt.Errorf("missing token")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(
		tokenString, claims,
		func(t *jwt.Token) (any, error) {
			return []byte(ENV.ACCESS_KEY), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	userID := int(claims["id"].(float64))
	roleID := int(claims["role_id"].(float64))

	return &model.WSUser{
		ID:     userID,
		RoleID: roleID,
	}, nil
}
