package gateway

import (
	"Yandex-Taxi-Clone/internal/cache"
	"Yandex-Taxi-Clone/internal/gateway/models"
	"net/http"
	"net/http/httputil"
)

type ApiGateway struct {
	ReverseProxy *httputil.ReverseProxy
	Port         string
	Cache        cache.Repository
	Services     []models.Service
}

func New(
	proxy *httputil.ReverseProxy,
	port string,
	repo cache.Repository,
	services []models.Service,
) *ApiGateway {
	return &ApiGateway{
		ReverseProxy: proxy,
		Port:         port,
		Cache:        repo,
		Services:     services,
	}
}

func (apiGateway *ApiGateway) SetTransport(transport http.RoundTripper) {
	apiGateway.ReverseProxy.Transport = transport
}
