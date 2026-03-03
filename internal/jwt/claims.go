package jwt

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	ID      string `json:"id"`
	IsValid bool   `json:"is_valid"`
	jwt.RegisteredClaims
}

func (c *Claims) IsRefresh() bool {
	return c.Subject == "refresh"
}
