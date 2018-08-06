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

const userPwPepper = "P4P]tV6$LZc;,bu5"
const hmacSecretKey = "E4j!STJ$??cc]UhQ"

// User represent the User model stored in our database.
// This is used for user accounts, storing both an email
// address and a password so users can log in and gain
// access to their content.
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

// UserService is a set of methods used to manipulate and
// work with the user model.
type UserService interface {
	// Authenticate verifies the provided email address and
	// password are correct. If they are correct, the user
	// corresponding to that email will be returned. Otherwise
	// you will receive either:
	// ErrNotFound, ErrInvalidPassword, or another error if
	// something goes wrong.
	Authenticate(email, password string) (*User, error)
	UserDB
}

// UserDB is used to interact with the users database.
//
// For all single user queries:
// a. If the user is found, return a nil error.
// b. If the user is not found, return ErrNotFound.
// c. If there is another error, return an error with
// more info about what went wrong. This may not be
// an error generated by the models package.
//
// For single user queries, any error but ErrNotFound should
// probably result in a 500 error.
type UserDB interface {
	// Methods for querying for single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Used to close a DB connection
	Close() error

	// Migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

// NewUserService helps create a UserService with db info.
func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := &userValidator{
		hmac:   hmac,
		UserDB: ug,
	}
	return &userService{
		UserDB: uv,
	}, nil
}

var _ UserService = &userService{}

// UserService provides services for interacting with user model.
type userService struct {
	UserDB
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
func (us *userService) Authenticate(email, password string) (*User, error) {
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

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

var _ UserDB = &userValidator{}

type userValidator struct {
	UserDB
	hmac hash.HMAC
}

// ByRemember hashes the remember token and then call
// ByRemember on the subsequent UserDB layer.
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFuncs(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

// Create creates the provided user with hashed password and remember hash.
func (uv *userValidator) Create(user *User) error {
	err := runUserValFuncs(user,
		uv.bcryptPassword,
		uv.setRememberIfUnset,
		uv.hmacRemember,
	)
	if err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

// Create hashes a remember token if it is provided.
func (uv *userValidator) Update(user *User) error {
	err := runUserValFuncs(user,
		uv.bcryptPassword,
		uv.hmacRemember,
	)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

// Delete deletes the user with the provided ID.
func (uv *userValidator) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	return uv.UserDB.Delete(id)
}

// bcryptPassword hashes user password with predefined pepper (userPwPepper)
// and bcrypt if the Password field is not the empty string
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytpes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytpes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

var _ UserDB = &userGorm{}

func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	ug := userGorm{
		db: db,
	}
	return &ug, nil
}

// userGorm provides services for interacting with user model.
type userGorm struct {
	db *gorm.DB
}

// ByID looks up user with the id provided.
// 1. user, nil (user found)
// 2. nil, ErrNotFound (user not found)
// 3. nil, otherError (others)
func (ug *userGorm) ByID(id uint) (*User, error) {
	var u User
	db := ug.db.Where("id = ?", id)
	err := first(db, &u)
	return &u, err
}

// ByEmail looks up user with the email provided.
// 1. user, nil (user found)
// 2. nil, ErrNotFound (user not found)
// 3. nil, otherError (others)
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var u User
	db := ug.db.Where("email = ?", email)
	err := first(db, &u)
	return &u, err
}

// ByRemember looks up a user with the given remember token
// and returns that user. This method expects the remember
// token to already be hashed.
// Errors are the same as ByEmail
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var u User
	db := ug.db.Where("remember_hash = ?", rememberHash)
	err := first(db, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Create creates the provided user and backfill data
// like the ID, CreatedAt, and UpdatedAt fields.
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

// Update updates the provided user with all of the data
// in the provided user object.
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

// Delete deletes the user with the provided ID.
func (ug *userGorm) Delete(id uint) error {
	u := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&u).Error
}

// Close closes the userService database connection.
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

// DestructiveReset drops the user table and rebuilds it.
func (ug *userGorm) DestructiveReset() error {
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return ug.AutoMigrate()
}

// AutoMigrate attemps to automatically migrate the user table.
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
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
