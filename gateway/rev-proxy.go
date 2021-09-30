package main

import (
	v1 "Yandex-Taxi-Clone/pkg/api/v1"
	"context"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func main() {
	log.SetPrefix("[proxy] ")
	log.SetOutput(os.Stdout)
	url, _ := url.Parse("http://localhost:8086")
	path := "/*catchall"
	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", url.Host)
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
		req.Header.Set("Content-Type", "application/grpc+protobuf")

		wildcardIndex := strings.IndexAny(path, "*")
		proxyPath := singleJoiningSlash(url.Path, req.URL.Path[wildcardIndex:])
		if strings.HasSuffix(proxyPath, "/") && len(proxyPath) > 1 {
			proxyPath = proxyPath[:len(proxyPath)-1]
		}
		req.URL.Path = proxyPath
		req.Proto = "HTTP/2.0"
	}
	proxy := &Upstream{target: url, proxy: &httputil.ReverseProxy{
		Director: director,
		Transport: &CustomProtocol{
			Host:    "localhost:8086",
			Context: context.Background(),
		},
		ModifyResponse: func(response *http.Response) error {
			log.Printf("%+v", response)
			return nil
		},
	}}

	mux := http.NewServeMux()
	mux.HandleFunc("/", proxy.handle)
	log.Fatal(http.ListenAndServe(":9001", mux))
	/*http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.proxy.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":9001", nil))*/

}

// Upstream ...
type Upstream struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func (p *Upstream) handle(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

type CustomProtocol struct {
	Host    string
	Context context.Context
}

func (custom *CustomProtocol) RoundTrip(req *http.Request) (*http.Response, error) {
	conn, err := grpc.Dial(custom.Host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	//defer conn.Close()

	a := new(v1.CreateResponse)
	if err = conn.Invoke(custom.Context, "/v1.UrlShortnerService/Create", &v1.CreateRequest{
		Url: "https://www.google.com/search?client",
	}, a); err != nil {
		log.Println("error:", err)
		return nil, err
	}
	log.Println("AAAA: ", a)

	return &http.Response{}, nil
}
