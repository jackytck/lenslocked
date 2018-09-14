package views

import "github.com/jackytck/lenslocked/models"

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
	Yield interface{}
}

// SetAlert sets the alert from PublicError if possible.
func (d *Data) SetAlert(err error) {
	if pErr, ok := err.(PublicError); ok {
		d.AlertError(pErr.Public())
	} else {
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
