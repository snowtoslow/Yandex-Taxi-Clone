package main

import (
	"crypto/tls"
	"golang.org/x/net/http2"
	"io/ioutil"
	"log"
	"net"
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
	}
	proxy := &Upstream{target: url, proxy: &httputil.ReverseProxy{
		Director: director,
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				ta, err := net.ResolveTCPAddr(network, addr)
				if err != nil {
					return nil, err
				}

				return net.DialTCP(network, nil, ta)
			},
		},
		ModifyResponse: func(response *http.Response) error {
			defer response.Body.Close()
			bytes, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return err
			}
			log.Printf("%s", bytes)
			return nil
		},
	}}

	mux := http.NewServeMux()
	mux.HandleFunc("/", proxy.handle)
	log.Fatal(http.ListenAndServe(":9001", mux))

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

/*func main() {
	origin, _ := url.Parse("http://localhost:8086/")
	path := "/*catchall"

	//p := httputil.NewSingleHostReverseProxy(origin)
	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = origin.Scheme
		req.URL.Host = origin.Host

		wildcardIndex := strings.IndexAny(path, "*")
		proxyPath := singleJoiningSlash(origin.Path, req.URL.Path[wildcardIndex:])
		if strings.HasSuffix(proxyPath, "/") && len(proxyPath) > 1 {
			proxyPath = proxyPath[:len(proxyPath)-1]
		}
		req.URL.Path = proxyPath

	}

	modifyResponse := func(response *http.Response) error {
		return nil
	}

	proxy := &httputil.ReverseProxy{
		Director: director,
		ModifyResponse: modifyResponse,
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				ta, err := net.ResolveTCPAddr(network, addr)
				if err != nil {
					return nil, err
				}
				return net.DialTCP(network, nil, ta)
			},
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":9001", nil))
}*/
