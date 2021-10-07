package models

import (
	"github.com/golang-jwt/jwt/v4"
	"net/http/httputil"
)

type myClaims struct {
	auth []string
	jwt.RegisteredClaims
}

type ApiGateway struct {
	reverseProxy httputil.ReverseProxy
	Port         string
}
