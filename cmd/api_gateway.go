package cmd

import (
	"Yandex-Taxi-Clone/internal/cache/redis-storage"
	"Yandex-Taxi-Clone/internal/gateway"
	"Yandex-Taxi-Clone/internal/gateway/models"
	"Yandex-Taxi-Clone/internal/transport"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/http"
)

func Run(config models.Config) error {

	ctx := context.Background()

	//Create redis client
	redisClient := newRedisClient(config.Redis.Host, config.Redis.Port)

	//Create cache storage;
	cache := redis_storage.New(redisClient)

	//Create gateway;
	apiGateway := gateway.New(":9001", cache, config.Services, &transport.CustomTransport{
		Context: ctx,
	})

	// creates logic for httputil.ReverseProxy;
	apiGateway.CreateProxy()

	//Register auth service
	if err := apiGateway.RegisterService("auth"); err != nil {
		return err
	}

	if err := apiGateway.RegisterService("notification"); err != nil {
		return err
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		apiGateway.ReverseProxy.ServeHTTP(w, r)
	})

	return http.ListenAndServe(apiGateway.Port, nil)
}

func newRedisClient(host string, port int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}
