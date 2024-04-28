package models

import (
	"mime/multipart"
	"strconv"

	"github.com/Rhaqim/creds/internal/database"
	"github.com/Rhaqim/creds/internal/lib"
	"gorm.io/gorm"
)

// SaveFormat represents the format of the saved data.
type CredentialSaveFormat string

const (
	JSON  CredentialSaveFormat = "json"
	YAML  CredentialSaveFormat = "yaml"
	Plain CredentialSaveFormat = "plain"
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

func (O *CredentialFile) AppendDefaults(file *multipart.FileHeader, format string) error {
	filedata := make([]byte, file.Size)

	// Open the file
	src, err := file.Open()
	if err != nil {
		return err
	}

	// Read the file
	_, err = src.Read(filedata)
	if err != nil {
		return err
	}

	O.FileName = file.Filename
	O.FileSize = file.Size
	O.FileData = filedata

	// Set the file format
	switch format {
	case "json":
		O.FileFormat = JSON
	case "yaml":
		O.FileFormat = YAML
	case "plain":
		O.FileFormat = Plain
	default:
		O.FileFormat = Plain
	}

	return nil
}

func (O *CredentialFile) Save(credID uint) error {
	var encryptor lib.EncryptionService

	var parser lib.FileParser = lib.FileParser{
		FileFormat: string(O.FileFormat),
		FileData:   O.FileData,
	}

	// parse the file
	keyValues := parser.Parse()

	// Prepare a slice to hold CredentialField objects
	var credFields []CredentialField

	// Convert key-values to CredentialField objects
	for _, kv := range keyValues {
		encodedData, err := encryptor.Scramble(kv.Value.(string))
		if err != nil {
			return err
		}

		credField := CredentialField{
			CredentialID: credID,
			Key:          kv.Key,
			Value:        encodedData,
		}

		credFields = append(credFields, credField)
	}

	// Save the credential fields
	return database.DB.Create(&credFields).Error
}

func (O *CredentialFile) Process(user User, file *multipart.FileHeader, id string, format string) error {
	var cred Credential

	// Get the credential ID
	credID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}

	// Get the credential
	if err := cred.GetByID(uint(credID)); err != nil {
		return err
	}

	// Check if the user is a member of the organization
	if err := cred.IsMember(user); err != nil {
		return err
	}

	err = O.AppendDefaults(file, format)
	if err != nil {
		return err
	}

	// Save the file
	return O.Save(cred.ID)
}
