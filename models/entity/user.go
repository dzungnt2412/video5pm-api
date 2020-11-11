package entity

import "video5pm-api/pkg/database"

//User - model user add groups and references
type User struct {
	database.Model
	UserName string `json:"username" gorm:"column:username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
