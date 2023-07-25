package handler

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/kaleido-io/kaleido-fabric-go/domain"

	"net/http"
)

type CustomerHandler struct {
	CustomerService domain.CustomerService
	JwtAuth         *jwt.GinJWTMiddleware
}

// GetCustomer
// @Tags Customer
// @Summary Get Customer
// @Accept  json
// @Security ApiKeyAuth
// @Produce  application/json
// @Param id path string true "Get"
// @success 201 {object} dto.CustAdminResponse
// @Failure 400,401 {object} string	"error"
// @Failure 500 {object} string	"error"
// @Router /v1/cust/get [get]
func (h CustomerHandler) GetCustomer(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	id := fmt.Sprint(claims["audience"])

	data, err := h.CustomerService.FetchByID(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(
			200,
			gin.H{"data": data},
		)
	}
}

// GetWalletId
// @Tags Customer
// @Summary Wallet Id
// @Accept  json
// @Security ApiKeyAuth
// @Produce  application/json
// @Param id path string true "Get"
// @success 201 {object} string "eDUwOTo6Q049dXNlcjMsT1U9Y2xpZW50OjpDTj1mYWJyaWMtY2Etc2VydmVy"
// @Failure 400,401 {object} string	"error"
// @Failure 500 {object} string	"error"
// @Router /v1/cust/get/wallet [get]
func (h CustomerHandler) GetWalletId(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	id := fmt.Sprint(claims["audience"])

	data, err := h.CustomerService.GetWalletId(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(
			200,
			gin.H{"data": data},
		)
	}
}

// GetWalletBalance
// @Tags Customer
// @Summary Wallet Balance
// @Accept  json
// @Security ApiKeyAuth
// @Produce  application/json
// @Param id path string true "Get"
// @success 201 {object} string "100000"
// @Failure 400,401 {object} string	"error"
// @Failure 500 {object} string	"error"
// @Router /v1/cust/get/wallet/balance [get]
func (h CustomerHandler) GetWalletBalance(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	id := fmt.Sprint(claims["audience"])

	data, err := h.CustomerService.GetWalletBalance(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(
			200,
			gin.H{"data": data},
		)
	}
}

// Payment
// @Tags Customer
// @Summary Payment
// @Accept  json
// @Security ApiKeyAuth
// @Produce  application/json
// @Param id path string true "POST"
// @success 201 {object} string "Pembayaran Sukses"
// @Failure 400,401 {object} string	"error"
// @Failure 500 {object} string	"error"
// @Router /v1/cust/bayar/:idtujuan/:amount [POST]
func (h CustomerHandler) Payment(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	id := fmt.Sprint(claims["audience"])

	idtujuan := c.Param("idtujuan")
	amount := c.Param("amount")

	err := h.CustomerService.Payment(c, id, idtujuan, amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(
			200,
			gin.H{"data": "Pembayaran Sukses"},
		)
	}
}

func (h CustomerHandler) Minting(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	id := fmt.Sprint(claims["audience"])

	amount := c.Param("amount")

	err := h.CustomerService.Minting(c, id, amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(
			200,
			gin.H{"data": "Pembayaran Sukses"},
		)
	}
}

func NewCustomerHandler(r *gin.RouterGroup, us domain.CustomerService, jwtAuth *jwt.GinJWTMiddleware) {

	handler := &CustomerHandler{
		CustomerService: us,
		JwtAuth:         jwtAuth,
	}

	r.GET("/cust/get", handler.GetCustomer)
	r.GET("/cust/get/wallet", handler.GetWalletId)
	r.GET("/cust/get/wallet/balance", handler.GetWalletBalance)
	r.POST("/cust/bayar/:idtujuan/:amount", handler.Payment)
	r.POST("/cust/mint/:amount", handler.Minting)

}
