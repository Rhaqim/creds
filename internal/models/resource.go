package models

import (
	"github.com/Rhaqim/creds/internal/database"
	"gorm.io/gorm"
)

type Resource struct {
	gorm.Model
	OrganizationID uint   `json:"organization_id" form:"organization_id" query:"organizaiton_id"`
	Name           string `json:"name" form:"name" query:"name" gorm:"not null" binding:"required"`
	Description    string `json:"description" form:"description" query:"description"`
	URL            string `json:"url" form:"url" query:"url"`
	Category       string `json:"category" form:"category" query:"category" gorm:"not null" binding:"required"`
	Guide          string `json:"guide" form:"guide" query:"guide"`
}

func (O *Resource) Insert() error {
	return database.Insert(O)
}

func (O *Resource) GetResource() *gorm.DB {
	return database.DB.Model(O)
}

func (O *Resource) GetByID(id int) error {
	return O.GetResource().Where("id = ?", id).First(O).Error
}

func (O Resource) GetMultipleByOrgID(organizationID int) ([]Resource, error) {
	var orgs []Resource
	err := O.GetResource().Where("organization_id = ?", organizationID).Find(&orgs).Error
	return orgs, err
}

func (O *Resource) Update() error {
	_, err := database.Update(O, "id = ?", O.ID)
	if err != nil {
		return err
	}

	return nil
}

func (O *Resource) Delete() error {
	_, err := database.Delete(O, "id = ?", O.ID)
	if err != nil {
		return err
	}

	return nil

}
