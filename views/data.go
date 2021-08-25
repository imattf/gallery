package views

const (
	AlertLevelError   = "danger"
	AlterLevelWarning = "warning"
	AlterLevelInfo    = "info"
	AlterLevelSuccess = "success"

	// AlterMsgGeneric is displayed for random unexpected error
	AlertMsgGeneric = "Something went wrong... Please try again and contact us if the problem persists."
)

// Alter is used to render a Bootstrap Alert message
type Alert struct {
	Level   string
	Message string
}

// Data is the top level structure that views expect data
type Data struct {
	Alert *Alert
	Yield interface{}
}

func (d *Data) SetAlert(err error) {
	if pErr, ok := err.(PublicError); ok {
		d.Alert = &Alert{
			Level:   AlertLevelError,
			Message: pErr.Public(),
		}
	} else {
		d.Alert = &Alert{
			Level:   AlertLevelError,
			Message: AlertMsgGeneric,
		}
	}
}

type PublicError interface {
	error
	Public() string
}