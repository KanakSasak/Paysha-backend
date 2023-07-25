package domain

import (
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/kaleido-io/kaleido-fabric-go/dto"
	"time"
)

type Auth struct {
	AuthID        uint64 `json:"auth_id" gorm:"primaryKey;autoIncrement:true"`
	UID           string `json:"uid" gorm:"primaryKey"`
	SigningMethod string `json:"signing_method" gorm:"not null"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Password      string `json:"password"`
	Role          string `json:"role"`
	TokenFcm      string `json:"token_fcm"`
	UserType      string `json:"user_type"`
	DeviceType    string `json:"device_type" gorm:"not null"`
	DeviceID      string `json:"device_id" gorm:"not null"`
	LastLoginDate string `json:"last_login_date"`
	LastLoginTime string `json:"last_login_time"`
}

type Token struct {
	Expires time.Time   `json:"expires"`
	Token   string      `json:"token"`
	Data    interface{} `json:"data"`
}

func (c Token) ToDto() dto.LoginResponse {
	return dto.LoginResponse{
		Expires: c.Expires,
		Token:   c.Token,
		Data:    c.Data,
	}
}

type AuthRepository interface {
	ValidateToken(c *gin.Context, code string) (*auth.Token, error)
	Insert(c *gin.Context, data *Auth) error
	GetCustomerbyID(c *gin.Context, request dto.LoginRequest, id string) (data *Customer, err error)
	NewAuth(c *gin.Context, uid string, role string, usertype string, request dto.LoginRequest) (*Auth, error)
	GenerateToken(c *gin.Context, data *Auth, dataCustomer *Customer) (*Token, error)
	CreateUserFirebase(c *gin.Context, request dto.RegisterRequest) (string, error)
	RefreshToken(c *gin.Context) (*Token, error)
	CreateCustomer(c *gin.Context, request dto.RegisterRequest, uid string, wallet string) (*Customer, error)
	Register(c *gin.Context, request dto.RegisterRequest) (*Customer, error)
}

type AuthService interface {
	LoginCust(c *gin.Context, request dto.LoginRequest) (*dto.LoginResponse, error)
	WebLoginCust(c *gin.Context, request dto.LoginRequest) (*dto.LoginResponse, error)
	Register(c *gin.Context, request dto.RegisterRequest) (*dto.RegisterResponse, error)
	RefreshToken(c *gin.Context) (*dto.LoginResponse, error)
}
