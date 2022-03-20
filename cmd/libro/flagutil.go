package main

import (
	"errors"
	"io"
	"log"
	"os"
	"text/template"
)

var (
	// ErrUnknownStatus is returned when trying to set a logswitcher's status
	// tat is neither 'true' nor 'false'.
	ErrUnknownLogStatus = errors.New("unknown status")
)

// logswitcher wraps a set of log.Logger to implement flag.Value interface and
// get activation/de-activation from command-line.
// When activated, logger will print to os.Stderr.
type logswitcher struct {
	loggers []*log.Logger
}

// NewLogSwitcher creates a new Logger that prints nothing (output is
// io.Discard), waiting for user to trigger it through command-line flags.
func NewLogSwitcher(loggers ...*log.Logger) *logswitcher {
	return &logswitcher{
		loggers: loggers,
	}
}

// String proposes a human-friendly string representation of a logswitcher.
func (ls logswitcher) String() string {
	return ""
}

// Set implements flag.Value interface to set logswitcher's status from
// command-line. Status could be either true or false.
func (ls *logswitcher) Set(status string) error {
	switch status {
	case "true":
		for _, logger := range ls.loggers {
			logger.SetOutput(os.Stderr)
		}
		return nil

	case "false":
		for _, logger := range ls.loggers {
			logger.SetOutput(io.Discard)
		}
		return nil
	}

	return ErrUnknownLogStatus
}

// IsBoolFlag implements flag.Value interface to notify that logswitcher flag
// is boolean.
func (ls logswitcher) IsBoolFlag() bool {
	return true
}

// gotemplate wraps a text/template.Template to implement flag.Value interface
// and get customization through command-line.
type gotemplate struct {
	*template.Template
}

// NewGoTemplate creates a new gotemplate.
func NewGoTemplate(tmpl *template.Template) *gotemplate {
	return &gotemplate{
		Template: tmpl,
	}
}

// String proposes a human-friendly string representation of a formatter.
func (gotmpl gotemplate) String() string {
	if gotmpl.Template != nil {
		return gotmpl.Template.Root.String()
	}
	return ""
}

// Set implements flag.Value interface for a gotemplate.
func (gotmpl *gotemplate) Set(tmpl string) error {
	if _, err := gotmpl.Template.Parse(tmpl); err != nil {
		return err
	}

	return nil
}
