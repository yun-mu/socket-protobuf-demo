package token

import (
	"config"
	"constant"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func GetJWTToken(data map[string]interface{}) string {
	t := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	claims := t.Claims.(jwt.MapClaims)
	for key, value := range data {
		claims[key] = value
	}
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	jwtToken, _ := t.SignedString([]byte(config.Conf.Security.Secret))

	return jwtToken
}

func VerifyJWTToken(tokenStr string) bool {
	token, err := jwt.Parse(tokenStr, func(jwtToken *jwt.Token) (interface{}, error) {
		return []byte(config.Conf.Security.Secret), nil
	})
	return token.Valid && err == nil
}

func GetJWTClaim(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(jwtToken *jwt.Token) (interface{}, error) {
		return []byte(config.Conf.Security.Secret), nil
	})
	if err != nil || !token.Valid {
		return nil, constant.ErrorUnAuth
	}
	return token.Claims.(jwt.MapClaims), nil
}
