package gateway

import (
	"Yandex-Taxi-Clone/internal/cache"
	"Yandex-Taxi-Clone/internal/gateway/models"
	"Yandex-Taxi-Clone/internal/transport"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

type ApiGateway struct {
	ReverseProxy          *httputil.ReverseProxy
	Port                  string
	Cache                 cache.Repository
	Services              map[string]serviceInformation
	ServiceInfoFromConfig models.Services
	Transport             *transport.CustomTransport
}

type serviceInformation struct {
	Url    string
	Routes []models.Route
}

func New(
	port string,
	serviceFromConfig models.Services,
	customTransport *transport.CustomTransport,
) *ApiGateway {
	return &ApiGateway{
		Port:                  port,
		Services:              map[string]serviceInformation{},
		ServiceInfoFromConfig: serviceFromConfig,
		Transport:             customTransport,
	}
}

func (apiGateway *ApiGateway) RegisterService(identifier string) error {
	host, port, routes, err := apiGateway.ServiceInfoFromConfig.GetInfoFromServiceConfig(identifier)
	if err != nil {
		return err
	}

	apiGateway.Services[identifier] = serviceInformation{
		Url:    fmt.Sprintf("%s:%d", host, port),
		Routes: routes,
	}
	return nil
}

func (apiGateway *ApiGateway) SetTransport(transport http.RoundTripper) {
	apiGateway.ReverseProxy.Transport = transport
}

func (apiGateway *ApiGateway) CreateProxy() {
	apiGateway.ReverseProxy = &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			identifier := strings.Split(req.URL.Path, "/")[2]
			log.Println("IDENTIFIER: ", identifier)
			srvInfo, ok := apiGateway.Services[identifier]
			if !ok {
				log.Fatalf("Can't find provided service by identifier: %s", identifier)
			}
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", srvInfo.Url)
			req.URL.Scheme = "http"
			req.URL.Host = req.Host
			req.Header.Add("Access-Control-Allow-Origin", "*")
			if !strings.Contains(identifier, "/auth") {
				req.Proto = "HTTP/2.0"
				apiGateway.Transport.SetHost(srvInfo.Url)
				apiGateway.Transport.SetRoutes(srvInfo.Routes)
				apiGateway.SetTransport(apiGateway.Transport)
			}

		},
		ErrorHandler: func(rw http.ResponseWriter, r *http.Request, err error) {
			fmt.Printf("error was: %+v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
		},
	}
}
