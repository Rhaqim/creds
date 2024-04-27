package app

import "gorm.io/gorm"

type Creds struct {
	gorm.Model
	OrganizationID uint            `json:"organization_id" form:"organization_id" query:"organization_id" gorm:"not null"`
	EnvironmentID  uint            `json:"environment_id" form:"environment_id" query:"environment_id" gorm:"not null"`
	FileFormat     CredsSaveFormat `json:"file_format" form:"file_format" query:"file_format" gorm:"not null"`
}

type CredsFields struct {
	gorm.Model
	CredsID uint   `json:"creds_id" form:"creds_id" query:"creds_id" gorm:"not null"`
	Key     string `json:"key" form:"key" query:"key" gorm:"not null"`
	Value   string `json:"value" form:"value" query:"value" gorm:"not null"`
}
