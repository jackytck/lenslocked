package models

import (
	"errors"

	"github.com/jinzhu/gorm"

	// dialects: postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is returned when a resource cannot be found
	// in the database.
	ErrNotFound = errors.New("models: resource not found")
)

// UserService provides services for interacting with user model.
type UserService struct {
	db *gorm.DB
}

// NewUserService helps create a UserService with db info.
func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	us := UserService{
		db: db,
	}
	return &us, nil
}

// ByID looks up user by the id provided.
// 1. user, nil (user found)
// 2. nil, ErrNotFound (user not found)
// 3. nil, otherError (others)
func (us *UserService) ByID(id uint) (*User, error) {
	var u User
	err := us.db.Where("id = ?", id).First(&u).Error
	switch err {
	case nil:
		return &u, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Close closes the UserService database connection.
func (us *UserService) Close() {
	us.db.Close()
}

// DestructiveReset drops the user table and rebuilds it.
func (us *UserService) DestructiveReset() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
}

// User represent the User model.
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}
