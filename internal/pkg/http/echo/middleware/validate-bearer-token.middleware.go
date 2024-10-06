package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// ValidateBearerToken is an Echo middleware function that validates incoming HTTP requests
// for a Bearer token. It uses the JWT library to parse and validate the token.
//
// The function checks the APP_ENV environment variable to determine if the request is in a test environment.
// If the environment is "test", the middleware function bypasses token validation and calls the next handler.
//
// If the environment is not "test", the middleware function extracts the Bearer token from the request's Authorization
// header or access_token form value. If the token is missing or invalid, it returns an HTTP 401 Unauthorized error.
//
// If the token is valid, the middleware function parses the token using the JWT library and sets the parsed token
// as a value in the Echo context. Finally, it calls the next handler.
func ValidateBearerToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			env := os.Getenv("APP_ENV")
			if env == "test" {
				return next(c)
			}

			auth, ok := bearerAuth(c.Request())
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, errors.New("missing or invalid bearer token"))
			}
			token, err := jwt.ParseWithClaims(
				auth,
				&generates.JWTAccessClaims{},
				func(t *jwt.Token) (interface{}, error) {
					if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, echo.NewHTTPError(http.StatusUnauthorized, errors.New("parse signing method error"))
					}
					return []byte("secret"), nil
				},
			)

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}

			c.Set("token", token)
			return next(c)
		}
	}
}

// bearerAuth extracts and validates a bearer token from an HTTP request.
// It first checks the "Authorization" header for a bearer token with the format "Bearer <token>".
// If no valid token is found in the header, it then checks for a "access_token" form value.
// If a valid token is found, it returns the token along with a boolean value of true.
// If no valid token is found, it returns an empty string and a boolean value of false.
func bearerAuth(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	} else {
		token = r.FormValue("access_token")
	}

	return token, token != ""
}
