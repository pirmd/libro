package main

import (
	"bytes"
	"flag"
	"testing"
	"text/template"

	"github.com/pirmd/libro/libro"

	"github.com/pirmd/verify"
)

const (
	testdataBooks = "../../testdata/books"   //Use test data of the main package
	testFormat    = "{{ toPrettyJSON . }}\n" //Output to pretty JSON to easier testing failures analysis
)

type testApp struct {
	*App
	*verify.TestFolder

	out *bytes.Buffer
}

func newTestApp(tb testing.TB) *testApp {
	testOut := new(bytes.Buffer)
	testLog := verify.NewLogger(tb)
	testDir := verify.MustNewTestFolder(tb)

	testLib := libro.New()
	testLib.Root = testDir.Root
	testLib.Verbose, testLib.Debug = testLog, testLog

	tmpl := template.New("formatter").Funcs(SerializationFuncMap)
	app := &App{
		FlagSet:   flag.NewFlagSet("libro-testing", flag.ExitOnError),
		Library:   testLib,
		Formatter: template.Must(tmpl.Parse(testFormat)),
		Stdout:    testOut,
	}

	return &testApp{
		App:        app,
		TestFolder: testDir,
		out:        testOut,
	}
}
