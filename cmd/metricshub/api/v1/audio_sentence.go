package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"video5pm-api/core/constants"
	"video5pm-api/models/services"
	transloadit "video5pm-api/pkg/tranloadit"

	"github.com/gin-gonic/gin"
	"github.com/go-audio/wav"
)

type createAudioForm struct {
	Text   string `form:"text"`
	Title  string `form:"title"`
	UserID int64  `form:"user_id"`
}

func CreateVideoPreview1(audioService *services.AudioService, videoService *services.VideoService) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var form createAudioForm

		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("%v", err.Error())})
			// c.HTML(http.StatusBadRequest, "index.html", gin.H{
			// 	"message": fmt.Sprintf("%v", constants.MSG_BAD_REQUEST),
			// })
			c.Abort()
			return
		}

		if len(form.Text) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Text is required",
			})
			c.Abort()
			return
		}

		if len(form.Title) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Title is required",
			})
			c.Abort()
			return
		}

		if form.UserID < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "user_id is required",
			})
			c.Abort()
			return
		}

		video, err := videoService.CreateVideoDefault(form.UserID, form.Title)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err,
			})
			return
		}

		listAudio := strings.Split(form.Text, ".")

		pathSubtitle := constants.SUBTITLE_PATH + "/" + form.Title + ".srt"
		f, err := os.Create(pathSubtitle)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err,
			})
			return
		}
		defer f.Close()

		var start int64 = 0

		for i, v := range listAudio {
			audio, err := audioService.CreateAudioDefault(video.ID, v)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": err,
				})
				return
			}

			file_path, err := getAudioSentence(v, strconv.FormatInt(audio.ID, 10))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": err,
				})
				return
			}

			time, err := getLengthAudio(file_path)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": err,
				})
				return
			}

			if time == -1 {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "the total size is not available",
				})
				return
			}

			err = audioService.AddPathAndLengthAudio(audio, file_path, time)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": err,
				})
				return
			}

			content := createContentSrtFile(i, start, start+time+500, v)
			_, err = f.WriteString(content)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": err,
				})
				return
			}
			start = start + time + 500

		}

		f.Sync()

		// err := packageService.UpdateUserPackage(form.Uid, form.Pid)
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"message": err,
		// 	})
		// 	return
		// }

		layer := transloadit.CreatLayoutPreview()

		c.JSON(http.StatusOK, gin.H{
			"status": layer,
		})
	}

	return gin.HandlerFunc(fn)
}

func getAudioSentence(text string, file_name string) (string, error) {

	requestBody, err := json.Marshal(map[string]string{
		"text":              text,
		"voice":             "doanngocle",
		"without_filter":    "false",
		"speed":             "1.0",
		"tts_return_option": "2",
	})

	if err != nil {
		return "", err
	}

	url := constants.VIETTEL_API

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Add("token", constants.TOKEN_VIETTEL)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New(strconv.Itoa(resp.StatusCode))
	}

	path := constants.AUDIO_PATH + "/" + file_name + ".mp3"
	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return path, nil
}

func getLengthAudio(file_path string) (int64, error) {
	r, err := os.Open(file_path)
	if err != nil {
		return 0, err
	}
	d := wav.NewDecoder(r)
	t, err := d.Duration()
	if err != nil {
		return 0, err
	}

	return int64(t.Milliseconds()), nil
}

func trimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

func convertIntToTime(i int64) string {
	timeParse := strconv.FormatInt(i, 10) + "ms"
	t, _ := time.ParseDuration(timeParse)
	h := strconv.FormatInt(int64(t.Hours()), 10)
	m := strconv.FormatInt(int64(t.Minutes()), 10)
	s := strconv.FormatInt(int64(t.Seconds()), 10)
	ms := strconv.FormatInt(int64(t.Milliseconds()), 10)

	if int64(t.Hours()) < 10 {
		h = "0" + h
	}

	if int64(t.Minutes()) < 10 {
		m = "0" + m
	}

	if int64(t.Seconds()) < 10 {
		s = "0" + s
	}

	if int64(t.Milliseconds()) < 10 {
		ms = "00" + ms
	} else if 10 <= int64(t.Milliseconds()) && int64(t.Milliseconds()) < 100 {
		ms = "0" + ms
	} else if int64(t.Milliseconds()) >= 1000 {
		ms = trimLastChar(ms)
	}

	duration := h + ":" + m + ":" + s + "," + ms
	return duration
}

func createContentSrtFile(i int, s int64, e int64, text string) string {

	index := strconv.Itoa(i + 1)
	start := convertIntToTime(s)
	end := convertIntToTime(e)

	content := index + "\n" + start + " --> " + end + "\n" + text + "\n"
	return content
}
