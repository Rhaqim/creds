package models

import (
	"errors"

	"github.com/Rhaqim/creds/internal/database"
	"gorm.io/gorm"
)

type Privileges string

const (
	ManagePrivilege Privileges = "manage"
	ReadPrivilege   Privileges = "read"
)

type OrganizationTeam struct {
	gorm.Model
	Name           string     `json:"name" form:"name" query:"name" gorm:"not null" binding:"required"`
	ManagerID      uint       `json:"manager_id" form:"manager_id" query:"manager_id"`
	CreatedByID    uint       `json:"created_by_id" form:"created_by_id" query:"created_by_id" gorm:"not null"`
	OrganizationID uint       `json:"organization_id" form:"organization_id" query:"organization_id" gorm:"not null" binding:"required"`
	Privileges     Privileges `json:"privileges" form:"privileges" query:"privileges" gorm:"not null" oneof:"manage read"`
}

// Insert creates a new organization.
func (O *OrganizationTeam) Insert() error {
	return database.DB.Model(O).Create(O).Error
}

// GetnByID retrieves an organization by its ID.
func (O *OrganizationTeam) GetByID(id uint) error {
	return database.DB.Where("id = ?", id).First(O).Error
}

func (O *OrganizationTeam) GetByUserID(id uint) error {
	return database.DB.Where("user_id = ?", id).First(O).Error
}

// GetMultipleByUserID retrieves organizations by user ID.
func (O OrganizationTeam) GetMultipleByOrgID(orgID uint) ([]OrganizationTeam, error) {
	var orgs []OrganizationTeam
	err := database.DB.Where("organization_id = ?", orgID).Find(&orgs).Error
	return orgs, err
}
func (O OrganizationTeam) GetMultipleByOrgIDName(orgID uint, name string) error {
	err := database.DB.Where("organization_id = ? AND name = ?", orgID, name).First(&O).Error
	return err
}

// Update updates an organization.
func (O *OrganizationTeam) Update() error {
	return database.DB.Save(O).Error
}

// Delete deletes an organization.
func (O *OrganizationTeam) Delete() error {
	return database.DB.Delete(O).Error
}

func (O *OrganizationTeam) CreateTeam(user User) error {
	var err error

	// Check if the organization exists
	var org Organization
	if err := org.GetByID(uint(O.OrganizationID)); err != nil {
		return err
	}

	// Check if the user is an admin
	var member OrganizationMember
	if err = member.GetByUserID(user.ID); err != nil {
		return err
	}

	if member.Role != Admin {
		return errors.New("User is not an admin")
	}

	// Check if the team exists
	if err := O.GetMultipleByOrgIDName(org.ID, O.Name); err == nil {
		return errors.New("team already exists")
	}

	// Create the team
	O.CreatedByID = user.ID
	O.Privileges = ReadPrivilege

	if err := O.Insert(); err != nil {
		return err
	}

	return nil
}

type OrganaizationMemberRole string

const (
	Admin     OrganaizationMemberRole = "admin"
	Manager   OrganaizationMemberRole = "manager"
	Member    OrganaizationMemberRole = "member"
	SeniorDev OrganaizationMemberRole = "senior_dev"
	JuniorDev OrganaizationMemberRole = "junior_dev"
	QualityAs OrganaizationMemberRole = "quality_as"
)

type OrganaizationMemberStatus string

const (
	Pending  OrganaizationMemberStatus = "pending"
	Active   OrganaizationMemberStatus = "active"
	Inactive OrganaizationMemberStatus = "inactive"
)

type OrganizationMember struct {
	gorm.Model
	OrganizationID uint                      `json:"organization_id" form:"organization_id" query:"organization_id" gorm:"not null" binding:"required"`
	UserID         uint                      `json:"user_id" form:"user_id" query:"user_id" gorm:"not null" binding:"required"`
	Role           OrganaizationMemberRole   `json:"role" form:"role" query:"role" gorm:"not null" binding:"oneof=admin member"`
	TeamID         uint                      `json:"team_id" form:"team_id" query:"team_id"`
	Status         OrganaizationMemberStatus `json:"status" form:"status" query:"status" binding:"oneof=pending active inactive" default:"pending"`
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

func (O *OrganizationMember) InviteMember(user User, email string) error {
	var err error

	// Check if the organization exists
	var org Organization
	if err := org.GetByID(O.OrganizationID); err != nil {
		return err
	}

	// Check if the user is an admin
	if err = O.GetByUserID(user.ID); err != nil {
		return err
	}

	if O.Role != Admin {
		return errors.New("User is not an admin")
	}

	// check if the team exists
	var team OrganizationTeam
	if err := team.GetByID(O.TeamID); err != nil {
		return err
	}

	// Check if the user exists
	var usr User
	if err := usr.GetByEmail(email); err != nil {
		// if user doesn't exist send an email to join
		return nil
	}

	// Check if the user is already a member
	var member OrganizationMember
	if err := database.DB.Where("organization_id = ? AND user_id = ?", org.ID, usr.ID).First(&member).Error; err == nil {
		return errors.New("user is already a member")
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
