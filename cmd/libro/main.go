package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pirmd/libro/libro"
)

const (
	// defaultFormat is the default text/template for pretty-printing
	// retrieved information when printing to a pipeline.
	defaultFormat = "{{ toJSON . }}"

	// defaultFormatForHuman is the default text/template for pretty-printing
	// retrieved information when printing to a TTY.
	defaultFormatForHuman = "{{ . }}"
)

var (
	myname    = filepath.Base(os.Args[0])
	myversion = "v?.?.?-build" //should be set using: go build -ldflags "-X main.version=X.X.X"
)

// App is a wrapper around a libro.Libro object that implements command line
// facilities to interact with the user.
// App is also supposed to offer a not that complicated API for testing the
// app.
type App struct {
	*flag.FlagSet

	// Library points to the underlying libro.Libro object.
	Library *libro.Libro
	// Formatter is the go template used by UI to pretty print an object.
	Formatter *template.Template
	// Stdout is the standard output where to print app's result. It is usually
	// os.Sdtout except for test where you might want to capture output to a
	// buffer.
	Stdout io.Writer
}

func main() {
	tmpl := template.New("formatter").Funcs(SerializationFuncMap)

	app := &App{
		FlagSet:   flag.NewFlagSet(myname, flag.ExitOnError),
		Stdout:    os.Stdout,
		Library:   libro.New(),
		Formatter: template.Must(tmpl.Parse(defaultFormat)),
	}

	// if we are printing to a TTY, we use a format that is easier to read for a human.
	if fi, _ := os.Stdout.Stat(); (fi.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		template.Must(app.Formatter.Parse(defaultFormatForHuman))
	}

	app.Usage = func() {
		fmt.Fprintf(app.Output(), "Usage: %s [option...] <commands> [arguments]\n", myname)
		fmt.Fprintf(app.Output(), "Commands:\n")
		fmt.Fprintf(app.Output(), "    info       retrieve information from an epub\n")
		fmt.Fprintf(app.Output(), "    insert     insert an epub to the library\n")
		fmt.Fprintf(app.Output(), "    edit       edit information\n")
		fmt.Fprintf(app.Output(), "    version    print %s version\n", myname)
		fmt.Fprintf(app.Output(), "Options:\n")
		app.PrintDefaults()
	}

	app.Var(NewLogSwitcher(app.Library.Verbose), "verbose", "print messages of low interest")
	app.Var(NewLogSwitcher(app.Library.Debug, app.Library.Verbose), "debug", "print cryptic messages supposed to help the developer understand his/her mistakes")
	app.Var(NewGoTemplate(app.Formatter), "format", "set output format using golang text/template")

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
