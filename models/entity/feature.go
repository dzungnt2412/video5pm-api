package entity

import "lionnix-metrics-api/pkg/database"

//Feature - model feature
type Feature struct {
	database.Model
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Position    int    `json:"order" gotm:"position"`
	System      string `json:"system" gorm:"column:system"`
	Reference   string `json:"reference" gorm:"column:reference"`
	Status      int    `json:"status" gorm:"column:status"`
	Icon        string `json:"icon" gorm:"column:icon"`
}

//GroupFeature - relation table group_features
type GroupFeature struct {
	database.Model
	GroupID   int64 `json:"group_id"`
	FeatureID int64 `json:"feature_id"`
}
