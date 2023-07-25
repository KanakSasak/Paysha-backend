package middleware

import (
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"time"
)

type Login struct {
	IdUser string
}

var identityKey = "id"
var id = ""

func NewAuth() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "Production",
		Key:              []byte("R@k1!ws@2@YWWbPd4435342dbRMXXBVLu##u!@taBoXxXx@3GLoxm[I}>s{t@@I@Fwhhkp79@@"),
		Timeout:          time.Hour * 24,
		MaxRefresh:       time.Hour * 8760,
		IdentityKey:      identityKey,
		SigningAlgorithm: "RS256",
		//PrivKeyFile:      "/Users/laluraynaldi/TestingApps/jwtrsakey/jwtRS256_jwk_key.key",
		//PubKeyFile:       "/Users/laluraynaldi/TestingApps/jwtrsakey/jwtRS256_jwk_key.pub",
		PrivKeyFile: "/usr/src/jwtRS256_jwk_key.key",
		PubKeyFile:  "/usr/src/jwtRS256_jwk_key.pub",

		Authenticator: func(c *gin.Context) (interface{}, error) {

			return nil, nil

		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": "Unauthorized!",
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc:     time.Now,
		SendCookie:   true,
		SecureCookie: true,
	})

}
