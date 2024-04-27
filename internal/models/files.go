package models

import (
	"github.com/Rhaqim/creds/internal/database"
	"gorm.io/gorm"
)

type CredentialFile struct {
	gorm.Model
	CredentialID uint                 `json:"credential_id" form:"credential_id" query:"credential_id" gorm:"not null"`
	FileName     string               `json:"file_name" form:"file_name" query:"file_name" gorm:"not null"`
	FileSize     int64                `json:"file_size" form:"file_size" query:"file_size" gorm:"not null"`
	FileData     []byte               `json:"file_data" form:"file_data" query:"file_data" gorm:"not null"`
	FileFormat   CredentialSaveFormat `json:"file_format" form:"file_format" query:"file_format" gorm:"not null"`
}

func (O *CredentialFile) Insert() error {
	return database.DB.Model(O).Create(O).Error
}

func (O *CredentialFile) GetByID(id int) error {
	return database.DB.Where("id = ?", id).First(O).Error
}

func (O CredentialFile) GetMultipleByCredentialID(credentialID int) ([]CredentialFile, error) {
	var orgs []CredentialFile
	err := database.DB.Where("credential_id = ?", credentialID).Find(&orgs).Error
	return orgs, err
}

func (O *CredentialFile) Update() error {
	return database.DB.Save(O).Error
}

func (O *CredentialFile) Delete() error {
	return database.DB.Delete(O).Error
}

func (O *CredentialFile) GetByFileName(fileName string) error {
	return database.DB.Where("file_name = ?", fileName).First(O).Error
}

func (O *CredentialFile) GetByCredentialIDAndFileName(credentialID int, fileName string) error {
	return database.DB.Where("credential_id = ? AND file_name = ?", credentialID, fileName).First(O).Error
}
