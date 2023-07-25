package repository

import (
	"errors"
	"firebase.google.com/go/auth"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/kaleido-io/kaleido-fabric-go/domain"
	"github.com/kaleido-io/kaleido-fabric-go/dto"
	"github.com/kaleido-io/kaleido-fabric-go/fabric"
	"github.com/kaleido-io/kaleido-fabric-go/kaleido"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"os"

	"strings"
	"time"
)

type AuthRepository struct {
	client   *auth.Client
	clientdb *gorm.DB
	jwtAuth  *jwt.GinJWTMiddleware
	config   map[string]interface{}
	network  *kaleido.KaleidoNetwork
}

func (a AuthRepository) Register(c *gin.Context, request dto.RegisterRequest) (*domain.Customer, error) {
	uid, err := a.CreateUserFirebase(c, request)
	if err != nil {
		return nil, err
	}

	wallet, err := a.CreateUserFabric(c, uid)
	if err != nil {
		return nil, err
	}

	resp, err := a.CreateCustomer(c, request, uid, wallet)
	if err != nil {
		return nil, err
	}

	return resp, nil

}

func (a AuthRepository) CreateUserFabric(c *gin.Context, uid string) (walletresp string, err error) {
	sdk1 := newSDK(a.config)
	defer sdk1.Close()
	log.Println(a.network)
	log.Println(a.config)
	wallet := kaleido.NewWallet(uid, *a.network, sdk1)
	err = wallet.InitIdentity()
	if err != nil {
		return "", err
	}
	fabric.AddTlsConfig(a.config, wallet.Signer)

	sdk2 := newSDK(a.config)
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
	data, err := channel.ExecChaincode(ccname, "account", "", "")
	if err != nil {
		fmt.Printf("=> Failed to invoke chaincode %s\n", err)
	} else {
		fmt.Printf("=> Completed to invoke chaincode \n")
	}

	data2, err := channel.ExecChaincode(ccname, "mint", "", "1")
	if err != nil {
		fmt.Printf("=> Failed to invoke chaincode %s\n", err)
	} else {
		fmt.Printf("=> Completed to invoke chaincode \n")
	}

	fmt.Printf("\nAll Done!\n")

	log.Println(data2)

	return data, nil
}

func newSDK(config map[string]interface{}) *fabsdk.FabricSDK {
	configProvider, err := fabric.NewConfigProvider(config)
	if err != nil {
		fmt.Printf("Failed to create config provider from config map: %s\n", err)
		os.Exit(1)
	}

	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		fmt.Printf("Failed to instantiate an SDK: %s\n", err)
		os.Exit(1)
	}
	return sdk
}

func (a AuthRepository) GetCustomerbyID(c *gin.Context, request dto.LoginRequest, id string) (data *domain.Customer, err error) {
	if err := a.clientdb.Table("customers").Where("fire_uid = ?", id).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := a.client.DeleteUser(c, id)
			if err != nil {
				log.Println("error deleting user: %v\n", err)
			}
			return nil, errors.New("customer not found please register first")
		}
	} else if err != nil {
		return nil, errors.New("something wrong when get data customer")
	}

	return data, nil
}

func (a AuthRepository) RefreshToken(c *gin.Context) (*domain.Token, error) {
	token, exp, err := a.jwtAuth.RefreshToken(c)
	if err != nil {
		return nil, err
	}
	b := &domain.Token{
		Expires: exp,
		Token:   token,
	}

	return b, err
}

func (a AuthRepository) CreateUserFirebase(c *gin.Context, request dto.RegisterRequest) (string, error) {

	params := (&auth.UserToCreate{}).
		DisplayName(request.Name).
		Email(request.Email).
		Password(request.Password).
		PhoneNumber(request.Phone).
		EmailVerified(false).
		Disabled(false)
	data, err := a.client.CreateUser(c, params)
	if err != nil {
		log.Println("error creating user: %v\n", err)
		return "", err
	}

	return data.UID, nil

}

func (a AuthRepository) CreateCustomer(c *gin.Context, request dto.RegisterRequest, uid string, wallet string) (*domain.Customer, error) {

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return nil, err
	}
	t := time.Now().Local().In(loc)
	dateval := t.Format("2006-01-02")
	timeval := t.Format("15:04:05")

	data := &domain.Customer{
		FireUid:        uid,
		Name:           request.Name,
		Phone:          request.Phone,
		Email:          request.Email,
		Password:       "",
		Address:        "",
		Role:           "customer",
		Wallet:         wallet,
		CreateDate:     dateval,
		CreateTime:     timeval,
		CreatedBy:      "admin",
		LastChangeBy:   "admin",
		LastUpdateDate: dateval,
		LastUpdateTime: timeval,
		Active:         true,
	}

	// begin a transaction
	tx := a.clientdb.Begin()

	result := tx.Table("customers").Create(&data)
	if result.Error != nil {
		tx.Rollback()
		return nil, errors.New("something wrong when save data user")
	}

	tx.Commit()

	return data, nil

}

func (a AuthRepository) GenerateToken(c *gin.Context, data *domain.Auth, dataCustomer *domain.Customer) (*domain.Token, error) {
	a.jwtAuth.PayloadFunc = func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*domain.Auth); ok {
			return jwt.MapClaims{
				"id":        v.AuthID,
				"reference": v.UID,
				"role":      "user",
				"audience":  dataCustomer.ID,
				"wallet":    dataCustomer.Wallet,
			}
		}
		return jwt.MapClaims{}
	}
	token, exp, err := a.jwtAuth.TokenGenerator(data)
	if err != nil {
		return nil, err
	}

	b := &domain.Token{
		Expires: exp,
		Token:   token,
		Data:    dataCustomer,
	}

	return b, nil
}

func (a AuthRepository) NewAuth(c *gin.Context, uid string, role string, usertype string, data dto.LoginRequest) (*domain.Auth, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return nil, err
	}
	t := time.Now().Local().In(loc)
	dateval := t.Format("2006-01-02")
	timeval := t.Format("15:04:05")

	b := &domain.Auth{
		UID:           uid,
		SigningMethod: data.SigningMethod,
		Email:         data.Email,
		Phone:         data.Phone,
		Password:      data.Password,
		Role:          role,
		UserType:      usertype,
		TokenFcm:      data.TokenFcm,
		DeviceType:    data.DeviceType,
		DeviceID:      data.DeviceID,
		LastLoginDate: dateval,
		LastLoginTime: timeval,
	}

	return b, nil
}

func (a AuthRepository) Insert(c *gin.Context, data *domain.Auth) error {

	result := a.clientdb.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"token_fcm": &data.TokenFcm}),
	}).Create(&data)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (a AuthRepository) ValidateToken(c *gin.Context, code string) (*auth.Token, error) {

	tokenF, err := a.client.VerifyIDToken(c, code)
	if err != nil {
		return nil, errors.New("firebase verification failed")
	}

	tokenFirebase := tokenF

	return tokenFirebase, err
}

func NewAuthDBRepository(authClient *auth.Client, dbClient *gorm.DB, jwtAuth *jwt.GinJWTMiddleware, config map[string]interface{}, network *kaleido.KaleidoNetwork) domain.AuthRepository {
	return &AuthRepository{authClient, dbClient, jwtAuth, config, network}
}
