package services

import (
	"video5pm-api/models/entity"

	"github.com/jinzhu/gorm"
)

func NewAudioService(db *gorm.DB) *AudioService {
	return &AudioService{db: db}
}

type AudioService struct {
	db *gorm.DB
}

//UpdateUserPackage - service update package id of user
func (c *AudioService) CreateAudioDefault(video_id int64, text string) (*entity.Audio_sentence, error) {
	tx := c.db.Begin().LogMode(true)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	defaultAudioSentence := &entity.Audio_sentence{
		Video_id: video_id,
		Path:     "",
		Length:   0,
		Text:     text,
	}

	if err := tx.Model(entity.Audio_sentence{}).Create(defaultAudioSentence).Error; err != nil {
		tx.Rollback()
		return defaultAudioSentence, tx.Error
	}

	return defaultAudioSentence, tx.Commit().Error

}

func (c *AudioService) AddPathAndLengthAudio(audio_sentence *entity.Audio_sentence, path string, length int64) error {
	c.db.First(&audio_sentence)

	audio_sentence.Path = path
	audio_sentence.Length = length
	err := c.db.Save(&audio_sentence).Error

	if err != nil {
		return err
	}

	return nil
}
