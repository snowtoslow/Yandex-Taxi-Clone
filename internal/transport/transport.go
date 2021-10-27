package transport

import (
	"Yandex-Taxi-Clone/internal/cache"
	"Yandex-Taxi-Clone/internal/gateway/models"
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
	Cache   cache.Repository
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
			/*var protoReq interface{}
			var protoResp interface{}*/
			reqBytes, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			splits := strings.Split(req.URL.Path, "/")
			switch identifier := splits[2]; identifier {
			case "car":
				switch method := splits[3]; method {
				case "find":
					findCarReq, err := models.ToFindCarRequest(reqBytes)
					if err != nil {
						return nil, err
					}

					cachedData, err := custom.Cache.GetCachedData(custom.Context, findCarReq.Status)
					if err != nil && err.Error() != "redis: nil" {
						return nil, err
					}

					if len(cachedData) == 0 {
						protoReq, protoResp, err := models.FindCarRequestToProtoObject(findCarReq)
						if err != nil {
							return nil, err
						}
						if err = conn.Invoke(custom.Context, v.ServicePath, protoReq, protoResp); err != nil {
							log.Println("error invoking:", err)
							return nil, err
						}

						carBytes, err := models.ProtoRespToCarModelBytes(protoResp)
						if err != nil {
							return nil, err
						}

						if err = custom.Cache.
							SetCachedData(custom.Context, protoReq.Status.String(), carBytes); err != nil {
							return nil, err
						}

						return &http.Response{
							Status:     http.StatusText(http.StatusOK),
							StatusCode: http.StatusOK,
							Header: map[string][]string{
								"Content-Type": {"application/json"},
							},
							Body:    ioutil.NopCloser(ioutil.NopCloser(bytes.NewBufferString("To set cached data"))),
							Request: req,
						}, nil
					} else {
						log.Println("NOT CACHED DATA")
					}
					return &http.Response{
						Status:     http.StatusText(http.StatusOK),
						StatusCode: http.StatusOK,
						Header: map[string][]string{
							"Content-Type": {"application/json"},
						},
						Body:    ioutil.NopCloser(ioutil.NopCloser(bytes.NewBufferString(cachedData))),
						Request: req,
					}, nil

				}
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
