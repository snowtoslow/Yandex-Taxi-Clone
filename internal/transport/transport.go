package transport

import (
	"Yandex-Taxi-Clone/internal/gateway/models"
	v1 "Yandex-Taxi-Clone/pkg/api/v1"
	"bytes"
	"context"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type CustomTransport struct {
	Host    string
	Context context.Context
	Routes  []models.Route
}

func (custom *CustomTransport) SetHost(host string) {
	custom.Host = host
}

func (custom *CustomTransport) SetRoutes(routes []models.Route) {
	custom.Routes = routes
}

func (custom *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	conn, err := grpc.Dial(custom.Host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	for _, v := range custom.Routes {
		if strings.Contains(v.GatewayPath, req.URL.Path) {
			var protoReq interface{}
			var protoResp interface{}
			reqBytes, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			splits := strings.Split(req.URL.Path, "/")
			switch identifier := splits[2]; identifier {
			case "notification":
				switch method := splits[3]; method {
				case "create":
					createReq, err := models.ConvertToNotificationCreateRequest(reqBytes)
					if err != nil {
						return nil, err
					}

					protoReq = &v1.NotificationCreateRequest{
						From: &v1.Coordinates{
							Latitude:  createReq.From.Latitude,
							Longitude: createReq.From.Longitude,
						},
						To: &v1.Coordinates{
							Latitude:  createReq.To.Latitude,
							Longitude: createReq.To.Longitude,
						},
						Status: v1.NotificationStatus_CREATE_ORDER_ATTEMPT,
					}
					protoResp = new(v1.NotificationCreateResponse)
				}
			}

			if err = conn.Invoke(custom.Context, v.ServicePath, protoReq, protoResp); err != nil {
				log.Println("error invoking:", err)
				return nil, err
			}

		}
	}

	//Create custom httpResponses for status ok, unautorized and bad request
	return &http.Response{
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body:    ioutil.NopCloser(ioutil.NopCloser(bytes.NewBufferString("Hello World"))),
		Request: req,
	}, nil
}
