package views

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
	Yield interface{}
}
