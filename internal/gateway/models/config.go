package models

import (
	"fmt"
	"google.golang.org/grpc"
	"net/url"
	"sync"
)

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
	HostWithStatus    HostsWithStatuses `json:"hosts"`
	ServiceIdentifier string            `json:"service_identifier"`
	Routes            []Route           `json:"routes"`
}

type Services []Service

func (ss Services) GetInfoFromServiceConfig(identifier string) (HostsWithStatuses, []Route, error) {
	for _, v := range ss {
		if v.ServiceIdentifier == identifier {
			return v.HostWithStatus, v.Routes, nil
		}
	}
	return nil, nil, fmt.Errorf("no routes by provided identifier")
}

// Route struct which handles the mapping of gateway url path's to grpc Methods or auth paths;
type Route struct {
	GatewayPath string `json:"gateway_path"`
	ServicePath string `json:"service_path"`
}

func HostsWithStatusesToBackEnds(hostsAndStatuses HostsWithStatuses) ([]*Backend, error) {
	backEnds := make([]*Backend, 0, len(hostsAndStatuses))
	for _, v := range hostsAndStatuses {
		backEnd := hostWithStatusToBackEnd(v)
		if err := backEnd.SetConn(
			grpc.WithDefaultCallOptions(grpc.CallContentSubtype("myCodec")),
			grpc.WithInsecure(),
		); err != nil {
			return nil, err
		}
		backEnds = append(backEnds, backEnd)
	}
	return backEnds, nil
}

type HostsWithStatuses []HostWithStatus

func hostWithStatusToBackEnd(hostAndStatus HostWithStatus) *Backend {
	return &Backend{
		URL: &url.URL{
			Host: hostAndStatus.Host,
		},
		Alive:   hostAndStatus.Healthy,
		Mux:     sync.RWMutex{},
		Limiter: NewLimiter(10),
	}
}

type HostWithStatus struct {
	Host    string `json:"host"`
	Healthy bool   `json:"healthy"`
}
