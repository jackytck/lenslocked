package controllers

import "github.com/jackytck/lenslocked/views"

// NewStatic creates and returns a new static handler.
func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "static/home"),
		Contact: views.NewView("bootstrap", "static/contact"),
	}
}

// Static is the static http handler.
type Static struct {
	Home    *views.View
	Contact *views.View
}
