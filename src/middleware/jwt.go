package middleware

import (
	"config"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// CustomJWTConfig custom jwt config
func CustomJWTConfig(skipperPaths []string, authScheme string) middleware.JWTConfig {
	if authScheme == "" {
		authScheme = "Bearer"
	}
	return middleware.JWTConfig{
		SigningKey:  []byte(config.Conf.Security.Secret),
		TokenLookup: "header:" + echo.HeaderAuthorization,
		AuthScheme:  authScheme,
		Claims:      jwt.MapClaims{},
		ContextKey:  config.Conf.Security.Secret,
		Skipper:     CustomSkipper(skipperPaths),
	}
}

// CustomSkipper custom skipper
func CustomSkipper(skipperPaths []string) func(c echo.Context) bool {
	return func(c echo.Context) bool {
		for _, skipperPath := range skipperPaths {
			if c.Path() == skipperPath {
				return true
			}
		}
		return false
	}
}
