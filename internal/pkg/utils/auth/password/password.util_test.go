package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidPasswordStrength(t *testing.T) {
	validPassword := "Valid1@Password"
	err := ValidatePasswordStrength(validPassword)
	assert.Nil(t, err)
}

func TestPasswordTooShort(t *testing.T) {
	shortPassword := "Short1!"
	err := ValidatePasswordStrength(shortPassword)
	assert.Equal(t, ErrPasswordTooShort, err)
}

// Password with mixed case and symbols but meeting requirements returns nil
func TestPasswordStrengthMixedCaseAndSymbols(t *testing.T) {
	mixedCaseAndSymbolsPassword := "MixedCase1@Symbols"
	err := ValidatePasswordStrength(mixedCaseAndSymbolsPassword)
	assert.NoError(t, err)
}

// Password with whitespace characters but meeting requirements returns nil
func TestPasswordWithWhitespaceReturnsNil(t *testing.T) {
	passwordWithWhitespace := "Pass word1@"
	err := ValidatePasswordStrength(passwordWithWhitespace)
	assert.NoError(t, err)
}

// Password with repeated characters but meeting requirements returns nil
func TestPasswordWithRepeatedCharacters(t *testing.T) {
	repeatedPassword := "Repe4ted!Pass"
	err := ValidatePasswordStrength(repeatedPassword)
	assert.NoError(t, err)
}

// Password with non-ASCII characters but meeting requirements returns nil
func TestPasswordWithNonASCIICharacters(t *testing.T) {
	nonASCIIValidPassword := "NonASCII1@￣ﾃﾑ￣ﾂﾹ￣ﾃﾯ￣ﾃﾼ￣ﾃﾉ"
	err := ValidatePasswordStrength(nonASCIIValidPassword)
	assert.NoError(t, err)
}

// Password with excessive length but missing required characters returns ErrPasswordWeak
func TestExcessiveLengthMissingRequiredCharsReturnsErrPasswordWeak(t *testing.T) {
	longPassword := "LongPasswordWithoutRequiredChars123"
	err := ValidatePasswordStrength(longPassword)
	assert.EqualError(t, err, ErrPasswordWeak.Error())
}

// Password missing special character returns ErrPasswordWeak
func TestPasswordMissingSpecialCharacter(t *testing.T) {
	password := "WeakPassword123"
	err := ValidatePasswordStrength(password)
	assert.EqualError(t, err, ErrPasswordWeak.Error())
}

// Password with only numbers returns ErrPasswordWeak
func TestPasswordWithOnlyNumbersReturnsErrPasswordWeak(t *testing.T) {
	password := "12345678"
	err := ValidatePasswordStrength(password)
	assert.EqualError(t, err, ErrPasswordWeak.Error())
}

// Password missing lowercase letter returns ErrPasswordWeak
func TestPasswordMissingLowercase(t *testing.T) {
	password := "UPPERCASE123@"
	err := ValidatePasswordStrength(password)
	assert.EqualError(t, err, ErrPasswordWeak.Error())
}

// Password missing number returns ErrPasswordWeak
func TestPasswordMissingNumberReturnsErrPasswordWeak(t *testing.T) {
	password := "NoNumberPassword@"
	err := ValidatePasswordStrength(password)
	assert.EqualError(t, err, ErrPasswordWeak.Error())
}

// Password missing uppercase letter returns ErrPasswordWeak
func TestPasswordMissingUppercase(t *testing.T) {
	password := "weakpassword1@"
	err := ValidatePasswordStrength(password)
	assert.EqualError(t, err, ErrPasswordWeak.Error())
}

// Password with only special characters returns ErrPasswordWeak
func TestPasswordOnlySpecialCharacters(t *testing.T) {
	specialCharsPassword := "!@#$%^&*()"
	err := ValidatePasswordStrength(specialCharsPassword)
	assert.EqualError(t, err, ErrPasswordWeak.Error())
}

// Password longer than minimum length with all required characters returns nil
func TestPasswordStrengthValidLongPassword(t *testing.T) {
	longValidPassword := "LongValid1@Password"
	err := ValidatePasswordStrength(longValidPassword)
	assert.NoError(t, err)
}

func TestPasswordAtMinimumLengthWithRequiredCharacters(t *testing.T) {
	password := "Valid1@P"
	err := ValidatePasswordStrength(password)
	assert.NoError(t, err)
}
