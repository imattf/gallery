package views

const (
	AlertLevelError   = "danger"
	AlterLevelWarning = "warning"
	AlterLevelInfo    = "info"
	AlterLevelSuccess = "success"
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
