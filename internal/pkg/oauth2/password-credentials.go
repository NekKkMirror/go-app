package oauth2

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/NekKkMirror/go-app/internal/pkg/utils/hash"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	srv        *server.Server
	once       sync.Once
	manager    *manage.Manager
	privateKey = []byte(`secret`)
	clients    = []*models.Client{
		{ID: "clientId", Secret: "clientSecret"},
		{ID: "clientId2", Secret: "clientSecret2"},
	}
)

// User model
type User struct {
	ID        uuid.UUID       `json:"userId"   gorm:"primaryKey"`
	FirstName string          `json:"firstName"`
	LastName  string          `json:"lastName"`
	UserName  string          `json:"userName"`
	Email     string          `json:"email"     gorm:"index"`
	Password  string          `json:"password"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	DeletedAt *gorm.DeletedAt `json:"deletedAt"  gorm:"index"`
}

// init initializes the OAuth2 server and its dependencies.
func init() {
	manager = manage.NewDefaultManager()
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", privateKey, jwt.SigningMethodHS512))
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	once.Do(func() {
		srv = server.NewDefaultServer(manager)
	})
}

// RunOauthServer initializes and runs the OAuth2 server with Echo and GORM.
func RunOauthServer(e *echo.Echo, gorm *gorm.DB) {
	manager.MapClientStorage(clientStore(clients...))

	srv.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (userID string, err error) {
		u := User{}
		gorm.Where("user_name = ?", username).First(&u)

		isMatch := hash.ComparePasswords(u.Password, password)
		if isMatch {
			return u.ID.String(), nil
		}
		return
	})

	srv.SetClientScopeHandler(func(tgr *oauth2.TokenGenerateRequest) (allowed bool, err error) {
		if tgr.Scope == "all" {
			allowed = true
		}
		return
	})

	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	srv.SetInternalErrorHandler(func(err error) (res *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	e.GET("connect/token", Token)
	e.GET("validate-token", validateBearerToken)
}

// clientStore creates and initializes a new client store for OAuth2 clients.
func clientStore(clients ...*models.Client) oauth2.ClientStore {
	clientStore := store.NewClientStore()

	for _, c := range clients {
		if c != nil {
			err := clientStore.Set(c.ID, &models.Client{
				ID:     c.ID,
				Secret: c.Secret,
				Domain: c.Domain,
				Public: c.Public,
			})
			if err != nil {
				return nil
			}
		}
	}
	return clientStore
}

// Token is an HTTP handler function that handles the OAuth2 token endpoint.
func Token(c echo.Context) error {
	err := srv.HandleTokenRequest(c.Response().Writer, c.Request())
	if err != nil {
		return err
	}
	return nil
}

// validateBearerToken is an HTTP handler function that validates and returns a bearer token.
func validateBearerToken(c echo.Context) error {
	token, err := srv.ValidationBearerToken(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	return c.JSON(http.StatusOK, token)
}
