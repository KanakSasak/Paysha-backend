package service

import (
	"github.com/gin-gonic/gin"
	"github.com/kaleido-io/kaleido-fabric-go/domain"
	"github.com/kaleido-io/kaleido-fabric-go/dto"
)

type AuthService struct {
	authrepo domain.AuthRepository
}

func (a AuthService) Register(c *gin.Context, request dto.RegisterRequest) (*dto.RegisterResponse, error) {
	data, err := a.authrepo.Register(c, request)
	if err != nil {
		return nil, err
	}
	resp := dto.RegisterResponse{
		UID:   data.FireUid,
		Phone: data.Phone,
		Name:  data.Name,
		Email: data.Email,
	}

	return &resp, nil
}

func (a AuthService) LoginCust(c *gin.Context, request dto.LoginRequest) (*dto.LoginResponse, error) {
	token, err := a.authrepo.ValidateToken(c, request.Code)
	if err != nil {
		return nil, err
	}
	role := "user"
	usertype := "customer"

	auth, err := a.authrepo.NewAuth(c, token.UID, role, usertype, request)
	if err != nil {
		return nil, err
	}
	err = a.authrepo.Insert(c, auth)
	if err != nil {
		return nil, err
	}
	data, err := a.authrepo.GetCustomerbyID(c, request, auth.UID)
	if err != nil {
		return nil, err
	}
	tokenjwt, err := a.authrepo.GenerateToken(c, auth, data)
	if err != nil {
		return nil, err
	}
	b := tokenjwt.ToDto()

	return &b, nil
}

func (a AuthService) WebLoginCust(c *gin.Context, request dto.LoginRequest) (*dto.LoginResponse, error) {
	token, err := a.authrepo.ValidateToken(c, request.Code)
	if err != nil {
		return nil, err
	}

	role := "user"
	usertype := "customer"

	auth, err := a.authrepo.NewAuth(c, token.UID, role, usertype, request)
	if err != nil {
		return nil, err
	}

	err = a.authrepo.Insert(c, auth)
	if err != nil {
		return nil, err
	}

	data, err := a.authrepo.GetCustomerbyID(c, request, auth.UID)
	if err != nil {
		return nil, err
	}

	tokenjwt, err := a.authrepo.GenerateToken(c, auth, data)
	if err != nil {
		return nil, err
	}
	b := tokenjwt.ToDto()

	return &b, nil
}

func (a AuthService) RefreshToken(c *gin.Context) (*dto.LoginResponse, error) {

	tokenjwt, err := a.authrepo.RefreshToken(c)
	if err != nil {
		return nil, err
	}

	d := tokenjwt.ToDto()

	return &d, nil
}

func NewAuthService(a domain.AuthRepository) domain.AuthService {
	return &AuthService{
		authrepo: a,
	}
}
