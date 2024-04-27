package core

import "gorm.io/gorm"

type Creds struct {
	gorm.Model
	OrganizationID int
	EnvironmentID  int
	Format         CredsSaveFormat
}

type CredsFields struct {
	gorm.Model
	CredsID int
	Key     string
	Value   string
}
