package models

import (
	"strconv"

	"github.com/Rhaqim/creds/internal/database"
	err "github.com/Rhaqim/creds/internal/errors"
	"github.com/Rhaqim/creds/internal/lib"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CredsEnvironment int

const (
	Development CredsEnvironment = iota + 1
	Staging
	Preproduction
	Production
)

type Credential struct {
	gorm.Model
	Name           string           `json:"name" form:"name" query:"name" gorm:"not null" binding:"required"`
	OrganizationID uint             `json:"organization_id" form:"organization_id" query:"organization_id" gorm:"not null" binding:"required"`
	EncryptionKey  []byte           `json:"encryption_key,omitempty" form:"encryption_key,omitempty" query:"encryption_key" gorm:"not null"`
	Environment    CredsEnvironment `json:"environment" form:"environment" query:"environment" gorm:"not null" oneof:"1 2 3 4" binding:"required"`
	Version        string           `json:"version" form:"version" query:"version" gorm:"not null"`
}

type CredentialField struct {
	gorm.Model
	CredentialID     uint   `json:"creds_id" form:"creds_id" query:"creds_id" gorm:"not null"`
	CredentialFileID uint   `json:"file_id" form:"file_id" query:"file_id"`
	Key              string `json:"key" form:"key" query:"key" gorm:"not null"`
	Value            string `json:"value" form:"value" query:"value" gorm:"not null"`
}

type CredentialReturn struct {
	Credential Credential `json:"credential"`
	Fields     []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	File CredentialFile `json:"file,omitempty"`
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
func (O CredentialField) GetMultipleByCredentialID(credsID uint) ([]CredentialField, error) {
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

func (O Credential) GetMultipleByOrgIDWithoutEncK(orgID uint) ([]Credential, error) {
	var orgs []Credential
	err := database.DB.Select("id, name, organization_id, environment, version").Where("organization_id = ?", orgID).Find(&orgs).Error
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

func (O *Credential) IsMember(user User) error {
	var org Organization

	// validate organization
	if err := org.GetByID(O.OrganizationID); err != nil {
		return err
	}

	// validate member
	if !org.IsMember(user.ID) {
		return err.ErrNotMemberOfOrganization
	}

	return nil
}

func (O *Credential) CreateCredential(user User) error {

	var encryptor lib.EncryptionService

	// is user a member of the organization
	if err := O.IsMember(user); err != nil {
		return err
	}

	// Generate encryption key
	O.EncryptionKey = encryptor.GenerateEncryptionKey()
	O.Version = uuid.New().String()

	// Insert credential
	if err := O.Insert(); err != nil {
		return err
	}

	return nil
}

func (O *Credential) FetchFile(user User) CredentialFile {

	var file CredentialFile

	// Get file
	_ = file.GetByCredentialID(O.ID)

	return file
}

func (O *Credential) FetchFields() []CredentialField {
	var fields []CredentialField
	var field CredentialField

	// Get fields
	fields, _ = field.GetMultipleByCredentialID(O.ID)

	return fields
}

func (O *Credential) FetchCredential(user User, id string) (CredentialReturn, error) {
	var err error
	var resp CredentialReturn

	// convert id to uint
	credId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return resp, err
	}

	// Get credential
	if err = O.GetByID(uint(credId)); err != nil {
		return resp, err
	}

	// is user a member of the organization
	if err = O.IsMember(user); err != nil {
		return resp, err
	}

	// Get file
	file := O.FetchFile(user)

	// Get fields
	fields := O.FetchFields()

	// append fields to response
	if len(fields) == 0 {
		resp.Fields = make([]struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}, 0)
	} else {

		for _, f := range fields {
			resp.Fields = append(resp.Fields, struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			}{
				Key:   f.Key,
				Value: f.Value,
			})
		}
	}

	O.EncryptionKey = nil

	resp.Credential = *O

	if file.ID != 0 {
		resp.File = file
	} else {
		resp.File = CredentialFile{}
	}

	return resp, err
}
