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
