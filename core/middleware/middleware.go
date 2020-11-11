package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"strings"
	"video5pm-api/core/constants"
	"video5pm-api/core/utils"
	"video5pm-api/pkg/logger"
)

// MiddlewareJWT is jwt middleware
func MiddlewareJWT(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ok = false
		// var ok = true

		//read token from cookie
		// token, err := c.Cookie("lionnix_token")
		// if err != nil {
		// 	logger.Log.Errorf("Cannot get Auth Cookies")
		// 	c.JSON(http.StatusUnauthorized, gin.H{
		// 		"message": constants.MSG_PERMISSION_DENIED,
		// 	})
		// 	c.Abort()
		// 	return
		// }

		token := c.Request.Header.Get("Authorization")
		// token := strings.Split(cookie, constants.TOKEN_PREFIX)[1]

		// read referer from header
		referer := c.Request.Header.Get("referer")

		if token == "" {
			logger.Log.Infof("Invalid Token: Null")
		} else {
			claims, err := utils.ParseToken(token, secretKey)
			if err != nil {
				logger.Log.Infof("Invalid Token: %v", err)
			} else {
				if int64(claims["user_id"].(float64)) > 0 {
					ok = true

					if claims["references"] != nil {
						references := claims["references"].([]interface{})
						for _, r := range references {
							if strings.Contains(referer, r.(string)) {
								ok = true
							}
						}
					}

				}
			}
		}

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": constants.MSG_PERMISSION_DENIED,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
