package authorization

import "github.com/golang-jwt/jwt/v4"

type claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
