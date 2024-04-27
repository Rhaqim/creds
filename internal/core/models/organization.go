package core

import "gorm.io/gorm"

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
