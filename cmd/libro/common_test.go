package main

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/pirmd/verify"
)

const (
	testdata      = "../../testdata" //Use test data of the main package
	testdataBooks = testdata + "/books"
)

type testApp struct {
	*App
	*verify.TestFolder
}

func newTestApp(tb testing.TB) *testApp {
	testOut := new(bytes.Buffer)
	testLog := verify.NewLogger(tb)
	testDir := verify.MustNewTestFolder(tb)

	app := New()
	app.Stdout = testOut
	app.Verbose, app.Debug = testLog, testLog
	app.Library.Root = testDir.Root
	app.Formatter = template.Must(app.Formatter.Parse(`{{toPrettyJSON .}}`)) //Output to pretty JSON to easier testing failures analysis

	app.Var(NewGoTemplate(app.Formatter), "format", "set output format using golang text/template")

	return &testApp{
		App:        app,
		TestFolder: testDir,
	}
}

func (app *testApp) Out() string {
	return app.App.Stdout.(*bytes.Buffer).String()
}
