package v1

import (
	"net/http"
	"video5pm-api/core/constants"
	"video5pm-api/core/utils"
	"video5pm-api/models/services"

	"github.com/gin-gonic/gin"
)

//loginHeaderForm - struct parse form
type loginHeaderForm struct {
	Username string `form:"username" binding:"required,min=6"`
	Password string `form:"password" binding:"required,min=6"`
}

//LoginHeader - parse and authenticate form
func LoginHeader(authService *services.AuthService, featureService *services.FeatureService, groupService *services.GroupService, secretKey string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var form loginHeaderForm

		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Input Data"})
			c.Abort()
			return
		}

		user, err := authService.GetUserByUserName(form.Username)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid Usernam or Password",
			})
			return
		}

		match := utils.CheckPasswordHash(form.Password, user.Password)
		if !match {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid Usernam or Password",
			})
			return
		}

		groups := groupService.GetGroupIDByUserID(user.ID)
		features := featureService.GetFeaturesByUserID(user.ID)
		references := featureService.GetListReferences(features)

		token, err := utils.GenerateToken(user, groups, references, secretKey)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Error",
			})
			return
		}

		c.Header("Authorization", constants.TOKEN_PREFIX+token)

		c.JSON(http.StatusOK, gin.H{
			"access_token": token,
			"user":         user,
		})
	}

	return gin.HandlerFunc(fn)
}
