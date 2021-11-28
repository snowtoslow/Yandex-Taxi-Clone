package transport

import (
	"Yandex-Taxi-Clone/internal/gateway/models"
	"Yandex-Taxi-Clone/utils"
	"bytes"
	"context"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
	"strings"
)

type CustomTransport struct {
	ServiceInformation *models.ServiceInformation
}

func (custom *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-Origin-Host", req.URL.Host)
	ctx := context.Background()
	backEnd := custom.ServiceInformation.GetNextPeer()

	for _, v := range custom.ServiceInformation.Routes {
		if strings.Contains(v.GatewayPath, req.URL.Path) {
			reqBytes, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return createResponse(err, nil, req), nil
			}
			req.Header.Set("Content-Type", "application/grpc+rawCodec")
			if !backEnd.Limiter.Limit() {
				backEnd = custom.ServiceInformation.GetNextPeer()
			}
			clientStream, err := grpc.NewClientStream(ctx, &grpc.StreamDesc{
				ServerStreams: true,
				ClientStreams: true,
			}, backEnd.GrpcClientConn, v.ServicePath)
			if err != nil {
				return createResponse(err, nil, req), nil
			}
			response, err := utils.CreateBytesResponse(clientStream, reqBytes)
			if err != nil {
				//todo: may cause errors fix by adding response create
				return nil, err
			}
			return createResponse(nil, response, req), nil
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

func (custom *CustomTransport) SetServiceInformation(srvInfo *models.ServiceInformation) {
	custom.ServiceInformation = srvInfo
}
