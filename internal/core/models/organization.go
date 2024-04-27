package core

import "gorm.io/gorm"

type CredsOrganizationType string

const (
	Company  CredsOrganizationType = "company"
	Personal CredsOrganizationType = "personal"
)

type Organization struct {
	gorm.Model
	ID   int
	Name string
	Type CredsOrganizationType
}
