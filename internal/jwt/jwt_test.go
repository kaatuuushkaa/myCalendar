package jwt_test

import (
	"github.com/stretchr/testify/assert"
	"myCalendar/internal/jwt"
	"testing"
	"time"
)

func TestGenerateAndParseJWT(t *testing.T) {
	svc := jwt.New("testcesret")

	token := svc.GenerateJWT("uuid-stepa", true, jwt.Minute*10)
	assert.NotEmpty(t, token)

	claims, err := svc.ParseJWT(token)
	assert.NoError(t, err)
	assert.Equal(t, "uuid-stepa", claims.ID)
	assert.True(t, claims.IsValid)
	assert.False(t, claims.IsRefresh())
}

func TestGenerateRefreshToken(t *testing.T) {
	svc := jwt.New("testcesret")

	token, exp := svc.GenerateRefreshToken("uuid-stepa", true, jwt.Hour)
	assert.NotEmpty(t, token)
	assert.True(t, exp.After(time.Now()))

	claims, err := svc.ParseJWT(token)
	assert.NoError(t, err)

	assert.True(t, claims.IsRefresh())
}

func TestParseJWT_Expired(t *testing.T) {
	svc := jwt.New("testsecret")

	token := svc.GenerateJWT("uuid-stepa", true, -1)

	_, err := svc.ParseJWT(token)

	assert.ErrorIs(t, err, jwt.ErrTokenExpired)
}

func TestParseJWT_InvalidSecret(t *testing.T) {
	svc1 := jwt.New("secret1")
	svc2 := jwt.New("secret2")

	token := svc1.GenerateJWT("uuid-stepa", true, jwt.Hour)

	_, err := svc2.ParseJWT(token)
	assert.Error(t, err)
	assert.NotErrorIs(t, err, jwt.ErrTokenExpired)
}
