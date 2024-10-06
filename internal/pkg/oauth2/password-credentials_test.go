package oauth2

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	postgrescontainer "github.com/NekKkMirror/go-app/internal/pkg/container/test/postgres"
	"github.com/NekKkMirror/go-app/internal/pkg/utils/hash"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// initialize environment variables and valid scopes.
func init() {
	err := os.Setenv("HASH_PEPPER", "test_pepper")
	if err != nil {
		return
	}
	err = os.Setenv("VALID_SCOPES", "read:token")
	if err != nil {
		return
	}
}

// Setup initializes the Echo server and the PostgreSQL test container for testing.
func Setup(t *testing.T) (*echo.Echo, *gorm.DB) {
	ctx := context.Background()

	// Start PostgreSQL test container.
	db, _, err := postgrescontainer.Start(ctx, t)
	if err != nil {
		t.Fatalf("failed to start PostgreSQL container: %v", err)
	}

	e := echo.New()
	RunOauthServer(e, db)
	return e, db
}

// TestTokenEndpoint tests the /connect/token endpoint.
func TestTokenEndpoint(t *testing.T) {
	e, db := Setup(t)

	// Create test user.
	_ = createTestUser(t, db)

	req := httptest.NewRequest(http.MethodPost,
		"/connect/token",
		strings.NewReader(
			"grant_type=password&client_id=clientId&client_secret=clientSecret&username=testuser&password=Val]idP3[2E@eas$@rd&scope=all"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	if assert.NoError(t, Token(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "access_token")
	} else {
		t.Fatalf("Failed to handle token: %v", rec.Body.String())
	}
}

// createTestUser creates a user in the database for testing purposes.
func createTestUser(t *testing.T, db *gorm.DB) *User {
	encryptedPassword, err := hash.EncryptPassword("Val]idP3[2E@eas$@rd", bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to encrypt password: %v", err)
	}
	testUser := User{
		ID:        uuid.NewV4(),
		UserName:  "testuser",
		Password:  encryptedPassword,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		DeletedAt: nil,
	}
	if err := db.Create(&testUser).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	return &testUser
}
