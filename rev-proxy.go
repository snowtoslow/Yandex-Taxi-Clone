package main

/*func main() {
	log.SetPrefix("[proxy] ")
	log.SetOutput(os.Stdout)
	url, _ := url.Parse("http://localhost:8080")
	path := "/*catchall"
	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", url.Host)
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
		//ADD CORS - to allow access of magic-costet-front-end:
		req.Header.Add("Access-Control-Allow-Origin", "*")

		wildcardIndex := strings.IndexAny(path, "*")
		proxyPath := singleJoiningSlash(url.Path, req.URL.Path[wildcardIndex:])
		if strings.HasSuffix(proxyPath, "/") && len(proxyPath) > 1 {
			proxyPath = proxyPath[:len(proxyPath)-1]
		}
		req.URL.Path = proxyPath
	}
	proxy := &Upstream{target: url, proxy: &httputil.ReverseProxy{
		Director: director,
		ModifyResponse: func(response *http.Response) error {
			log.Printf("%+v", response)
			return nil
		},
	}}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/api/auth") {
			if tokenString := r.Header.Get("Authorization"); len(tokenString) != 0 {
				token, err := jwt.Parse(tokenString[7:], func(token *jwt.Token) (interface{}, error) {
					// Don't forget to validate the alg is what you expect:
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
					}

					return []byte{97, 100, 56, 100, 53, 98, 55, 56, 49, 49, 98, 49, 53, 57, 54, 53, 97, 101, 55, 49, 102, 102, 52, 55, 98, 51, 99, 57, 101, 102, 56, 97, 56, 50, 97, 102, 53, 99, 56, 53, 102, 55, 99, 55, 49, 102, 100, 100, 55, 50, 52, 102, 102, 102, 97, 100, 57, 100, 99, 52, 97, 57, 53, 57, 56, 97, 52, 49, 100, 54, 101, 57, 55, 98, 53, 97, 49, 52, 48, 98, 56, 98, 48, 56, 55, 98, 55, 102, 48, 56, 99, 50, 98, 100, 55, 56, 101, 56, 99, 98, 54, 56, 57, 49, 57, 99, 56, 55, 48, 48, 48, 48, 54, 57, 57, 97, 102, 98, 100, 50, 54, 52, 49, 98, 100, 98, 102, 56}, nil
				})

				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					fmt.Println("Claims: ", claims)
				} else {
					fmt.Println(err)
				}

			} else {
				log.Printf("EMPTY TOKEN")
			}
			r.Proto = "HTTP/2.0"
			proxy.proxy.Transport = &CustomProtocol{
				Host:    "localhost:8086",
				Context: context.Background(),
			}
		}

		proxy.proxy.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":9001", nil))

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
	}, nil
}*/
