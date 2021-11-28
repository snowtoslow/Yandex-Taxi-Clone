package cmd

import (
	"Yandex-Taxi-Clone/internal/gateway"
	"Yandex-Taxi-Clone/internal/gateway/models"
	"Yandex-Taxi-Clone/internal/transport"
	"net/http"
)

func Run(config models.Config) error {
	//Create gateway;
	apiGateway := gateway.New(":9001", config.Services, &transport.CustomTransport{})

	// creates logic for httputil.ReverseProxy;
	apiGateway.CreateProxy()

	//Register auth service
	if err := apiGateway.RegisterServices(); err != nil {
		return err
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		apiGateway.ReverseProxy.ServeHTTP(w, r)
	})

	return http.ListenAndServe(apiGateway.Port, nil)
}
