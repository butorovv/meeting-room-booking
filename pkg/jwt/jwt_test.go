package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken("user1", "admin", "secret")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken(t *testing.T) {
	token, _ := GenerateToken("user1", "admin", "secret")
	claims, err := ValidateToken(token, "secret")
	assert.NoError(t, err)
	assert.Equal(t, "user1", claims.UserID)
	assert.Equal(t, "admin", claims.Role)
}

func TestValidateToken_Invalid(t *testing.T) {
	_, err := ValidateToken("invalid", "secret")
	assert.Error(t, err)
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	token, _ := GenerateToken("user1", "admin", "secret1")
	_, err := ValidateToken(token, "secret2")
	assert.Error(t, err)
}

func TestValidateToken_Expired(t *testing.T) {
	claims := Claims{
		UserID: "user1",
		Role:   "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("secret"))

	_, err := ValidateToken(tokenString, "secret")
	assert.Error(t, err)
}
