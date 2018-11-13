package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	llctx "github.com/jackytck/lenslocked/context"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/jackytck/lenslocked/models"
	"golang.org/x/oauth2"
)

// NewOAuths is used to create a new OAuths controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup.
func NewOAuths(os models.OAuthService, configs map[string]*oauth2.Config) *OAuths {
	return &OAuths{
		os:      os,
		configs: configs,
	}
}

// OAuths represent a set of users.
type OAuths struct {
	os      models.OAuthService
	configs map[string]*oauth2.Config
}

// Connect redirects to the dropbox oauth2 endpoint page.
func (o *OAuths) Connect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	oauthConfig, ok := o.configs[service]
	if !ok {
		http.Error(w, "Invalid OAuth2 Service", http.StatusBadRequest)
		return
	}

	state := csrf.Token(r)
	cookie := http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	url := oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

// Callback handles the oauth callback from dropbox.
func (o *OAuths) Callback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	oauthConfig, ok := o.configs[service]
	if !ok {
		http.Error(w, "Invalid OAuth2 Service", http.StatusBadRequest)
		return
	}

	r.ParseForm()
	state := r.FormValue("state")
	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if cookie == nil || cookie.Value != state {
		http.Error(w, "Invalid state provided", http.StatusBadRequest)
		return
	}
	cookie.Value = ""
	cookie.Expires = time.Now()
	http.SetCookie(w, cookie)

	code := r.FormValue("code")
	token, err := oauthConfig.Exchange(context.TODO(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := llctx.User(r.Context())
	existing, err := o.os.Find(user.ID, service)
	if err == models.ErrNotFound {
		// no op
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		o.os.Delete(existing.ID)
	}
	userOAuth := models.OAuth{
		UserID:  user.ID,
		Token:   *token,
		Service: service,
	}
	err = o.os.Create(&userOAuth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%+v", token)
}

// DropboxTest handles the tests.
func (o *OAuths) DropboxTest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]

	r.ParseForm()
	path := r.FormValue("path")

	user := llctx.User(r.Context())
	userOAuth, err := o.os.Find(user.ID, service)
	if err != nil {
		panic(err)
	}
	token := userOAuth.Token

	config := dropbox.Config{
		Token: token.AccessToken,
	}
	dbx := files.New(config)
	res, err := dbx.ListFolder(&files.ListFolderArg{
		Path: path,
	})
	if err != nil {
		panic(err)
	}
	for _, entry := range res.Entries {
		switch meta := entry.(type) {
		case *files.FolderMetadata:
			fmt.Fprintln(w, "FolderMetadata=", meta)
		case *files.FileMetadata:
			fmt.Fprintln(w, "FileMetadata=", meta)
		}
	}
}
