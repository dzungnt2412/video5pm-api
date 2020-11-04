package services

import (
	"lionnix-metrics-api/models/entity"
	"lionnix-metrics-api/pkg/logger"

	"github.com/jinzhu/gorm"
)

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

type AuthService struct {
	db *gorm.DB
}

//GetUserByUserName - query user by username
func (c *AuthService) GetUserByUserName(username string) (*entity.User, error) {
	var u entity.User
	err := c.db.Where("username = ?", username).Or("email = ?", username).First(&u).Error

	if err != nil {
		logger.Log.Errorf("%v", err)
		return nil, err
	}

	return &u, nil
}
