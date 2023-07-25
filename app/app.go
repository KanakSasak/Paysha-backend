package app

import (
	"github.com/gin-gonic/gin"
	"github.com/kaleido-io/kaleido-fabric-go/database"
	_ "github.com/kaleido-io/kaleido-fabric-go/docs"
	"github.com/kaleido-io/kaleido-fabric-go/handler"
	"github.com/kaleido-io/kaleido-fabric-go/middleware"
	"github.com/kaleido-io/kaleido-fabric-go/repository"
	"github.com/kaleido-io/kaleido-fabric-go/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func Start() {
	gin.ForceConsoleColor()
	gin.SetMode(gin.DebugMode)

	r := gin.Default()
	r.Use(CORSMiddleware())

	database.FirebaseConnect()
	database.PostgresConnect()
	database.BlockchainConnect()

	authfb := database.GetfirebaseAuth()
	dbpostg := database.GetPostgresdb()
	blockchainConfig := database.GetBlockchainConfig()
	blockchainNetwork := database.GetBlockchainNetwork()

	api := r.Group("/v1")

	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	jwtauth, err := middleware.NewAuth()
	if err != nil {
		panic("JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := jwtauth.MiddlewareInit()

	if errInit != nil {
		panic("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	AuthRepo := repository.NewAuthDBRepository(authfb, dbpostg, jwtauth, blockchainConfig, blockchainNetwork)
	AuthService := service.NewAuthService(AuthRepo)

	CustomerRepo := repository.NewCustomerRepository(dbpostg, jwtauth, blockchainConfig, blockchainNetwork)
	CustomerService := service.NewCustomerService(CustomerRepo)

	handler.NewAuthHandler(api, AuthService, jwtauth)
	api.Use(jwtauth.MiddlewareFunc())
	{
		handler.NewCustomerHandler(api, CustomerService, jwtauth)
	}

	log.Fatal(r.Run(":8888"))
	log.Println("listen on port :8888")
}
