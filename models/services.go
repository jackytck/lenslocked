package models

import "github.com/jinzhu/gorm"

// NewServices creates the Services.
func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	s := Services{
		User:    NewUserService(db),
		Gallery: NewGalleryService(db),
		db:      db,
	}
	return &s, nil
}

// Services represents all the services of the server.
type Services struct {
	Gallery GalleryService
	User    UserService
	db      *gorm.DB
}

// Close closes the database connection.
func (s *Services) Close() error {
	return s.db.Close()
}

// DestructiveReset drops all the tables and rebuilds them.
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

// AutoMigrate attemps to automatically migrate the tables.
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}).Error
}
