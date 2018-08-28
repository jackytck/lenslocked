package controllers

import (
	"github.com/jackytck/lenslocked/models"
	"github.com/jackytck/lenslocked/views"
)

// NewGalleries is used to create a new Galleries controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup.
func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New: views.NewView("bootstrap", "galleries/new"),
		gs:  gs,
	}
}

// Galleries represent a Galleries controller.
type Galleries struct {
	New *views.View
	gs  models.GalleryService
}
