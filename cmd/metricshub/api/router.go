package api

import (
	v1 "video5pm-api/cmd/metricshub/api/v1"
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

	audioService := services.NewAudioService(mysqlConn)
	videoService := services.NewVideoService(mysqlConn)

	//auth
	r.POST("/create-video-preview", v1.CreateVideoPreview(audioService, videoService))
	r.POST("/upload-video", v1.UploadVideo(audioService, videoService))

	return r
}
