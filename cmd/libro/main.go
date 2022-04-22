package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pirmd/libro/book"
	"github.com/pirmd/libro/libro"
	"github.com/pirmd/libro/util"
)

var (
	myname    = filepath.Base(os.Args[0])
	myversion = "v?.?.?-?" //should be set using: go build -ldflags "-X main.myversion=X.X.X"

	//go:embed gotmpl/*
	bookTmplDir embed.FS
)

// App is a wrapper around a libro.Libro object that implements command line
// facilities to interact with the user.
// App is also supposed to offer a not that complicated API for testing the
// app.
type App struct {
	*flag.FlagSet

	// Verbose is the logger for providing low interest messages to the user.
	Verbose *log.Logger

	// Debug is the logger for providing messages supposed to help the
	// developer understand his/her mistakes.
	Debug *log.Logger

	// Library points to the underlying libro.Libro object.
	Library *libro.Libro

	// Formatter is the go template used by UI to pretty print an object.
	Formatter *template.Template

	// Stdout is the standard output where to print app's result. It is usually
	// os.Sdtout except for test where you might want to capture output to a
	// buffer.
	Stdout io.Writer
}

// New creates a new App
func New() *App {
	tmpl := template.New("formatter").Option("missingkey=error")
	tmpl = tmpl.Funcs(util.SerializationFuncMap).Funcs(util.StringsFuncMap)
	tmpl = template.Must(tmpl.ParseFS(bookTmplDir, "gotmpl/*"))

	app := &App{
		FlagSet:   flag.NewFlagSet(myname, flag.ExitOnError),
		Verbose:   log.New(io.Discard, "", 0),
		Debug:     log.New(io.Discard, "debug:", 0),
		Stdout:    os.Stdout,
		Library:   libro.New(),
		Formatter: template.Must(tmpl.Parse(`{{toJSON .}}`)),
	}

	app.Library.Verbose, app.Library.Debug = app.Verbose, app.Debug
	book.Verbose, book.Debug = app.Verbose, app.Debug

	return app
}

func main() {
	app := New()

	// if we are printing to a TTY, we use a format that is easier to read for a human.
	if fi, _ := os.Stdout.Stat(); (fi.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		template.Must(app.Formatter.Parse(`{{template "plaintext" .}}`))
	}

	app.Usage = func() {
		fmt.Fprintf(app.Output(), "Usage: %s [option...] <commands> [arguments]\n", myname)
		fmt.Fprintf(app.Output(), "Commands:\n")
		fmt.Fprintf(app.Output(), "    info       retrieve information from an epub\n")
		fmt.Fprintf(app.Output(), "    insert     insert an epub to the library\n")
		fmt.Fprintf(app.Output(), "    edit       edit information about an epub\n")
		fmt.Fprintf(app.Output(), "    version    print %s version\n", myname)
		fmt.Fprintf(app.Output(), "Options:\n")
		app.PrintDefaults()
	}

	app.Var(NewLogSwitcher(app.Verbose), "verbose", "print messages of low interest")
	app.Var(NewLogSwitcher(app.Debug, app.Verbose), "debug", "print cryptic messages supposed to help the developer understand his/her mistakes")
	app.Var(NewGoTemplate(app.Formatter), "format", "set output format using golang text/template")
	app.Var(NewGoTemplateFS(app.Formatter), "format-tmpl", "loads user-defined format template(s) from golang text/template definition files")

	if err := app.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(app.Output(), "err: wrong arguments\nRun %s -help\n", app.Name())
		os.Exit(1)
	}

	switch cmd := app.Arg(0); cmd {
	case "version":
		fmt.Fprintf(app.Output(), "%s version %s\n", myname, myversion)

	case "info":
		if err := app.RunInfoSubcmd(app.Args()[1:]); err != nil {
			fmt.Fprintf(app.Output(), "err: %v\n", err)
			os.Exit(1)
		}

	case "insert", "add":
		if err := app.RunInsertSubcmd(app.Args()[1:]); err != nil {
			fmt.Fprintf(app.Output(), "err: %v\n", err)
			os.Exit(1)
		}

	case "edit":
		if err := app.RunEditSubcmd(app.Args()[1:]); err != nil {
			fmt.Fprintf(app.Output(), "err: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(app.Output(), "err: '%[1]s %s' unknown command\nRun %[1]s -help\n", app.Name(), cmd)
		os.Exit(1)
	}
}
