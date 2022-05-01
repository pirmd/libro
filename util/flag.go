package util

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
)

var (
	// ErrUnknownLogStatus is returned when trying to set a logswitcher's
	// status tat is neither 'true' nor 'false'.
	ErrUnknownLogStatus = errors.New("unknown status")
)

// Logswitcher wraps a set of log.Logger to implement flag.Value interface and
// get activation/de-activation from command-line.
// When activated, logger will print to os.Stderr.
type Logswitcher struct {
	loggers []*log.Logger
}

// NewLogSwitcher creates a new Logger that prints nothing (output is
// io.Discard), waiting for user to trigger it through command-line flags.
func NewLogSwitcher(loggers ...*log.Logger) *Logswitcher {
	return &Logswitcher{
		loggers: loggers,
	}
}

// String proposes a human-friendly string representation of a logswitcher.
func (ls Logswitcher) String() string {
	return ""
}

// Set implements flag.Value interface to set logswitcher's status from
// command-line. Status could be either true or false.
func (ls *Logswitcher) Set(status string) error {
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
func (ls Logswitcher) IsBoolFlag() bool {
	return true
}

// GoTemplate wraps a text/template.Template to implement flag.Value interface
// and get customization through command-line.
type GoTemplate struct {
	*template.Template
}

// NewGoTemplate creates a new GoTemplate.
func NewGoTemplate(tmpl *template.Template) *GoTemplate {
	return &GoTemplate{
		Template: tmpl,
	}
}

// String proposes a human-friendly string representation of a go template.
func (gotmpl GoTemplate) String() string {
	if gotmpl.Template != nil {
		return gotmpl.Template.Root.String()
	}
	return ""
}

// Set implements flag.Value interface for a GoTemplate.
func (gotmpl *GoTemplate) Set(tmpl string) error {
	if _, err := gotmpl.Template.Parse(tmpl); err != nil {
		return err
	}

	return nil
}

// GoTemplateFS wraps a text/template.Template that parses its template
// definition from files in order to implement flag.Value interface and get
// customization through command-line.
type GoTemplateFS struct {
	*template.Template
}

// NewGoTemplateFS creates a new GoTemplateFS.
func NewGoTemplateFS(tmpl *template.Template) *GoTemplateFS {
	return &GoTemplateFS{
		Template: tmpl,
	}
}

// String proposes a human-friendly string representation of a go template.
func (gotmplFS GoTemplateFS) String() string {
	if gotmplFS.Template != nil {
		return gotmplFS.Template.Root.String()
	}
	return ""
}

// Set implements flag.Value interface for a GoTemplateFS.
func (gotmplFS *GoTemplateFS) Set(pattern string) error {
	if _, err := gotmplFS.Template.ParseGlob(pattern); err != nil {
		return err
	}

	return nil
}

// KV wraps a map to implement flag.Value interface and get ability to allow
// user to define (key, value) through command-line.
type KV struct {
	kv map[string]string
}

// NewKV creates a new Map.
func NewKV(m map[string]string) *KV {
	return &KV{
		kv: m,
	}
}

// String proposes a human-friendly string representation of a collection of
// (key,value).
func (kv KV) String() string {
	return fmt.Sprint(kv.kv)
}

// Set implements flag.Value interface for a KV.
// Command-line flag format is key=value
func (kv *KV) Set(s string) error {
	arg := strings.SplitN(s, "=", 2)
	if len(arg) != 2 {
		return fmt.Errorf("argument is not in key=value format")
	}

	kv.kv[arg[0]] = arg[1]
	return nil
}
