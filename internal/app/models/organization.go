package app

import (
	"github.com/Rhaqim/creds/internal/database"
	"gorm.io/gorm"
)

type CredsOrganizationType string

const (
	Company  CredsOrganizationType = "company"
	Personal CredsOrganizationType = "personal"
)

type Organization struct {
	gorm.Model
	UserID           int                   `json:"user_id" form:"user_id" query:"user_id" gorm:"not null"`
	OrganizationName string                `json:"organization_name" form:"organization_name" query:"organization_name" gorm:"not null"`
	OrganizationType CredsOrganizationType `json:"organization_type" form:"organization_type" query:"organization_type" gorm:"not null"`
}

// Insert creates a new organization.
func (O *Organization) Insert() error {
	return database.DB.Model(O).Create(O).Error
}

// GetnByID retrieves an organization by its ID.
func (O *Organization) GetByID(id int) error {
	return database.DB.Where("id = ?", id).First(O).Error
}

// GetMultipleByUserID retrieves organizations by user ID.
func (O Organization) GetMultipleByUserID(userID int) ([]Organization, error) {
	var orgs []Organization
	err := database.DB.Where("user_id = ?", userID).Find(&orgs).Error
	return orgs, err
}

// Update updates an organization.
func (O *Organization) Update() error {
	return database.DB.Save(O).Error
}

// Delete deletes an organization.
func (O *Organization) Delete() error {
	return database.DB.Delete(O).Error
}
