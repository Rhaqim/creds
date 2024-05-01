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
	CreatorID        uint                  `json:"creator_id,omitempty" form:"creator_id,omitempty" query:"creator_id" gorm:"not null"`
	OrganizationName string                `json:"organization_name" form:"organization_name" query:"organization_name" gorm:"not null" binding:"required"`
	Description      string                `json:"description" form:"description" query:"description" gorm:"not null" binding:"required"`
	OrganizationType CredsOrganizationType `json:"organization_type" form:"organization_type" query:"organization_type" gorm:"not null" oneof:"company personal" binding:"required"`
	Members          []OrganizationMember  `json:"members" form:"members" query:"members" gorm:"foreignKey:OrganizationID"`
}

// Insert creates a new organization.
func (O *Organization) Insert() error {
	return database.DB.Model(O).Create(O).Error
}

func (O *Organization) GetOrganization(filter string, args ...interface{}) *gorm.DB {
	return database.DB.Preload("Members").Where(filter, args...)
}

// GetnByID retrieves an organization by its ID.
func (O *Organization) GetByID(id uint) error {
	return O.GetOrganization("id = ?", id).First(O).Error
}

// GetByOrganizationName retrieves an organization by its name.
func (O *Organization) GetByOrganizationName(name string) error {
	return O.GetOrganization("organization_name = ?", name).First(O).Error
}

// GetMultipleByUserID retrieves organizations by user ID.
func (O Organization) GetMultipleByUserID(userID uint) ([]Organization, error) {
	var orgs []Organization
	err := database.DB.Joins("Members").Where("user_id = ?", userID).Find(&orgs).Error
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

func (O *Organization) IsMember(userID uint) bool {
	for _, member := range O.Members {
		if member.UserID == userID {
			return true
		}
	}

	return false
}

func (O *Organization) CreateOrganization(user User) error {
	var err error
	var member OrganizationMember

	O.CreatorID = user.ID

	err = O.Insert()
	if err != nil {
		return err
	}

	member.OrganizationID = O.ID
	member.UserID = user.ID
	member.Role = Admin

	err = member.Insert()
	if err != nil {
		return err
	}

	return err
}
