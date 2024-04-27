package models

import (
	"github.com/Rhaqim/creds/internal/database"
	"gorm.io/gorm"
)

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

// Insert creates a new CredsFields.
func (O *CredsFields) Insert() error {
	return database.DB.Model(O).Create(O).Error
}

// GetnByID retrieves an CredsFields by its ID.
func (O *CredsFields) GetByID(id int) error {
	return database.DB.Where("id = ?", id).First(O).Error
}

// GetMultipleByUserID retrieves CredsFieldss by user ID.
func (O CredsFields) GetMultipleByCredsID(credsID int) ([]CredsFields, error) {
	var orgs []CredsFields
	err := database.DB.Where("creds_id = ?", credsID).Find(&orgs).Error
	return orgs, err
}

// Update updates an CredsFields.
func (O *CredsFields) Update() error {
	return database.DB.Save(O).Error
}

// Delete deletes an CredsFields.
func (O *CredsFields) Delete() error {
	return database.DB.Delete(O).Error
}

// Insert creates a new Creds.
func (O *Creds) Insert() error {
	return database.DB.Model(O).Create(O).Error
}

// GetnByID retrieves an Creds by its ID.
func (O *Creds) GetByID(id int) error {
	return database.DB.Where("id = ?", id).First(O).Error
}

// GetMultipleByOrgrID retrieves Credss by user ID.
func (O Creds) GetMultipleByOrgID(orgID int) ([]Creds, error) {
	var orgs []Creds
	err := database.DB.Where("organization_id = ?", orgID).Find(&orgs).Error
	return orgs, err
}

// GetMultipleByEnvID retrieves Credss by user ID.
func (O Creds) GetMultipleByEnvID(envID int) ([]Creds, error) {
	var orgs []Creds
	err := database.DB.Where("environment_id = ?", envID).Find(&orgs).Error
	return orgs, err
}

// Update updates an Creds.
func (O *Creds) Update() error {
	return database.DB.Save(O).Error
}

// Delete deletes an Creds.
func (O *Creds) Delete() error {
	return database.DB.Delete(O).Error
}
