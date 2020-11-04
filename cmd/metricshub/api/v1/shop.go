package v1

import (
	"fmt"
	"lionnix-metrics-api/models/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShopPoint(userService *services.UserService) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		username := c.Query("username")
		fmt.Printf("%v", username)

		u, err := userService.FindUserByUsername(username)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Username does not exist",
			})
			return
		}

		svn, err := userService.GetShopPointVN(username)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Cannot found shop in VN",
			})
			return
		}

		sg, err := userService.GetShopPointGlobal(username)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Cannot found shop in Global",
			})
			return
		}

		fmt.Println(u, svn, sg)

		c.JSON(http.StatusOK, gin.H{
			"user":       u,
			"shopVn":     svn,
			"shopGlobal": sg,
		})
	}

	return gin.HandlerFunc(fn)
}
