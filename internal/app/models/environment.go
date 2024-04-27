package app

import "gorm.io/gorm"

type CredsEnvironment int

const (
	Development CredsEnvironment = iota
	Staging
	Preproduction
	Production
)

type DevelopmentEnvironment struct {
	gorm.Model
	OrgnaizationID uint             `json:"organization_id" form:"organization_id" query:"organization_id" gorm:"not null"`
	Environment    CredsEnvironment `json:"environment" form:"environment" query:"environment" gorm:"not null" binding:"oneof=0 1 2 3"`
}
