package services

import (
	"video5pm-api/models/entity"

	"github.com/jinzhu/gorm"
)

func NewPackageService(db *gorm.DB) *PackageService {
	return &PackageService{db: db}
}

type PackageService struct {
	db *gorm.DB
}

//UpdateUserPackage - service update package id of user
func (c *PackageService) UpdateUserPackage(uid, pid int64) error {
	var user entity.User
	c.db.First(&user, uid)
	err := c.db.Save(&user).Error

	return err

}

func (c *PackageService) UpdateUserVnPackage(uid, pid int64) error {
	var user entity.User
	c.db.First(&user, uid)
	err := c.db.Save(&user).Error

	return err

}
