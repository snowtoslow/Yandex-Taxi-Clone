package transport

import (
	"context"
	"net/http"
)

type CustomTransport struct {
	Host    string
	Context context.Context
}

func New(host string, ctx context.Context) CustomTransport {
	return CustomTransport{
		Host:    host,
		Context: ctx,
	}
}

func (custom CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	/*conn, err := grpc.Dial(custom.Host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	a := new(v1.CreateResponse)
	if err = conn.Invoke(custom.Context, "/v1.UrlShortnerService/Create", &v1.CreateRequest{
		Url: "https://www.google.com/search?client",
	}, a); err != nil {
		log.Println("error:", err)
		return nil, err
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
	}, nil*/
	return nil, nil
}
