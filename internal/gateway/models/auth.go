package models

type Register struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	IsDriver bool   `json:"isDriver"`
}

type Authenticate struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
