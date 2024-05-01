package models

import (
	"strconv"

	"github.com/Rhaqim/creds/internal/database"
	errs "github.com/Rhaqim/creds/internal/errors"
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
	Description      string                `json:"description,omitempty" form:"description,omitempty" query:"description,omitempty"`
	OrganizationType CredsOrganizationType `json:"organization_type" form:"organization_type" query:"organization_type" gorm:"not null" oneof:"company personal" binding:"required"`
	Members          []OrganizationMember  `json:"members" form:"members" query:"members" gorm:"foreignKey:OrganizationID"`
}

type OrganizationReturn struct {
	Organization
	MembersUser      []User       `json:"members_user"`
	MembersCount     int          `json:"members_count"`
	Credentials      []Credential `json:"credentials"`
	CredentialsCount int          `json:"credentials_count"`
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
	err := O.GetOrganization("creator_id = ?", userID).Find(&orgs).Error
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

func (O *Organization) FetchOrganizations(user User) ([]OrganizationReturn, error) {
	var err error
	var resp []OrganizationReturn
	var orgs []Organization
	var cred Credential
	var creds []Credential

	// Fetch organizations
	orgs, err = O.GetMultipleByUserID(user.ID)
	if err != nil {
		return resp, err
	}

	// Fetch credentials
	for _, org := range orgs {
		creds, err = cred.GetMultipleByOrgIDWithoutEncK(org.ID)
		if err != nil {
			return resp, err
		}
	}

	// Prepare response
	for _, org := range orgs {
		resp = append(resp, OrganizationReturn{
			Organization:     org,
			MembersCount:     len(org.Members),
			Credentials:      creds,
			CredentialsCount: len(creds),
		})
	}

	return resp, err
}

func (O *Organization) FetchMembers() ([]User, error) {
	var err error
	var user User
	var members []User

	var user_ids []uint

	for _, member := range O.Members {
		user_ids = append(user_ids, member.UserID)
	}

	// Fetch members
	members, err = user.GetMemberUsers(user_ids)
	if err != nil {
		return members, err
	}

	return members, err
}

func (O *Organization) FetchOrganization(user User, id string) (OrganizationReturn, error) {
	var err error
	var cred Credential
	var creds []Credential
	var resp OrganizationReturn

	// convert id to uint
	orgId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return resp, err
	}

	// Fetch FetchOrganization
	err = O.GetByID(uint(orgId))
	if err != nil {
		return resp, err
	}

	if !O.IsMember(user.ID) {
		return resp, errs.ErrNotMemberOfOrganization
	}

	// Fetch members
	membersUser, err := O.FetchMembers()
	if err != nil {
		return resp, err
	}

	// Fetch credentials
	creds, err = cred.GetMultipleByOrgIDWithoutEncK(O.ID)
	if err != nil {
		return resp, err
	}

	// Prepare response
	resp = OrganizationReturn{
		Organization:     *O,
		MembersCount:     len(O.Members),
		MembersUser:      membersUser,
		Credentials:      creds,
		CredentialsCount: len(creds),
	}

	return resp, err
}
