package entity

import "video5pm-api/pkg/database"

//Group - model group
type Group struct {
	database.Model
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
}

//UserGroup - relation table user_groups
type UserGroup struct {
	database.Model
	UserID  int64 `json:"user_id"`
	GroupID int64 `json:"group_id"`
}
