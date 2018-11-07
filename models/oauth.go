package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
)

// OAuth stores the oauth token of a user.
type OAuth struct {
	gorm.Model
	UserID  uint   `gorm:"not_null;unique_index:user_id_service"`
	Service string `gorm:"not_null;unique_index:user_id_service"`
	oauth2.Token
}

type OAuthService interface {
	OAuthDB
}

type OAuthDB interface {
	Find(userID uint, service string) (*OAuth, error)
	Create(oauth *OAuth) error
	Delete(id uint) error
}

func NewOAuthService(db *gorm.DB) OAuthService {
	return nil
}
