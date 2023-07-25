package domain

import (
	"github.com/gin-gonic/gin"
	"github.com/kaleido-io/kaleido-fabric-go/dto"
)

type Customer struct {
	ID             uint64 `json:"id" gorm:"primaryKey;autoIncrement:true"`
	FireUid        string `json:"fire_uid"`
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Address        string `json:"address"`
	Role           string `json:"role"`
	Wallet         string `json:"wallet"`
	CreateDate     string `json:"create_date"`
	CreateTime     string `json:"create_time"`
	CreatedBy      string `json:"created_by"`
	LastChangeBy   string `json:"last_change_by"`
	LastUpdateDate string `json:"last_update_date"`
	LastUpdateTime string `json:"last_update_time"`
	Active         bool   `json:"active"`
}

func (d Customer) ToClientDto() dto.CustClientResp {
	return dto.CustClientResp{
		ID:      d.ID,
		FireUid: d.FireUid,
		Name:    d.Name,
		Phone:   d.Phone,
		Email:   d.Email,
		Address: d.Address,
		Role:    d.Role,
		Wallet:  d.Wallet,
		Active:  d.Active,
	}
}

type CustomerRepository interface {
	FetchByID(c *gin.Context, id string) (*Customer, error)
	GetWalletId(c *gin.Context, id string) (string, error)
	GetWalletBalance(c *gin.Context, id string) (string, error)
	Payment(c *gin.Context, id string, idtujuan string, amount string) error
	Minting(c *gin.Context, id string, amount string) error
}

type CustomerService interface {
	FetchByID(c *gin.Context, id string) (interface{}, error)
	GetWalletId(c *gin.Context, id string) (string, error)
	GetWalletBalance(c *gin.Context, id string) (string, error)
	Payment(c *gin.Context, id string, idtujuan string, amount string) error
	Minting(c *gin.Context, id string, amount string) error
}
