package models

import (
	"github.com/Rhaqim/creds/internal/database"
	"gorm.io/gorm"
)

type OrganaizationMemberRole string

const (
	Admin  OrganaizationMemberRole = "admin"
	Member OrganaizationMemberRole = "member"
)

type OrganizationMember struct {
	gorm.Model
	OrganizationID uint                    `json:"organization_id" form:"organization_id" query:"organization_id" gorm:"not null" binding:"required"`
	UserID         uint                    `json:"user_id" form:"user_id" query:"user_id" gorm:"not null" binding:"required"`
	Role           OrganaizationMemberRole `json:"role" form:"role" query:"role" gorm:"not null" binding:"oneof=admin member"`
}

// Insert creates a new organization.
func (O *OrganizationMember) Insert() error {
	return database.DB.Model(O).Create(O).Error
}

// GetnByID retrieves an organization by its ID.
func (O *OrganizationMember) GetByID(id int) error {
	return database.DB.Where("id = ?", id).First(O).Error
}

// GetMultipleByUserID retrieves organizations by user ID.
func (O Organization) GetMultipleByOrgrID(orgID int) ([]Organization, error) {
	var orgs []Organization
	err := database.DB.Where("organization_id = ?", orgID).Find(&orgs).Error
	return orgs, err
}

// Update updates an organization.
func (O *OrganizationMember) Update() error {
	return database.DB.Save(O).Error
}

// Delete deletes an organization.
func (O *OrganizationMember) Delete() error {
	return database.DB.Delete(O).Error
}
