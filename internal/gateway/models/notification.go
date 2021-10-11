package models

import "encoding/json"

func ConvertToNotificationCreateRequest(reqBytes []byte) (NotificationCreate, error) {
	var createReq NotificationCreate
	if err := json.Unmarshal(reqBytes, &createReq); err != nil {
		return NotificationCreate{}, err
	}
	return createReq, nil
}

func ConvertToNotificationSetStatusRequest(reqBytes []byte) (SetStatus, error) {
	var createReq SetStatus
	if err := json.Unmarshal(reqBytes, &createReq); err != nil {
		return SetStatus{}, err
	}
	return createReq, nil
}

type NotificationCreate struct {
	From Geolocation `json:"from"`
	To   Geolocation `json:"to"`
}

type Geolocation struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type SetStatus struct {
	NotificationID     string `json:"notification_id"`
	NotificationStatus string `json:"status"`
}
