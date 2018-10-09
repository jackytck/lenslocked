package views

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/jackytck/lenslocked/models"
)

const (
	// AlertLvlError represents Bootstrap danger alert.
	AlertLvlError = "danger"
	// AlertLvlWarning represents Bootstrap warning alert.
	AlertLvlWarning = "warning"
	// AlertLvlInfo represents Bootstrap info alert.
	AlertLvlInfo = "info"
	// AlertLvlSuccess represents Bootstrap success alert.
	AlertLvlSuccess = "success"
	// AlertMsgGeneric is displayed when any random error
	// is encountered by our backend.
	AlertMsgGeneric = "Something went wrong. Please try again, and contact us if the problem persists."
)

// Alert is used to render Bootstrap Alert messages in templates.
type Alert struct {
	Level   string
	Message string
}

// Data is the top level structure that views expect data to come in.
type Data struct {
	Alert *Alert
	User  *models.User
	CSRF  template.HTML
	Yield interface{}
}

// SetAlert sets the alert from PublicError if possible.
func (d *Data) SetAlert(err error) {
	if pErr, ok := err.(PublicError); ok {
		d.AlertError(pErr.Public())
	} else {
		log.Println(err)
		d.AlertError(AlertMsgGeneric)
	}
}

// AlertError sets an alert level error.
func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

// PublicError shows the public error only.
type PublicError interface {
	error
	Public() string
}

func persistAlert(w http.ResponseWriter, alert Alert) {
	expiresAt := time.Now().Add(5 * time.Minute)
	lv := http.Cookie{
		Name:     "alert_level",
		Value:    alert.Level,
		Expires:  expiresAt,
		HttpOnly: true,
	}
	msg := http.Cookie{
		Name:     "alert_message",
		Value:    alert.Message,
		Expires:  expiresAt,
		HttpOnly: true,
	}
	http.SetCookie(w, &lv)
	http.SetCookie(w, &msg)
}

func clearAlert(w http.ResponseWriter) {
	lv := http.Cookie{
		Name:     "alert_level",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	msg := http.Cookie{
		Name:     "alert_message",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &lv)
	http.SetCookie(w, &msg)
}

func getAlert(r *http.Request) *Alert {
	lv, err := r.Cookie("alert_level")
	if err != nil {
		return nil
	}
	msg, err := r.Cookie("alert_message")
	if err != nil {
		return nil
	}
	alert := Alert{
		Level:   lv.Value,
		Message: msg.Value,
	}
	return &alert
}

// RedirectAlert accepts all the normal params for an
// http.Redirect and performs a redirect, but only after
// persisting the provided alert in a cookie so that it can
// be displayed when the new page is loaded.
func RedirectAlert(w http.ResponseWriter, r *http.Request, url string, code int, alert Alert) {
	persistAlert(w, alert)
	http.Redirect(w, r, url, code)
}
