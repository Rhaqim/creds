package models

import (
	"errors"
	"strconv"

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
func (O *OrganizationMember) GetByID(id uint) error {
	return database.DB.Where("id = ?", id).First(O).Error
}

func (O *OrganizationMember) GetByUserID(id uint) error {
	return database.DB.Where("user_id = ?", id).First(O).Error
}

// GetMultipleByUserID retrieves organizations by user ID.
func (O OrganizationMember) GetMultipleByOrgID(orgID uint) ([]OrganizationMember, error) {
	var orgs []OrganizationMember
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

func (O *OrganizationMember) InviteMember(user User, email, id string) error {
	// convert id to uint
	orgId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}

	// Check if the organization exists
	var org Organization
	if err := org.GetByID(uint(orgId)); err != nil {
		return err
	}

	// Check if the user is an admin
	if err = O.GetByUserID(user.ID); err != nil {
		return err
	}

	if O.Role != Admin {
		return errors.New("User is not an admin")
	}

	// Check if the user exists
	var usr User
	if err := usr.GetByEmail(email); err != nil {
		// if user doesn't exist send an email to join
		return nil
	}

	// Check if the user is already a member
	var member OrganizationMember
	if err := database.DB.Where("organization_id = ? AND user_id = ?", id, usr.ID).First(&member).Error; err == nil {
		return errors.New("User is already a member")
	}

	// Create the member
	member = OrganizationMember{
		OrganizationID: org.ID,
		UserID:         usr.ID,
		Role:           Member,
	}

	if err := member.Insert(); err != nil {
		return err
	}

	return nil
}
