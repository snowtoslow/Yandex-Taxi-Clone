package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type myClaims struct {
	auth []string
	jwt.RegisteredClaims
}

// Config structure handle the whole information about the configs;
type Config struct {
	Redis struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"redis"`
	Services []Service `json:"services"`
}

// Service struct which holds the information about a service registered in gateway;
type Service struct {
	Host              string  `json:"host"`
	Port              int     `json:"port"`
	ServiceIdentifier string  `json:"service_identifier"`
	Routes            []Route `json:"routes"`
}

// Route struct which handles the mapping of gateway url path's to grpc Methods or auth paths;
type Route struct {
	GatewayPath string `json:"gateway_path"`
	ServicePath string `json:"service_path"`
}
