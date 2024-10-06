package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"os"

	passwordu "github.com/NekKkMirror/go-app/internal/pkg/utils/auth/password"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// EncryptPassword hashes a plain text password using bcrypt with a configurable cost.
func EncryptPassword(password string, cost int) (string, error) {
	// Validate password strength
	if err := passwordu.ValidatePasswordStrength(password); err != nil {
		return "", err
	}

	pepper, err := LoadPepper()
	if err != nil {
		log.Println("Error loading pepper:", err)
		return "", err
	}

	// Append pepper for extra security
	password = password + pepper

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return "", err
	}
	return string(hashedPassword), nil
}

// LoadPepper loads the pepper from environment variables.
func LoadPepper() (string, error) {
	pepper := os.Getenv("HASH_PEPPER")
	if pepper == "" {
		return "", errors.New("pepper environment variable not set")
	}
	return pepper, nil
}

// ComparePasswords compares a hashed password with a plain text password.
// Returns true if they match, false otherwise.
func ComparePasswords(hashedPassword, password string) bool {
	pepper, err := LoadPepper()
	if err != nil {
		log.Println("Error loading pepper:", err)
		return false
	}
	password = password + pepper

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println("Password comparison failed:", err)
		return false
	}
	return true
}

// LoadHMACSecretKey loads the secret key from environment variables for HMAC.
func LoadHMACSecretKey() (string, error) {
	secretKey := os.Getenv("HMAC_SECRET_KEY")
	if secretKey == "" {
		return "", errors.New("secret key environment variable not set")
	}
	return secretKey, nil
}

// HashDataSHA256 hashes input data using SHA-256 algorithm.
func HashDataSHA256(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	hashedData := hash.Sum(nil)
	return hex.EncodeToString(hashedData)
}

// CreateHMAC creates an HMAC with SHA-256 using a secret key.
func CreateHMAC(data string, key string) (string, error) {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	mac := h.Sum(nil)
	return hex.EncodeToString(mac), nil
}

// EncodeBase64 encodes input data to Base64 URL-safe encoding.
func EncodeBase64(data string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(data))
	return encoded
}

// DecodeBase64 decodes Base64 URL-safe encoded input data.
func DecodeBase64(encodedData string) (string, error) {
	decodedData, err := base64.URLEncoding.DecodeString(encodedData)
	if err != nil {
		return "", err
	}
	return string(decodedData), nil
}

// ValidateHMAC validates an HMAC for given data and key.
func ValidateHMAC(data, receivedMAC, key string) bool {
	expectedMAC, err := CreateHMAC(data, key)
	if err != nil {
		log.Println("Error creating HMAC:", err)
		return false
	}
	return hmac.Equal([]byte(expectedMAC), []byte(receivedMAC))
}
