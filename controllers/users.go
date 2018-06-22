package controllers

import (
	"fmt"
	"net/http"

	"github.com/jackytck/lenslocked/views"
)

// NewUsers return new users instance.
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
	}
}

// Users represent a set of users.
type Users struct {
	NewView *views.View
}

// New render the form where a user can create a new user account.
//
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}

// Create processes the signup form when a user tries to create a new user account.
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is a fake message. Pretend that we created the user account!")
}
