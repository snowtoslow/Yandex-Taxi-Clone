package models

import "github.com/golang-jwt/jwt/v4"

type Register struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	IsDriver bool   `json:"isDriver"`
}

type Authenticate struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type myClaims struct {
	Auth []string `json:"auth,omitempty"`
	jwt.RegisteredClaims
}
