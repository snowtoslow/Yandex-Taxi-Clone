package models

import (
	"fmt"
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
	Services Services `json:"services"`
}

// Service struct which holds the information about a service registered in gateway;
type Service struct {
	Host              string  `json:"host"`
	Port              int     `json:"port"`
	ServiceIdentifier string  `json:"service_identifier"`
	Routes            []Route `json:"routes"`
}

type Services []Service

func (ss Services) GetInfoFromServiceConfig(identifier string) (string, int, []Route, error) {
	for _, v := range ss {
		if v.ServiceIdentifier == identifier {
			return v.Host, v.Port, v.Routes, nil
		}
	}
	return "", 0, nil, fmt.Errorf("no routes by provided identifier")
}

// Route struct which handles the mapping of gateway url path's to grpc Methods or auth paths;
type Route struct {
	GatewayPath string `json:"gateway_path"`
	ServicePath string `json:"service_path"`
}
