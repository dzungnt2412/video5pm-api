package utils

import (
	"strings"
	"time"
	"video5pm-api/core/constants"
	"video5pm-api/models/entity"

	"github.com/dgrijalva/jwt-go"
)

// GenerateToken generate tokens used for auth
func GenerateToken(u *entity.User, groups []int64, references []string, secretKey string) (string, error) {
	now := time.Now()
	exp := now.Add(constants.TOKEN_EXPIRE * time.Minute).Unix()

	claims := jwt.MapClaims{}
	claims["user_id"] = u.ID
	claims["iss"] = constants.JWT_ISSUER
	claims["exp"] = exp
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// ParseToken parsing token
func ParseToken(tokenString, secretKey string) (jwt.MapClaims, error) {
	tokenString = strings.Replace(tokenString, constants.TOKEN_PREFIX, "", -1)
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		// if err != nil {
		// 	return nil, err
		// }
	}

	return claims, nil
}
