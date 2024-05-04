package models

import "gorm.io/gorm"

type Resource struct {
	gorm.Model
	Name        string `json:"name" form:"name" query:"name" gorm:"not null" binding:"required"`
	Description string `json:"description" form:"description" query:"description"`
	URL         string `json:"url" form:"url" query:"url"`
	Category    string `json:"category" form:"category" query:"category" gorm:"not null" binding:"required"`
	Guide       string `json:"guide" form:"guide" query:"guide"`
}

type OnBoarding struct {
	gorm.Model
	UserID     uint `json:"user_id" form:"user_id" query:"user_id" gorm:"not null" binding:"required"`
	ResourceID uint `json:"resource_id" form:"resource_id" query:"resource_id" gorm:"not null" binding:"required"`
	Completed  bool `json:"completed" form:"completed" query:"completed" gorm:"not null" binding:"required" default:"false"`
}
