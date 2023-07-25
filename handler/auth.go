package handler

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/kaleido-io/kaleido-fabric-go/domain"
	"github.com/kaleido-io/kaleido-fabric-go/dto"
	"log"
	"net/http"
)

type AuthHandler struct {
	AuthService domain.AuthService
	JwtAuth     *jwt.GinJWTMiddleware
}

// LoginCust
// @Tags Auth
// @Summary Auth Customer
// @Produce  json
// @Param Auth body dto.LoginRequest true "Login"</span>
// @success 201 {object} dto.LoginResponse
// @Failure 400,401 {object} string	"error"
// @Failure 500 {object} string	"error"
// @Router /v1/cust/auth/login [post]
func (h AuthHandler) LoginCust(c *gin.Context) {
	var request dto.LoginRequest
	//body, _ := ioutil.ReadAll(c.Request.Body)
	//println(string(body))
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("masukk", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("masukk")

	response, err := h.AuthService.LoginCust(c, request)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(
			200,
			gin.H{"data": response},
		)
	}
}

func (h AuthHandler) WebLoginCust(c *gin.Context) {
	var request dto.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.AuthService.LoginCust(c, request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(
			200,
			gin.H{"data": response},
		)
	}
}

// RegisterCust
// @Tags Auth
// @Summary Auth Customer
// @Produce  json
// @Param Auth body dto.RegisterRequest true "Register"</span>
// @success 201 {object} dto.RegisterResponse
// @Failure 400,401 {object} string	"error"
// @Failure 500 {object} string	"error"
// @Router /v1/cust/auth/register [post]
func (h AuthHandler) RegisterCust(c *gin.Context) {
	var request dto.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.AuthService.Register(c, request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(
			200,
			gin.H{"data": response},
		)
	}
}

// Refresh
// @Tags Auth
// @Summary Auth Refresh
// @Security ApiKeyAuth
// @Produce  application/json
// @success 201 {object} dto.LoginResponse
// @Failure 400,401 {object} string	"error"
// @Failure 500 {object} string	"error"
// @Router /v1/cust/auth/refresh [get]
func (h AuthHandler) Refresh(c *gin.Context) {
	response, err := h.AuthService.RefreshToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(
			200,
			gin.H{"data": response},
		)
	}
}

func NewAuthHandler(r *gin.RouterGroup, us domain.AuthService, jwtAuth *jwt.GinJWTMiddleware) {
	handler := &AuthHandler{
		AuthService: us,
		JwtAuth:     jwtAuth,
	}
	r.POST("/cust/auth/login", handler.LoginCust)
	r.POST("/web/auth/login", handler.WebLoginCust)
	r.POST("/cust/auth/register", handler.RegisterCust)
	r.GET("/cust/auth/refresh", jwtAuth.RefreshHandler)
}
