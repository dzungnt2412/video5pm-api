package services

import (
	"github.com/jinzhu/gorm"
	"lionnix-metrics-api/models/entity"
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
	user.PackageID = pid
	err := c.db.Save(&user).Error

	return err

}

func (c *PackageService) UpdateUserVnPackage(uid, pid int64) error {
	var user entity.User
	c.db.First(&user, uid)
	user.PackageVnID = pid
	err := c.db.Save(&user).Error

	return err

}
