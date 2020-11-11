package services

import (
	"github.com/jinzhu/gorm"
	"video5pm-api/models/entity"
)

func NewGroupService(db *gorm.DB) *GroupService {
	return &GroupService{db: db}
}

type GroupService struct {
	db *gorm.DB
}

//GetGroupByUserID - get all groups of user
func (c *GroupService) GetGroupByUserID(uid int64) ([]entity.Group, error) {
	var groups []entity.Group
	err := c.db.Where("group_id = ?", uid).Find(&groups).Error

	if err != nil {
		return nil, err
	}

	return groups, nil
}

//GetGroupIDByUserID - get all group_ids of user
func (c *GroupService) GetGroupIDByUserID(uid int64) []int64 {
	var groups []int64
	c.db.Table("user_groups").Select("group_id").Where("user_id = ?", uid).Order("id asc").Find(&groups)

	return groups
}
