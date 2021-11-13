package gateway

import (
	"Yandex-Taxi-Clone/internal/cache"
	"Yandex-Taxi-Clone/internal/gateway/models"
	"Yandex-Taxi-Clone/internal/transport"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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
	Url     *url.URL
	Routes  []models.Route
	grpcCon *grpc.ClientConn
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
	host, routes, err := apiGateway.ServiceInfoFromConfig.GetInfoFromServiceConfig(identifier)
	if err != nil {
		return err
	}

	apiGateway.Services[identifier] = serviceInformation{
		Url: &url.URL{
			Host: host,
		},
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
			srvInfo, ok := apiGateway.Services[identifier]
			if !ok {
				log.Fatalf("Can't find provided service by identifier: %s", identifier)
			}
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", srvInfo.Url.Host)
			req.URL.Scheme = "http"
			req.URL.Host = req.Host
			req.Header.Add("Access-Control-Allow-Origin", "*")
			if !strings.Contains(identifier, "/auth") {
				req.Proto = "HTTP/2.0"
				apiGateway.Transport.SetGrpcConnection(srvInfo.grpcCon)
				apiGateway.Transport.SetRoutes(srvInfo.Routes)
				apiGateway.SetTransport(apiGateway.Transport)
			} else {
				apiGateway.SetTransport(http.DefaultTransport)
			}
		},
		ErrorHandler: func(rw http.ResponseWriter, r *http.Request, err error) {
			fmt.Printf("error was: %+v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
		},
	}
}
