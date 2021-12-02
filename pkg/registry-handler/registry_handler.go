package registry_handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const registryServiceUrl = "http://localhost:8086/api/service-info/list"

func GetServiceInformationByIdentifier() ([]serviceInformation, error) {
	response, err := http.Get(registryServiceUrl)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var srvsInfo []serviceInformation

	if err := json.Unmarshal(respBytes, &srvsInfo); err != nil {
		return nil, err
	}
	return srvsInfo, nil
}

type serviceInformation struct {
	ID            uint   `json:"id"`
	Identifier    string `json:"identifier"`
	ServiceRoutes []struct {
		ID                   uint   `json:"id"`
		GatewayPath          string `json:"gateway_path"`
		ServicePath          string `json:"service_path"`
		ServiceInformationID uint   `json:"service_information_id"`
	} `json:"service_routes"`
	Services []struct {
		ID                   uint   `json:"id"`
		Host                 string `json:"host"`
		Protocol             string `json:"protocol"`
		Priority             string `json:"priority"`
		ServiceInformationID uint
	} `json:"services"`
}
