package core

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
	OrgnaizationID int
	Environment    CredsEnvironment
}
