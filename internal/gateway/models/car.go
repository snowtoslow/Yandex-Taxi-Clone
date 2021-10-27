package models

import (
	v1 "Yandex-Taxi-Clone/pkg/api/v1"
	"encoding/json"
	"fmt"
)

type Car struct {
	Number   string  `json:"number"`
	Model    string  `json:"model"`
	Color    string  `json:"color"`
	Status   string  `json:"status"`
	Type     string  `json:"type"`
	Location float32 `json:"location"`
}

type FindCarRequest struct {
	Status   string  `json:"status,omitempty"`
	Type     string  `json:"type"`
	Location float32 `json:"location"`
}

func ToFindCarRequest(reqBytes []byte) (*FindCarRequest, error) {
	var findCarRequest FindCarRequest
	if err := json.Unmarshal(reqBytes, &findCarRequest); err != nil {
		return nil, err
	}
	return &findCarRequest, nil
}

func FindCarRequestToProtoObject(findCarRequest *FindCarRequest) (*v1.FindCarRequest, *v1.FindCarResponse, error) {
	protoCarType, err := toProtoCarType(findCarRequest.Type)
	if err != nil {
		return nil, nil, err
	}

	return &v1.FindCarRequest{
		Status:   v1.Status_Free,
		CarType:  protoCarType,
		Location: findCarRequest.Location,
	}, &v1.FindCarResponse{}, nil
}

func ProtoRespToCarModelBytes(resp *v1.FindCarResponse) ([]byte, error) {
	carModel := Car{
		Number:   resp.FoundCar.CarNumber,
		Model:    resp.FoundCar.Model,
		Color:    "color will be added",
		Status:   resp.FoundCar.Status.String(),
		Type:     resp.FoundCar.CarType.String(),
		Location: resp.FoundCar.Location,
	}

	carBytes, err := json.Marshal(carModel)
	if err != nil {
		return nil, err
	}
	return carBytes, nil
}

func toProtoCarType(typeString string) (v1.Type, error) {
	switch typeString {
	case "standard":
		return v1.Type_Standard, nil
	case "comfort":
		return v1.Type_Comfort, nil
	default:
		return v1.Type_UnknownType, fmt.Errorf("invalid car type provided: %s", typeString)
	}
}
