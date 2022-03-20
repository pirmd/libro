package epubcheck

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"path/filepath"
)

var (
	Executable = "epubcheck"
)

type Location struct {
	Path   string
	Line   int
	Column int
}

type Message struct {
	ID         string
	Severity   string
	Message    string
	Suggestion string
	Locations  []Location
}

type Report struct {
	Messages []*Message
}

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

func CanRun() (bool, error) {
	_, err := exec.LookPath(Executable)
	return (err == nil), err
}

func Check(path string) (*Report, error) {
	buf := new(bytes.Buffer)
	cmd := exec.Command(Executable, "-json", "-", path)
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	r := new(Report)
	if err := json.Unmarshal(buf.Bytes(), r); err != nil {
		return nil, err
	}

	return r, nil
}
