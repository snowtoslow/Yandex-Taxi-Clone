package transport

import (
	"Yandex-Taxi-Clone/internal/cache"
	"Yandex-Taxi-Clone/internal/gateway/models"
	"bytes"
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	"io"
	"io/ioutil"
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
	encoding.RegisterCodec(rawCodec{})
	conn, err := grpc.Dial(
		"localhost:8082",
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype("myCodec")),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	for _, v := range custom.Routes {
		if strings.Contains(v.GatewayPath, req.URL.Path) {
			reqBytes, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return createResponse(err, nil, req), nil
			}

			req.Header.Set("Content-Type", "application/grpc+rawCodec")
			clientStream, err := grpc.NewClientStream(custom.Context, &grpc.StreamDesc{
				ServerStreams: true,
				ClientStreams: true,
			}, conn, v.ServicePath)
			if err != nil {
				return createResponse(err, nil, req), nil
			}

			errCh1 := sendRequestToBackend(clientStream, reqBytes)
			errCh2, retChan := retrieveResponseFromBackEnd(clientStream)
			for i := 0; i < 3; i++ {
				select {
				case magErr := <-errCh1:
					if !errors.Is(magErr, io.EOF) {
						return createResponse(magErr, nil, req), nil
					}
					clientStream.CloseSend()
				case magErr2 := <-errCh2:
					if !errors.Is(magErr2, io.EOF) {
						return createResponse(magErr2, nil, req), nil
					}
				case response := <-retChan:
					return createResponse(nil, response, req), nil
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

func retrieveResponseFromBackEnd(stream grpc.ClientStream) (chan error, chan []byte) {
	errCh := make(chan error, 1)
	retChan := make(chan []byte, 1)
	var a []byte
	go func() {
		for {
			if err := stream.RecvMsg(&a); err != nil {
				if errors.Is(err, io.EOF) {
					retChan <- a
				}
				errCh <- err
				break
			}
		}

	}()

	return errCh, retChan
}

func sendRequestToBackend(stream grpc.ClientStream, info []byte) chan error {
	ret := make(chan error, 1)
	go func() {
		for {
			if err := stream.SendMsg(&info); err != nil {
				ret <- err

				break
			}
		}

	}()
	return ret
}
