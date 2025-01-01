package severity

import "github.com/fatih/color"

type Severity int

const (
	Info Severity = iota
	Warning
	Error
	Critical
)

func (s Severity) String() string {
	switch s {
	case Info:
		return "Info"
	case Warning:
		return "Warning"
	case Error:
		return "Error"
	case Critical:
		return "Critical"
	default:
		return "Unknown"
	}
}

func (sev Severity) Color() *color.Color {
	switch sev {
	case Severity(Error):
		return color.New(color.FgRed)
	case Severity(Warning):
		return color.New(color.FgYellow)
	case Severity(Info):
		return color.New(color.FgCyan)
	case Severity(Critical):
		return color.New(color.FgHiMagenta)
	default:
		return color.New(color.Reset)
	}
}
