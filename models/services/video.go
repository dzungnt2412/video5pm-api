package services

import (
	"video5pm-api/models/entity"

	"github.com/jinzhu/gorm"
)

func NewVideoService(db *gorm.DB) *VideoService {
	return &VideoService{db: db}
}

type VideoService struct {
	db *gorm.DB
}

//UpdateUserPackage - service update package id of user
func (c *VideoService) CreateVideoDefault(user_id int64, title string) (*entity.Video, error) {
	tx := c.db.Begin().LogMode(true)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	defaultVideo := &entity.Video{
		User_id:  user_id,
		Title:    title,
		Path:     "",
		Status:   "",
		Subtitle: "",
	}

	if err := tx.Model(entity.Video{}).Create(defaultVideo).Error; err != nil {
		tx.Rollback()
		return defaultVideo, tx.Error
	}

	return defaultVideo, tx.Commit().Error

}

func (c *VideoService) CreateVideoPreviewDefault(video_id int64) (*entity.Video_previews, error) {
	tx := c.db.Begin().LogMode(true)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	defaultVideo := &entity.Video_previews{
		Video_id: video_id,
		Length:   0,
		Path:     "",
	}

	if err := tx.Model(entity.Video_previews{}).Create(defaultVideo).Error; err != nil {
		tx.Rollback()
		return defaultVideo, tx.Error
	}

	return defaultVideo, tx.Commit().Error

}

func (c *VideoService) AddPathAudioToVideo(video *entity.Video, path string) error {
	c.db.First(&video)

	video.Path_audio = path
	err := c.db.Save(&video).Error

	if err != nil {
		return err
	}

	return nil
}

func (c *VideoService) AddPathVideo(video *entity.Video, path string, length int64, sub string) (*entity.Video, error) {
	c.db.First(&video)

	video.Path = path
	if length != 0 {
		video.Length = length
	}
	if sub != "" {
		video.Subtitle = sub
	}
	err := c.db.Save(&video).Error

	if err != nil {
		return nil, err
	}

	return video, nil
}

func (c *VideoService) AddPathLengthToVideoPreview(video *entity.Video_previews, path string, length int64) (*entity.Video_previews, error) {
	c.db.First(&video)

	video.Path = path
	if length != 0 {
		video.Length = length
	}
	err := c.db.Save(&video).Error

	if err != nil {
		return nil, err
	}

	return video, nil
}

func (c *VideoService) FindVideo(id int64) (entity.Video, error) {
	var video entity.Video
	db := c.db.LogMode(true)
	err := db.First(&video, "id=?", id).Error
	if err != nil {
		return video, err
	}
	return video, nil
}

func (c *VideoService) FindVideoPreview(id int64) (entity.Video_previews, error) {
	var videoPreview entity.Video_previews
	db := c.db.LogMode(true)
	err := db.First(&videoPreview, "video_id = ?", id).Error
	if err != nil {
		return videoPreview, err
	}
	return videoPreview, nil
}
