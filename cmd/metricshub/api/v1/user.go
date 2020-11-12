package v1

import (
	"net/http"
	"strconv"
	"video5pm-api/models/services"

	"github.com/gin-gonic/gin"
)

//Find user - function find and return user by id
func FindUser(userService *services.UserService) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		uid, err := strconv.ParseInt(c.Query("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid UserID",
			})
			return

		}

		u, err := userService.FindUser(uid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid UserID",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": u,
		})
	}

	return gin.HandlerFunc(fn)
}
