package models

import "github.com/jinzhu/gorm"

// Gallery is our image container resources that visitors view.
type Gallery struct {
	gorm.Model
	UserID uint     `gorm:"not_null;index"`
	Title  string   `gorm:"not_null"`
	Images []string `gorm:"-"`
}

type GalleryService interface {
	GalleryDB
}

type GalleryDB interface {
	ByID(id uint) (*Gallery, error)
	ByUserID(id uint) ([]Gallery, error)
	Create(gallery *Gallery) error
	Update(gallery *Gallery) error
	Delete(id uint) error
}

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{&galleryGorm{db}},
	}
}

type galleryService struct {
	GalleryDB
}

type galleryValidator struct {
	GalleryDB
}

func (gv *galleryValidator) Create(gallery *Gallery) error {
	err := runGalleryValFuncs(gallery,
		gv.userIDRequired,
		gv.titleRequired,
	)
	if err != nil {
		return err
	}

	return gv.GalleryDB.Create(gallery)
}

func (gv *galleryValidator) Update(gallery *Gallery) error {
	err := runGalleryValFuncs(gallery,
		gv.userIDRequired,
		gv.titleRequired,
	)
	if err != nil {
		return err
	}

	return gv.GalleryDB.Update(gallery)
}

// Delete deletes the gallery with the provided ID.
func (gv *galleryValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}
	return gv.GalleryDB.Delete(id)
}

func (gv *galleryValidator) userIDRequired(g *Gallery) error {
	if g.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (gv *galleryValidator) titleRequired(g *Gallery) error {
	if g.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

type galleryValFunc func(*Gallery) error

func runGalleryValFuncs(gallery *Gallery, fns ...galleryValFunc) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
			return err
		}
	}
	return nil
}

var _ GalleryDB = &galleryGorm{}

type galleryGorm struct {
	db *gorm.DB
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}

func (gg *galleryGorm) Update(gallery *Gallery) error {
	return gg.db.Save(gallery).Error
}

// Delete deletes the gallery with the provided ID.
func (gg *galleryGorm) Delete(id uint) error {
	g := Gallery{Model: gorm.Model{ID: id}}
	return gg.db.Delete(&g).Error
}

func (gg *galleryGorm) ByID(id uint) (*Gallery, error) {
	var g Gallery
	db := gg.db.Where("id = ?", id)
	err := first(db, &g)
	return &g, err
}

func (gg *galleryGorm) ByUserID(id uint) ([]Gallery, error) {
	var galleries []Gallery
	gg.db.Where("user_id = ?", id).Find(&galleries)
	return galleries, nil
}
