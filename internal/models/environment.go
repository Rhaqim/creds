package models

import (
	"github.com/Rhaqim/creds/internal/database"
	"gorm.io/gorm"
)

type CredsEnvironment int

const (
	Development CredsEnvironment = iota
	Staging
	Preproduction
	Production
)

type DevelopmentEnvironment struct {
	gorm.Model
	OrgnaizationID uint             `json:"organization_id" form:"organization_id" query:"organization_id" gorm:"not null"`
	Environment    CredsEnvironment `json:"environment" form:"environment" query:"environment" gorm:"not null" binding:"oneof=0 1 2 3"`
}

// Insert creates a new CredsEnvironment.
func (O *CredsEnvironment) Insert() error {
	return database.DB.Model(O).Create(O).Error
}

// GetnByID retrieves an CredsEnvironment by its ID.
func (O *CredsEnvironment) GetByID(id int) error {
	return database.DB.Where("id = ?", id).First(O).Error
}

// GetMultipleByUserID retrieves CredsEnvironments by user ID.
func (O CredsEnvironment) GetMultipleByOrgID(orgID int) ([]CredsEnvironment, error) {
	var orgs []CredsEnvironment
	err := database.DB.Where("organization_id = ?", orgID).Find(&orgs).Error
	return orgs, err
}

// Update updates an CredsEnvironment.
func (O *CredsEnvironment) Update() error {
	return database.DB.Save(O).Error
}

// Delete deletes an CredsEnvironment.
func (O *CredsEnvironment) Delete() error {
	return database.DB.Delete(O).Error
}
