package models

import (
	"github.com/jinzhu/gorm"

	// dialects: postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// User represent the User model
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}
