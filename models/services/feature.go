package services

import (
	"github.com/jinzhu/gorm"
)

type FeatureResult struct {
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	Reference string `json:"reference"`
}

func NewFeatureService(db *gorm.DB) *FeatureService {
	return &FeatureService{db: db}
}

type FeatureService struct {
	db *gorm.DB
}

//GetFeaturesByGroupID - get all features in 1 group
func (c *FeatureService) GetFeaturesByGroupID(gid int64) ([]FeatureResult, error) {
	var features []FeatureResult
	err := c.db.Where("group_id = ?", gid).Find(&features).Error

	if err != nil {
		return nil, err
	}

	return features, nil
}

//GetFeaturesByGroupID - get all features of 1 user
func (c *FeatureService) GetFeaturesByUserID(uid int64) []FeatureResult {
	var features []FeatureResult
	c.db.Table("features").Select("name ,icon, reference").Where("id in (select feature_id from group_features where status = 1 and group_id in (select group_id from user_groups where user_id = ?))", uid).Order("position asc").Find(&features)

	return features
}

//GetListReferences - extract list references
func (c *FeatureService) GetListReferences(features []FeatureResult) []string {
	var result []string
	for _, f := range features {
		result = append(result, f.Reference)
	}
	return result
}
