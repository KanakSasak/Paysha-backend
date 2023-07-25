package repository

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/kaleido-io/kaleido-fabric-go/domain"
	"github.com/kaleido-io/kaleido-fabric-go/fabric"
	"github.com/kaleido-io/kaleido-fabric-go/kaleido"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"strings"
)

type CustomerRepository struct {
	client  *gorm.DB
	jwtAuth *jwt.GinJWTMiddleware
	config  map[string]interface{}
	network *kaleido.KaleidoNetwork
}

func (c2 CustomerRepository) Minting(c *gin.Context, id string, amount string) error {
	var data *domain.Customer
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	// Get record by id
	result := c2.client.Debug().Table("customers").Where("id=?", idInt).Find(&data)
	if result.Error != nil {
		return result.Error
	}

	sdk1 := newSDK(c2.config)
	defer sdk1.Close()

	wallet := kaleido.NewWallet(data.FireUid, *c2.network, sdk1)
	err = wallet.InitIdentity()
	if err != nil {
		return err
	}
	fabric.AddTlsConfig(c2.config, wallet.Signer)

	sdk2 := newSDK(c2.config)
	defer sdk2.Close()

	channel := kaleido.NewChannel("default-channel", sdk2)
	err = channel.Connect(wallet.Signer.Identifier())
	if err != nil {
		fmt.Printf("Failed to connect to channel: %s\n", err)
		os.Exit(1)
	}

	initChaincode := os.Getenv("INIT_CC")
	if initChaincode == "" {
		initChaincode = "false"
	} else {
		initChaincode = strings.ToLower(initChaincode)
	}
	ccname := "payshatest"
	resp, err := channel.ExecChaincode(ccname, "mint", "", amount)
	if err != nil {
		fmt.Printf("=> Failed to invoke chaincode %s\n", err)
	} else {
		fmt.Printf("=> Completed to invoke chaincode \n")
	}

	log.Println(resp)

	fmt.Printf("\nAll Done!\n")

	return nil
}

func (c2 CustomerRepository) GetWalletId(c *gin.Context, id string) (string, error) {
	var data *domain.Customer
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return "", err
	}

	// Get record by id
	result := c2.client.Debug().Table("customers").Where("id=?", idInt).Find(&data)
	if result.Error != nil {
		return "", result.Error
	}

	sdk1 := newSDK(c2.config)
	defer sdk1.Close()

	wallet := kaleido.NewWallet(data.FireUid, *c2.network, sdk1)
	err = wallet.InitIdentity()
	if err != nil {
		return "", err
	}
	fabric.AddTlsConfig(c2.config, wallet.Signer)

	sdk2 := newSDK(c2.config)
	defer sdk2.Close()

	channel := kaleido.NewChannel("default-channel", sdk2)
	err = channel.Connect(wallet.Signer.Identifier())
	if err != nil {
		fmt.Printf("Failed to connect to channel: %s\n", err)
		os.Exit(1)
	}

	initChaincode := os.Getenv("INIT_CC")
	if initChaincode == "" {
		initChaincode = "false"
	} else {
		initChaincode = strings.ToLower(initChaincode)
	}
	ccname := "payshatest"
	resp, err := channel.ExecChaincode(ccname, "account", "", "")
	if err != nil {
		fmt.Printf("=> Failed to invoke chaincode %s\n", err)
	} else {
		fmt.Printf("=> Completed to invoke chaincode \n")
	}

	fmt.Printf("\nAll Done!\n")

	return resp, nil
}

func (c2 CustomerRepository) GetWalletBalance(c *gin.Context, id string) (string, error) {
	var data *domain.Customer
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return "", err
	}

	// Get record by id
	result := c2.client.Debug().Table("customers").Where("id=?", idInt).Find(&data)
	if result.Error != nil {
		return "", result.Error
	}

	sdk1 := newSDK(c2.config)
	defer sdk1.Close()

	wallet := kaleido.NewWallet(data.FireUid, *c2.network, sdk1)
	err = wallet.InitIdentity()
	if err != nil {
		return "", err
	}
	fabric.AddTlsConfig(c2.config, wallet.Signer)

	sdk2 := newSDK(c2.config)
	defer sdk2.Close()

	channel := kaleido.NewChannel("default-channel", sdk2)
	err = channel.Connect(wallet.Signer.Identifier())
	if err != nil {
		fmt.Printf("Failed to connect to channel: %s\n", err)
		os.Exit(1)
	}

	initChaincode := os.Getenv("INIT_CC")
	if initChaincode == "" {
		initChaincode = "false"
	} else {
		initChaincode = strings.ToLower(initChaincode)
	}
	ccname := "payshatest"
	resp, err := channel.ExecChaincode(ccname, "clientaccountbalance", "", "")
	if err != nil {
		fmt.Printf("=> Failed to invoke chaincode %s\n", err)
	} else {
		fmt.Printf("=> Completed to invoke chaincode \n")
	}

	fmt.Printf("\nAll Done!\n")

	return resp, nil
}

func (c2 CustomerRepository) Payment(c *gin.Context, id string, idtujuan string, amount string) error {
	var data *domain.Customer
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	// Get record by id
	result := c2.client.Debug().Table("customers").Where("id=?", idInt).Find(&data)
	if result.Error != nil {
		return result.Error
	}

	sdk1 := newSDK(c2.config)
	defer sdk1.Close()

	wallet := kaleido.NewWallet(data.FireUid, *c2.network, sdk1)
	err = wallet.InitIdentity()
	if err != nil {
		return err
	}
	fabric.AddTlsConfig(c2.config, wallet.Signer)

	sdk2 := newSDK(c2.config)
	defer sdk2.Close()

	channel := kaleido.NewChannel("default-channel", sdk2)
	err = channel.Connect(wallet.Signer.Identifier())
	if err != nil {
		fmt.Printf("Failed to connect to channel: %s\n", err)
		os.Exit(1)
	}

	initChaincode := os.Getenv("INIT_CC")
	if initChaincode == "" {
		initChaincode = "false"
	} else {
		initChaincode = strings.ToLower(initChaincode)
	}
	ccname := "payshatest"
	resp, err := channel.ExecChaincode(ccname, "transfer", idtujuan, amount)
	if err != nil {
		fmt.Printf("=> Failed to invoke chaincode %s\n", err)
	} else {
		fmt.Printf("=> Completed to invoke chaincode \n")
	}

	log.Println(resp)
	fmt.Printf("\nAll Done!\n")

	return nil
}

func (c2 CustomerRepository) FetchByID(c *gin.Context, id string) (*domain.Customer, error) {
	var data *domain.Customer

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	// Get record by id
	result := c2.client.Debug().Table("customers").Where("id=?", idInt).Find(&data)
	if result.Error != nil {
		return nil, result.Error
	}

	return data, nil
}

func NewCustomerRepository(client *gorm.DB, jwtAuth *jwt.GinJWTMiddleware, config map[string]interface{}, network *kaleido.KaleidoNetwork) domain.CustomerRepository {
	return &CustomerRepository{client, jwtAuth, config, network}
}
