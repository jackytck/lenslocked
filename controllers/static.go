package controllers

import "github.com/jackytck/lenslocked/views"

// NewStatic creates and returns a new static handler.
func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "views/static/home.gohtml"),
		Contact: views.NewView("bootstrap", "views/static/contact.gohtml"),
	}
}

// Static is the static http handler.
type Static struct {
	Home    *views.View
	Contact *views.View
}
