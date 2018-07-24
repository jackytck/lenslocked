package models

import (
	"errors"

	"github.com/jackytck/lenslocked/hash"
	"github.com/jackytck/lenslocked/rand"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	// dialects: postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is returned when a resource cannot be found
	// in the database.
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID is provided
	// to a method like Delete.
	ErrInvalidID = errors.New("models: ID provided was invalid")

	// ErrInvalidPassword is returned when an invalid password
	// is used when attempting to authenticate a user.
	ErrInvalidPassword = errors.New("models: incorrect password provided")
)

// UserService provides services for interacting with user model.
type UserService struct {
	db   *gorm.DB
	hmac hash.HMAC
}

const userPwPepper = "P4P]tV6$LZc;,bu5"
const hmacSecretKey = "E4j!STJ$??cc]UhQ"

// NewUserService helps create a UserService with db info.
func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	hmac := hash.NewHMAC(hmacSecretKey)
	us := UserService{
		db:   db,
		hmac: hmac,
	}
	return &us, nil
}

// ByID looks up user with the id provided.
// 1. user, nil (user found)
// 2. nil, ErrNotFound (user not found)
// 3. nil, otherError (others)
func (us *UserService) ByID(id uint) (*User, error) {
	var u User
	db := us.db.Where("id = ?", id)
	err := first(db, &u)
	return &u, err
}

// ByEmail looks up user with the email provided.
// 1. user, nil (user found)
// 2. nil, ErrNotFound (user not found)
// 3. nil, otherError (others)
func (us *UserService) ByEmail(email string) (*User, error) {
	var u User
	db := us.db.Where("email = ?", email)
	err := first(db, &u)
	return &u, err
}

// ByRemember looks up a user with the given remember token
// and returns that user. This method will handle hashing
// the token for us.
// Errors are the same as ByEmail
func (us *UserService) ByRemember(token string) (*User, error) {
	var u User
	rememberHash := us.hmac.Hash(token)
	db := us.db.Where("remember_hash = ?", rememberHash)
	err := first(db, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Authenticate can be used to authenticate a user with the
// provided email address and password.
// If the email address provided is invalid, this will return
// 	nil, ErrNotFound
// If the password provided is invalid, this will return
// 	nil, ErrInvalidPassword
// If the email and password are both valid, this will return
// 	user, nil
// Otherwise if another error is encountered this will return
// 	nil, error
func (us *UserService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}
	return foundUser, nil
}

// first will query using the provided gorm.DB and it will
// get the first item returned and place it into dst. If
// nothing is found in the query, it will return ErrNotFound.
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// Create creates the provided user and backfill data
// like the ID, CreatedAt, and UpdatedAt fields.
func (us *UserService) Create(user *User) error {
	pwBytpes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytpes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = us.hmac.Hash(user.Remember)
	return us.db.Create(user).Error
}

// Update updates the provided user with all of the data
// in the provided user object.
func (us *UserService) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = us.hmac.Hash(user.Remember)
	}
	return us.db.Save(user).Error
}

// Delete deletes the user with the provided ID.
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	u := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&u).Error
}

// Close closes the UserService database connection.
func (us *UserService) Close() error {
	return us.db.Close()
}

// DestructiveReset drops the user table and rebuilds it.
func (us *UserService) DestructiveReset() error {
	if err := us.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return us.AutoMigrate()
}

// AutoMigrate attemps to automatically migrate the user table.
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// User represent the User model.
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}
