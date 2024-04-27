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

// OrganizationService is a service that provides operations on organizations.
type OrganizationService struct {
	db *gorm.DB
}

// NewOrganizationService creates a new organization service.
func NewOrganizationService() *OrganizationService {
	return &OrganizationService{db: database.DB}
}

// CreateOrganization creates a new organization.
func (s *OrganizationService) CreateOrganization(org *Organization) error {
	return s.db.Create(org).Error
}

// GetOrganizationByID retrieves an organization by its ID.
func (s *OrganizationService) GetOrganizationByID(id int) (*Organization, error) {
	org := new(Organization)
	err := s.db.First(org, id).Error
	return org, err
}

// GetOrganizationsByUserID retrieves organizations by user ID.
func (s *OrganizationService) GetOrganizationsByUserID(userID int) ([]*Organization, error) {
	var orgs []*Organization
	err := s.db.Where("user_id = ?", userID).Find(&orgs).Error
	return orgs, err
}

// UpdateOrganization updates an organization.
func (s *OrganizationService) UpdateOrganization(org *Organization) error {
	return s.db.Save(org).Error
}

// DeleteOrganization deletes an organization.
func (s *OrganizationService) DeleteOrganization(org *Organization) error {
	return s.db.Delete(org).Error
}
