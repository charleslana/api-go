package services

import (
	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

const Secret = "secret"

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(Secret), nil)
}

func GetTokenAuth() *jwtauth.JWTAuth {
	return tokenAuth
}
