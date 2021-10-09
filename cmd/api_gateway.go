package cmd

import (
	"Yandex-Taxi-Clone/internal/cache/redis-storage"
	"Yandex-Taxi-Clone/internal/gateway"
	"Yandex-Taxi-Clone/internal/gateway/models"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

func Run(config models.Config) error {
	log.Printf("%+v", config)
	proxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			switch a := strings.Split(req.URL.Path, "/")[2]; a {
			case "auth":
				for _, v := range config.Services {
					if v.ServiceIdentifier == a {
						host := fmt.Sprintf("%s:%d", v.Host, v.Port)
						req.Header.Add("X-Forwarded-Host", req.Host)
						req.Header.Add("X-Origin-Host", host)
						req.URL.Scheme = "http"
						req.URL.Host = host
						//ADD CORS - to allow access of magic-costet-front-end:
						req.Header.Add("Access-Control-Allow-Origin", "*")
						break
					}
				}

			}

		},
		ModifyResponse: func(response *http.Response) error {
			return nil
			//
			// purposefully return an error so ErrorHandler gets called
			//return errors.New("uh-oh")
		},
		ErrorHandler: func(rw http.ResponseWriter, r *http.Request, err error) {
			fmt.Printf("error was: %+v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
		},
	}

	redisClient := newRedisClient(config.Redis.Host, config.Redis.Port)

	cache := redis_storage.New(redisClient)

	apiGateway := gateway.New(&proxy, ":9001", cache, config.Services)

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
