package migrations

import (
	"github.com/Rhaqim/creds/internal/database"
	"github.com/Rhaqim/creds/internal/models"
)

func Migrate() {
	// Migrate the database
	database.DB.AutoMigrate(&models.User{}, &models.Organization{},
		&models.OrganizationMember{}, &models.Credential{},
		&models.CredentialFile{}, &models.CredentialField{})
}
