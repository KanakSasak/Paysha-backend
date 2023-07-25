package dto

type LoginRequest struct {
	Code          string `json:"code"`
	DeviceType    string `json:"device_type"`
	DeviceID      string `json:"device_id"`
	SigningMethod string `json:"signing_method"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	TokenFcm      string `json:"token_fcm"`
}
