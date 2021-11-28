package gateway

import (
	"Yandex-Taxi-Clone/internal/gateway/models"
	"Yandex-Taxi-Clone/internal/transport"
	"Yandex-Taxi-Clone/utils"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

type ApiGateway struct {
	ReverseProxy          *httputil.ReverseProxy
	Port                  string
	ServerPool            map[string]*models.ServiceInformation
	ServiceInfoFromConfig models.Services
	Transport             *transport.CustomTransport
}

func New(
	port string,
	serviceFromConfig models.Services,
	customTransport *transport.CustomTransport,
) *ApiGateway {
	return &ApiGateway{
		Port:                  port,
		ServerPool:            map[string]*models.ServiceInformation{},
		ServiceInfoFromConfig: serviceFromConfig,
		Transport:             customTransport,
	}
}

func (apiGateway *ApiGateway) RegisterServices() error {
	for _, fromConfig := range apiGateway.ServiceInfoFromConfig {
		checkedHostsWithStatuses := healthCheck(fromConfig.HostWithStatus, fromConfig.ServiceIdentifier)
		backEnds, err := models.HostsWithStatusesToBackEnds(checkedHostsWithStatuses)
		if err != nil {
			return err
		}
		apiGateway.ServerPool[fromConfig.ServiceIdentifier] = &models.ServiceInformation{
			BackEnds: backEnds,
			Routes:   fromConfig.Routes,
		}
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
			srvInfo, ok := apiGateway.ServerPool[identifier]
			if !ok {
				log.Fatalf("Can't find provided service by identifier: %s", identifier)
			}
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.URL.Scheme = "http"
			req.URL.Host = req.Host
			req.Header.Add("Access-Control-Allow-Origin", "*")
			if !strings.Contains(identifier, "/auth") {
				req.Proto = "HTTP/2.0"
				apiGateway.Transport.SetServiceInformation(srvInfo)
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

func healthCheck(statuses models.HostsWithStatuses, serviceIdentifier string) models.HostsWithStatuses {
	for i := range statuses {
		if strings.Contains(serviceIdentifier, "auth") {
			ok, err := checkHttp(statuses[i].Host)
			if err != nil {
				log.Println("ERROR during check http: ", err)
			}
			statuses[i].Healthy = ok
		} else {
			ok, err := checkRPC(statuses[i].Host, serviceIdentifier)
			if err != nil {
				log.Println("ERROR during connection to rpc service: ", err)
			}
			statuses[i].Healthy = ok
		}
	}

	return statuses
}

func checkRPC(host, identifier string) (bool, error) {
	conn, err := grpc.Dial(
		host,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype("myCodec")))
	defer conn.Close()
	if err != nil {
		return false, err
	}

	stream, err := grpc.NewClientStream(context.Background(), &grpc.StreamDesc{
		ClientStreams: true,
		ServerStreams: true,
	}, conn, fmt.Sprintf("/v1.%sService/Health", strings.Title(identifier)))

	//can skip here cause we are marshalling an empty struct
	reqBytes, _ := json.Marshal(&models.HealthRequest{})

	res, err := utils.CreateBytesResponse(stream, reqBytes)
	if err != nil {
		//todo: add a custom error which will return that healtz endpoint doesn't exist
		return false, err
	}

	var healthResponse models.HealthResponse
	if err := json.Unmarshal(res, &healthResponse); err != nil {
		return false, err
	}

	if healthResponse.Status != "SERVING" {
		return false, nil
	}

	return true, nil
}

func checkHttp(host string) (bool, error) {
	hostWithHttp := fmt.Sprintf("http://%s/health", host)
	if _, err := http.Get(hostWithHttp); err != nil {
		return false, err
	}
	return true, nil
}
