package core

import "gorm.io/gorm"

type Organization struct {
	gorm.Model
	ID   int
	Name string
}
