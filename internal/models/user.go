package models

import (
	"github.com/Rhaqim/creds/internal/database"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	DisplayName  string `json:"display_name" binding:"required"`
	Email        string `json:"email" gorm:"unique" binding:"required"`
	OAuthID      string `json:"oauth_id" gorm:"unique" binding:"required"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// Insert creates a new organization.
func (O *User) Insert() error {
	return database.DB.Model(O).Create(O).Error
}

// GetnByID retrieves an User by its ID.
func (O *User) GetByID(id int) error {
	return database.DB.Where("id = ?", id).First(O).Error
}

// GetnByEmail retrieves an User by its ID.
func (O *User) GetByEmail(email string) error {
	return database.DB.Where("email = ?", email).First(O).Error
}

// GetByOAuthID retrieves an User by its OAuth ID.
func (O *User) GetByOAuthID(id string) error {
	return database.DB.Where("oauth_id = ?", id).First(O).Error
}

// GetMultipleByUserID retrieves Users by user ID.
func (O User) GetAll() ([]string, error) {
	var orgs []User

	var dis []string
	err := database.DB.Where("1 = 1").Find(&orgs).Error

	for _, v := range orgs {
		dis = append(dis, v.DisplayName)
	}

	return dis, err
}

// GetMemberUsers retrieves users for the user IDs
func (O User) GetMemberUsers(userIDs []uint) ([]User, error) {
	var orgs []User
	err := database.DB.Select("id, display_name").Where("id IN ?", userIDs).Find(&orgs).Error
	return orgs, err
}

// Update updates an User.
func (O *User) Update() error {
	return database.DB.Save(O).Error
}

// Delete deletes an User.
func (O *User) Delete() error {
	return database.DB.Delete(O).Error
}
