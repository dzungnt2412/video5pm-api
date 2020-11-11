package entity

import "video5pm-api/pkg/database"

//User - model user add groups and references
type Video struct {
	database.Model
	User_id    int64  `json:"user_id" gorm:"column:user_id"`
	Title      string `json:"title"`
	Path       string `json:"path"`
	Status     string `json:"status"`
	Subtitle   string `json:"subtitle"`
	Path_audio string `json:"path_audio"`
	Length     int64  `json:"length"`
}

type Video_previews struct {
	database.Model
	Video_id   int64  `json:"video_id" gorm:"column:video_id"`
	Path       string `json:"path"`
	Length     int64  `json:"length"`
}
