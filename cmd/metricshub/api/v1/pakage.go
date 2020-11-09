package v1

import (
	"fmt"
	"net/http"
	"video5pm-api/models/services"

	"github.com/gin-gonic/gin"
)

type updatePackageForm struct {
	Uid int64 `form:"uid" binding:"required,numeric"`
	Pid int64 `form:"pid" binding:"required,numeric"`
}

type updatePackageVnForm struct {
	Uid int64 `form:"uid" binding:"required,numeric"`
	Pid int64 `form:"pid" binding:"required,numeric"`
}

func UpdateUserPackage(packageService *services.PackageService) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var form updatePackageForm

		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("%v", err.Error())})
			// c.HTML(http.StatusBadRequest, "index.html", gin.H{
			// 	"message": fmt.Sprintf("%v", constants.MSG_BAD_REQUEST),
			// })
			c.Abort()
			return
		}

		if form.Pid < 1 || form.Pid > 3 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid Package id",
			})
			c.Abort()
			return
		}

		err := packageService.UpdateUserPackage(form.Uid, form.Pid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "true",
		})
	}

	return gin.HandlerFunc(fn)
}

func UpdateUserVnPackage(packageService *services.PackageService) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var form updatePackageVnForm

		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("%v", err.Error())})
			// c.HTML(http.StatusBadRequest, "index.html", gin.H{
			// 	"message": fmt.Sprintf("%v", constants.MSG_BAD_REQUEST),
			// })
			c.Abort()
			return
		}

		if form.Pid < 6 || form.Pid > 7 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid Package id",
			})
			c.Abort()
			return
		}

		err := packageService.UpdateUserVnPackage(form.Uid, form.Pid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "true",
		})
	}

	return gin.HandlerFunc(fn)
}
