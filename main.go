package main

import (
	"Yandex-Taxi-Clone/cmd"
	"Yandex-Taxi-Clone/internal/gateway"
	"Yandex-Taxi-Clone/internal/gateway/models"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/encoding"
	"log"
	"os"
)

var config models.Config

func init() {
	configFile, err := os.Open("config-dev.json")
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	encoding.RegisterCodec(gateway.RawCodec{})
}

func main() {
	if err := cmd.Run(config); err != nil {
		log.Fatalln("Error starting apy gateway: ", err)
	}
}
