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

// ValidateBearerToken validates incoming HTTP requests for a Bearer token.
func ValidateBearerToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if isTestEnvironment() {
				return next(c)
			}

			authToken, ok := extractBearerToken(c.Request())
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid bearer token")
			}

			token, err := parseJWT(authToken)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			c.Set("token", token)
			return next(c)
		}
	}
}

// isTestEnvironment checks if the application is running in a test environment.
func isTestEnvironment() bool {
	return os.Getenv("APP_ENV") == "tests"
}

// extractBearerToken retrieves the Bearer token from the request header or form data.
func extractBearerToken(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")
	if token := extractTokenFromHeader(auth); token != "" {
		return token, true
	}
	return r.FormValue("access_token"), r.FormValue("access_token") != ""
}

// extractTokenFromHeader extracts the token from the "Authorization" header.
func extractTokenFromHeader(auth string) string {
	const prefix = "Bearer "
	if auth == "" || !strings.HasPrefix(auth, prefix) {
		return ""
	}
	return auth[len(prefix):]
}

// parseJWT parses and validates the JWT token.
func parseJWT(authToken string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(
		authToken,
		&generates.JWTAccessClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte("secret"), nil
		},
	)
}
