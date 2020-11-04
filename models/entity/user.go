package entity

import "lionnix-metrics-api/pkg/database"

//User - model user add groups and references
type User struct {
	database.Model
	UserName    string `json:"username" gorm:"column:username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	FullName    string `json:"full_name"`
	Facebook    string `json:"facebook"`
	Role        string `json:"role"`
	PackageID   int64  `json:"package_id"`
	PackageVnID   int64  `json:"package_vn_id"`
	RefID       int64  `json:"ref_id"`
	RefCode     string `json:"ref_code"`
	CanRefer    bool   `json:"can_refer"`
	Status      string `json:"status"`
	Point       int64  `json:"point"`
}
