package models

import (
	"github.com/Rhaqim/creds/internal/database"
	"github.com/Rhaqim/creds/internal/lib"
	"gorm.io/gorm"
)

type CredsEnvironment int

const (
	Development CredsEnvironment = iota
	Staging
	Preproduction
	Production
)

type Credential struct {
	gorm.Model
	OrganizationID uint             `json:"organization_id" form:"organization_id" query:"organization_id" gorm:"not null"`
	EncryptionKey  []byte           `json:"encryption_key" form:"encryption_key" query:"encryption_key" gorm:"not null"`
	Environment    CredsEnvironment `json:"environment" form:"environment" query:"environment" gorm:"not null" binding:"oneof=0 1 2 3"`
	Version        string           `json:"version" form:"version" query:"version" gorm:"not null"`
}

type CredentialField struct {
	gorm.Model
	CredentialID uint   `json:"creds_id" form:"creds_id" query:"creds_id" gorm:"not null"`
	Key          string `json:"key" form:"key" query:"key" gorm:"not null"`
	Value        string `json:"value" form:"value" query:"value" gorm:"not null"`
}

// Insert creates a new CredentialField.
func (O *CredentialField) Insert() error {
	return database.DB.Model(O).Create(O).Error
}

// GetnByID retrieves an CredentialField by its ID.
func (O *CredentialField) GetByID(id int) error {
	return database.DB.Where("id = ?", id).First(O).Error
}

// GetMultipleByUserID retrieves CredentialFields by user ID.
func (O CredentialField) GetMultipleByCredentialID(credsID int) ([]CredentialField, error) {
	var orgs []CredentialField
	err := database.DB.Where("creds_id = ?", credsID).Find(&orgs).Error
	return orgs, err
}

// Update updates an CredentialField.
func (O *CredentialField) Update() error {
	return database.DB.Save(O).Error
}

// Delete deletes an CredentialField.
func (O *CredentialField) Delete() error {
	return database.DB.Delete(O).Error
}

// Insert creates a new Credential.
func (O *Credential) Insert() error {
	return database.DB.Model(O).Create(O).Error
}

// GetnByID retrieves an Credential by its ID.
func (O *Credential) GetByID(id uint) error {
	return database.DB.Where("id = ?", id).First(O).Error
}

// GetMultipleByOrgrID retrieves Credentials by user ID.
func (O Credential) GetMultipleByOrgID(orgID uint) ([]Credential, error) {
	var orgs []Credential
	err := database.DB.Where("organization_id = ?", orgID).Find(&orgs).Error
	return orgs, err
}

// GetMultipleByEnvID retrieves Credentials by user ID.
func (O Credential) GetMultipleByEnvID(envID uint) ([]Credential, error) {
	var orgs []Credential
	err := database.DB.Where("environment_id = ?", envID).Find(&orgs).Error
	return orgs, err
}

// Update updates an Credential.
func (O *Credential) Update() error {
	return database.DB.Save(O).Error
}

// Delete deletes an Credential.
func (O *Credential) Delete() error {
	return database.DB.Delete(O).Error
}

func (O *Credential) Create() error {
	var encryptor lib.EncryptionService

	// Generate encryption key
	O.EncryptionKey = encryptor.GenerateEncryptionKey()

	// Insert credential
	if err := O.Insert(); err != nil {
		return err
	}

	return nil
}
