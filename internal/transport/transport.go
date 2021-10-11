package transport

import (
	"Yandex-Taxi-Clone/internal/gateway/models"
	"bytes"
	"context"
	"google.golang.org/grpc"
	"io/ioutil"
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
			/*a := new(v1.CreateResponse)
			if err = conn.Invoke(custom.Context, "/v1.UrlShortnerService/Create", &v1.CreateRequest{
				Url: "https://www.google.com/search?client",
			}, a); err != nil {
				log.Println("error:", err)
				return nil, err
			}*/
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
