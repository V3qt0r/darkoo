package middleware

import (
	"log"
	"os"
	"strings"
	"time"
	"darkoo/models"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)


var identityKey = "id"

type User struct {
	ID uint
	UUID string
	UserName string
	Email string
}


type login struct {
	Email 		string `form:"email" json:"email" binding:"required"`
	Password 	string `form:"password" json:"password" binding:"required"`
}


func MiddleWare(userService models.IUserService) (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:			os.Getenv("REALM"),
		Key: 			[]byte(os.Getenv("SECRET")),
		Timeout:		time.Hour,
		MaxRefresh:		time.Hour,
		IdentityKey:	identityKey,
		PayloadFunc:    func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UUID,
				}
			}
			return jwt.MapClaims{}
		},

		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			n := string(claims[identityKey].(string))

			log.Printf("Parsed: %v\n", n)
			user, err := userService.GetUserByUUID(n)

			if err != nil {
				log.Printf("Middlerware couldnt get user by UUID %v\n", n)
				return err
			}

			return &User {
				UUID: 	user.UUID.String(),
				ID: 	user.ID,
			}
		},
 
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			userID := loginVals.Email
			password := loginVals.Password

			path := c.FullPath()

			var user *models.User
			var err error

			if strings.EqualFold("/api/user/one-time-login", path) {
				user, err = userService.LoginWithOneTimePassword(userID, password)
			} else {
				user, err = userService.Login(userID, password)
			}

			if err == nil {
				return &User {
					UUID: user.UUID.String(),
					UserName: user.UserName,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},

		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H {
				"code": code,
				"message": message,
			})
		},


		TokenLookup:  	"header: Authorization, query: token, cookie: jwt",
		TokenHeadName:	"Bearer",
		TimeFunc:		time.Now,
	})

	if err != nil {
		log.Print("Failed Authorization and Authentication!")
		return nil, err
	}

	return authMiddleware, nil
}