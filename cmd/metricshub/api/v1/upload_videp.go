package v1

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"video5pm-api/core/constants"
	"video5pm-api/models/services"
	transloadit "video5pm-api/pkg/tranloadit"

	"github.com/gin-gonic/gin"
)

func UploadVideo(audioService *services.AudioService, videoService *services.VideoService) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		VideoID := c.PostForm("video_id")
		l := c.PostForm("length")

		Length, err := time.ParseDuration(l + "s")

		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"message": "Length is required",
		// 	})
		// 	c.Abort()
		// 	return
		// }

		if l == "0" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "length is required",
			})
			c.Abort()
			return
		}
		if len(VideoID) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "VideoID is required",
			})
			c.Abort()
			return
		}

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "1",
			})
			c.Abort()
			return
		}

		filename := VideoID + header.Filename
		path := constants.VIDEO_SENTENCE_PATH + "/" + filename
		out, err := os.Create(path)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "2",
			})
			c.Abort()
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "3",
			})
			c.Abort()
			return
		}

		fmt.Print("done")
		videoID, _ := strconv.ParseInt(VideoID, 10, 64)
		video, err := videoService.FindVideo(videoID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "4",
			})
			c.Abort()
			return
		}
		videoPreview, err := videoService.FindVideoPreview(videoID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "5",
			})
			c.Abort()
			return
		}

		timeVideoPreview := strconv.FormatInt(videoPreview.Length, 10) + "ms"
		t1, _ := time.ParseDuration(timeVideoPreview)

		timeVideo := strconv.FormatInt(video.Length, 10) + "ms"
		t2, _ := time.ParseDuration(timeVideo)

		v, vp, lengthPreview := transloadit.Concatenate_video(path, Length, videoPreview.Path, t1, t2, video.Path_audio, video.Subtitle)

		new_video, err := videoService.AddPathVideo(&video, v, 0, "")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "5",
			})
			c.Abort()
			return
		}

		_, err = videoService.AddPathLengthToVideoPreview(&videoPreview, vp, lengthPreview)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "5",
			})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"video": new_video,
		})
	}

	return gin.HandlerFunc(fn)
}
