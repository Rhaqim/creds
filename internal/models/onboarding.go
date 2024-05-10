package models

import (
	"github.com/Rhaqim/creds/internal/database"
	"gorm.io/gorm"
)

type OnBoarding struct {
	gorm.Model
	UserID    uint     `json:"user_id" form:"user_id" query:"user_id" gorm:"not null" binding:"required"`
	Resource  Resource `json:"resource" form:"resource" query:"resource" binding:"required" gorm:"foreignKey:ID"`
	Completed bool     `json:"completed" form:"completed" query:"completed" gorm:"not null" binding:"required" default:"false"`
}

func (O *OnBoarding) Insert() error {
	return database.Insert(O)
}

func (O *OnBoarding) GetOnBoarding() *gorm.DB {
	return database.DB.Model(O).Preload("Resource")
}

func (O *OnBoarding) GetByID(id int) error {
	return O.GetOnBoarding().Where("id = ?", id).First(O).Error
}

func (O OnBoarding) GetMultipleByCredentialID(credentialID int) ([]OnBoarding, error) {
	var orgs []OnBoarding
	err := O.GetOnBoarding().Where("credential_id = ?", credentialID).Find(&orgs).Error
	return orgs, err
}

func (O *OnBoarding) Update() error {
	_, err := database.Update(O, "id = ?", O.ID)
	if err != nil {
		return err
	}

	return nil
}

func (O *OnBoarding) Delete() error {
	_, err := database.Delete(O, "id = ?", O.ID)
	if err != nil {
		return err
	}

	return nil

}
