// Package epubcheck provides basic bindings to run and interpret `epubcheck`
// tool.
// Features are compatible with epubcheck v4 and correspond to a subset of
// https://github.com/w3c/epubcheck/tree/main/docs  and
// https://www.w3.org/publishing/epubcheck/docs/messages/
package epubcheck

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	// Executable contains the path to EPUBcheck binary.
	Executable = "epubcheck"
)

// Location contains a description of where a Message is associated to.
type Location struct {
	Path   string
	Line   int
	Column int
}

// Message represents an EPUBcheck report message.
type Message struct {
	ID         string
	Severity   string
	Message    string
	Suggestion string
	Locations  []Location
}

// String provides a text representation of a Message.
func (msg Message) String() string {
	var loc []string
	for _, l := range msg.Locations {
		loc = append(loc, fmt.Sprintf("%s (%d, %d)", l.Path, l.Line, l.Column))
	}
	location := strings.Join(loc, ": ")

	return fmt.Sprintf("%s (%s): %s: %s; %s", msg.Severity, msg.ID, location, msg.Message, msg.Suggestion)
}

// Report represents an EPUBcheck's report.
type Report struct {
	Messages []*Message
}

// FilterMessagesByID filters a Report to keep only Messages whose ID match the
// given glob pattern.
func (r Report) FilterMessagesByID(pattern string) *Report {
	filteredReport := new(Report)

	for _, m := range r.Messages {
		matched, err := filepath.Match(pattern, m.ID)
		if err != nil {
			continue
		}

		if matched {
			filteredReport.Messages = append(filteredReport.Messages, m)
		}
	}

	return filteredReport
}

// FilterMessagesBySeverity filters a Report to keep only Messages whose
// Severity match the given glob pattern.
func (r Report) FilterMessagesBySeverity(pattern string) *Report {
	filteredReport := new(Report)

	for _, m := range r.Messages {
		matched, err := filepath.Match(pattern, m.Severity)
		if err != nil {
			continue
		}

		if matched {
			filteredReport.Messages = append(filteredReport.Messages, m)
		}
	}

	return filteredReport
}

// Run executes EPUBcheck on the given EPUB.
// Additional options are added to EPUBcheck command line. By default, Run uses
// "--json -" command line to capture EPUBcheck report, therefore any additional
// command line argument that is not consistent with it shall be avoided.
func Run(path string, options ...string) (*Report, error) {
	args := append(options, "--quiet", "--json", "-", path)

	buf := new(bytes.Buffer)
	//#nosec G204 -- Executable is a constant and args is controlled internally (no user input).
	cmd := exec.Command(Executable, args...)
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil && buf.Len() == 0 {
		return nil, err
	}

	r := new(Report)
	if err := json.NewDecoder(buf).Decode(r); err != nil {
		return nil, err
	}

	return r, nil
}
