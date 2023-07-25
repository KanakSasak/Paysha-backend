package dto

import "time"

type LoginResponse struct {
	Expires time.Time   `json:"expires"`
	Token   string      `json:"token"`
	Data    interface{} `json:"data"`
}
