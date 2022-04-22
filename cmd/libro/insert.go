package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pirmd/libro/book"
)

// RunInsertSubcmd executes the "insert" sub-command.
func (app *App) RunInsertSubcmd(args []string) error {
	fs := flag.NewFlagSet("insert", flag.ExitOnError)

	// re-use generic flag from main command
	// Note: generic flag are not post-processed by the sub-command therefore
	// the needed logic to process them should be handled by the flag Parse()
	// [for example using the flag.Value interface].
	app.FlagSet.VisitAll(func(f *flag.Flag) {
		fs.Var(f.Value, f.Name, f.Usage)
	})

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: %s %s [option...] BOOKinJSON\n", app.Name(), fs.Name())
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	fs.StringVar(&app.Library.Root, "root", app.Library.Root, "root folder where the books library is to be found")
	fs.Var(NewGoTemplate(app.Library.PathTmpl), "rename", "sets filename format using golang text/template")
	fs.Var(NewGoTemplateFS(app.Library.PathTmpl), "rename-tmpl", "loads user-defined filename template(s) from golang text/template definition files")

	if err := fs.Parse(args); err != nil {
		return err
	}

	var bookJSON io.Reader
	switch fs.NArg() {
	case 0:
		if fi, _ := os.Stdin.Stat(); (fi.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
			return fmt.Errorf("invalid number of argument(s)\nRun %s %s -help", app.Name(), fs.Name())
		}
		bookJSON = os.Stdin
	case 1:
		bookJSON = strings.NewReader(fs.Arg(0))
	default:
		return fmt.Errorf("invalid number of argument(s)\nRun %s %s -help", app.Name(), fs.Name())
	}

	b := book.New()
	if err := json.NewDecoder(bookJSON).Decode(&b); err != nil {
		return fmt.Errorf("fail to decode book's JSON: %v", err)
	}

	if err := app.Library.Create(b); err != nil {
		return fmt.Errorf("fail to add new book: %v", err)
	}

	if err := app.Formatter.Execute(app.Stdout, b); err != nil {
		return fmt.Errorf("fail to display book information: %v", err)
	}
	fmt.Fprint(app.Stdout)

	return nil
}
