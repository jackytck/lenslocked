package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jackytck/lenslocked/context"
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

// GalleryForm represents the form data of create gallery page.
type GalleryForm struct {
	Title string `schema:"title"`
}

// Create processes the gallery form when a user tries to create a new gallery.
//
// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}

	user := context.User(r.Context())
	fmt.Println("Create got the user:", user)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	fmt.Fprintln(w, gallery)
}
