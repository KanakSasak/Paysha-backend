package service

import (
	"github.com/gin-gonic/gin"
	"github.com/kaleido-io/kaleido-fabric-go/domain"
)

type CustomerService struct {
	CustomerRepo domain.CustomerRepository
}

func (c2 CustomerService) Minting(c *gin.Context, id string, amount string) error {
	err := c2.CustomerRepo.Minting(c, id, amount)
	if err != nil {
		return err
	}

	return nil
}

func (c2 CustomerService) GetWalletId(c *gin.Context, id string) (string, error) {
	data, err := c2.CustomerRepo.GetWalletId(c, id)
	if err != nil {
		return "", err
	}

	resp := data
	return resp, nil
}

func (c2 CustomerService) GetWalletBalance(c *gin.Context, id string) (string, error) {
	data, err := c2.CustomerRepo.GetWalletBalance(c, id)
	if err != nil {
		return "", err
	}

	resp := data
	return resp, nil
}

func (c2 CustomerService) Payment(c *gin.Context, id string, idtujuan string, amount string) error {
	err := c2.CustomerRepo.Payment(c, id, idtujuan, amount)
	if err != nil {
		return err
	}

	return nil
}

func (c2 CustomerService) FetchByID(c *gin.Context, id string) (interface{}, error) {
	data, err := c2.CustomerRepo.FetchByID(c, id)
	if err != nil {
		return nil, err
	}

	resp := data.ToClientDto()
	return &resp, nil

}

func NewCustomerService(repo domain.CustomerRepository) domain.CustomerService {
	return &CustomerService{
		CustomerRepo: repo,
	}
}
