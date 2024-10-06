package hash

import (
	"os"
	"testing"

	"github.com/NekKkMirror/go-app/internal/pkg/utils/auth/password"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	// Setting up environment variables before running tests
	os.Setenv("HASH_PEPPER", "test_pepper")
	os.Setenv("HMAC_SECRET_KEY", "test_secret_key")
}

func TestValidatePasswordStrength(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		password string
		expected error
	}{
		{"", password.ErrPasswordTooShort},
		{"short", password.ErrPasswordTooShort},
		{"noNumber!", password.ErrPasswordWeak},
		{"NoSymbol8", password.ErrPasswordWeak},
		{"ValidPassw0rd!", nil},
	}

	for _, test := range tests {
		err := password.ValidatePasswordStrength(test.password)
		assert.Equal(test.expected, err)
	}
}

func TestHashPassword(t *testing.T) {
	assert := assert.New(t)

	password := "ValidPassw0rd!"
	cost := bcrypt.DefaultCost

	hashedPassword, err := EncryptPassword(password, cost)
	assert.NoError(err)
	assert.NotEmpty(hashedPassword)

	// Trying to hash a weak password should return an error
	weakPassword := "weak"
	_, err = EncryptPassword(weakPassword, cost)
	assert.Error(err)
}

func TestComparePasswords(t *testing.T) {
	assert := assert.New(t)

	password := "ValidPassw0rd!"
	cost := bcrypt.DefaultCost

	hashedPassword, err := EncryptPassword(password, cost)
	assert.NoError(err)

	// Correct password should match
	assert.True(ComparePasswords(hashedPassword, password))

	// Incorrect password should not match
	assert.False(ComparePasswords(hashedPassword, "WrongPassword"))
}

func TestLoadPepper(t *testing.T) {
	assert := assert.New(t)

	pepper, err := LoadPepper()
	assert.NoError(err)
	assert.Equal("test_pepper", pepper)
}

func TestLoadHMACSecretKey(t *testing.T) {
	assert := assert.New(t)

	secretKey, err := LoadHMACSecretKey()
	assert.NoError(err)
	assert.Equal("test_secret_key", secretKey)
}

func TestHashDataSHA256(t *testing.T) {
	assert := assert.New(t)

	data := "test data"
	expectedHash := "916f0027a575074ce72a331777c3478d6513f786a591bd892da1a577bf2335f9"

	actualHash := HashDataSHA256(data)
	assert.Equal(expectedHash, actualHash)
}

func TestCreateHMAC(t *testing.T) {
	assert := assert.New(t)

	data := "test data"
	expectedMac := "9f8a075e0cc4eda7c6aff29369ab5503115a7b34a136a16c3421829790f795ce"

	key, err := LoadHMACSecretKey()
	assert.NoError(err)

	actualMac, err := CreateHMAC(data, key)
	assert.NoError(err)
	assert.Equal(expectedMac, actualMac)
}

func TestEncodeBase64(t *testing.T) {
	assert := assert.New(t)

	data := "test data"
	expectedEncodedData := "dGVzdCBkYXRh"

	actualEncodedData := EncodeBase64(data)
	assert.Equal(expectedEncodedData, actualEncodedData)
}

func TestDecodeBase64(t *testing.T) {
	assert := assert.New(t)

	encodedData := "dGVzdCBkYXRh"
	expectedDecodedData := "test data"

	actualDecodedData, err := DecodeBase64(encodedData)
	assert.NoError(err)
	assert.Equal(expectedDecodedData, actualDecodedData)
}

func TestValidateHMAC(t *testing.T) {
	assert := assert.New(t)

	data := "test data"
	mac := "9f8a075e0cc4eda7c6aff29369ab5503115a7b34a136a16c3421829790f795ce"

	key, err := LoadHMACSecretKey()
	assert.NoError(err)

	isValidMAC := ValidateHMAC(data, mac, key)
	assert.True(isValidMAC)
}
