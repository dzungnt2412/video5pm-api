package api

import (
	v1 "video5pm-api/cmd/metricshub/api/v1"
	"video5pm-api/core/middleware"
	"video5pm-api/models/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

// InitRouter initialize routing information
func InitRouter(mysqlConn *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Authorization"},
		ExposeHeaders:    []string{"Content-Length, Authorization"},
		AllowCredentials: true,
	}))

	gin.SetMode(viper.GetString("server.run_mode"))

	secretKey := viper.GetString("server.secret_key")
	userService := services.NewUserService(mysqlConn)
	authService := services.NewAuthService(mysqlConn)
	packageService := services.NewPackageService(mysqlConn)
	featureService := services.NewFeatureService(mysqlConn)
	groupService := services.NewGroupService(mysqlConn)
	audioService := services.NewAudioService(mysqlConn)
	videoService := services.NewVideoService(mysqlConn)

	//auth
	r.POST("/auth", v1.LoginHeader(authService, featureService, groupService, secretKey))

	//group
	apiv1 := r.Group("/video5m/v1")
	apiv1.Use(middleware.MiddlewareJWT(secretKey))
	{
		apiv1.GET("/user", v1.FindUser(userService))
		apiv1.GET("/shop", v1.ShopPoint(userService))
		apiv1.POST("/package/update", v1.UpdateUserPackage(packageService))
		apiv1.POST("/package/updateVn", v1.UpdateUserVnPackage(packageService))
		apiv1.POST("/create-video-preview-1", v1.CreateVideoPreview1(audioService, videoService))
	}

	return r
}
