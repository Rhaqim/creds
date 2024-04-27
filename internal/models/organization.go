package models

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
	OwnerID          uint                  `json:"owner_id" form:"owner_id" query:"owner_id" gorm:"not null" binding:"required"`
	OrganizationName string                `json:"organization_name" form:"organization_name" query:"organization_name" gorm:"not null" binding:"required"`
	OrganizationType CredsOrganizationType `json:"organization_type" form:"organization_type" query:"organization_type" gorm:"not null" binding:"oneof=company personal"`
}

// Insert creates a new organization.
func (O *Organization) Insert() error {
	return database.DB.Model(O).Create(O).Error
}

// GetnByID retrieves an organization by its ID.
func (O *Organization) GetByID(id int) error {
	return database.DB.Where("id = ?", id).First(O).Error
}

// GetByOrganizationName retrieves an organization by its name.
func (O *Organization) GetByOrganizationName(name string) error {
	return database.DB.Where("organization_name = ?", name).First(O).Error
}

// GetMultipleByUserID retrieves organizations by user ID.
func (O Organization) GetMultipleByUserID(userID int) ([]Organization, error) {
	var orgs []Organization
	err := database.DB.Where("owner_id = ?", userID).Find(&orgs).Error
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

func (O *Organization) CreateOrganization(user User) error {
	O.OwnerID = user.ID

	return O.Insert()
}
