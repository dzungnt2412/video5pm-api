package entity

import (
	"video5pm-api/pkg/database"
)

//User - model user add groups and references
type Audio_sentence struct {
	database.Model
	Video_id int64  `json:"video_id" gorm:"column:video_id"`
	Path     string `json:"path"`
	Length   int64  `json:"Length"`
	Text     string `json:"text"`
}
