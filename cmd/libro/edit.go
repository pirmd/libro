package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pirmd/libro/book"
)

// RunEditSubcmd executes the "edit" sub-command.
func (app *App) RunEditSubcmd(args []string) error {
	fs := flag.NewFlagSet("edit", flag.ExitOnError)

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

	auto := fs.Bool("auto", false, "do not trigger an editor if libro get the impression that Book's information is good enough")
	editor := fs.String("editor", os.Getenv("EDITOR"), "sets editor's name to use for editing Book's information")

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

	app.Library.Verbose.Printf("Edit book's information")

	switch {
	case *editor == "":
		app.Library.Verbose.Printf("no editor has been defined. Set $EDITOR global var or use -editor command line flag")
	case *auto && b.IsComplete():
		app.Library.Verbose.Printf("no need to edit book's information that seems good enough to me")
	default:
		var err error
		if b, err = editBook(*editor, b); err != nil {
			return fmt.Errorf("fail to edit book: %v", err)
		}
	}

	if err := app.Formatter.Execute(app.Stdout, b); err != nil {
		return fmt.Errorf("fail to display book information: %v", err)
	}
	fmt.Fprintln(app.Stdout)

	return nil
}

func editBook(editor string, b *book.Book) (*book.Book, error) {
	tmpfile, err := os.CreateTemp("", "*.json")
	if err != nil {
		return nil, err
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()

	prettyJSON := json.NewEncoder(tmpfile)
	prettyJSON.SetIndent("", "  ")
	if err := prettyJSON.Encode(b); err != nil {
		return nil, err
	}

	if err := runInTTY(editor, tmpfile.Name()); err != nil {
		return nil, err
	}

	if _, err := tmpfile.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	edbook := book.New()
	if err := json.NewDecoder(tmpfile).Decode(&edbook); err != nil {
		return nil, err
	}

	return edbook, nil
}

func runInTTY(name string, arg ...string) error {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer tty.Close()

	c := exec.Command(name, arg...)
	c.Stdin = tty
	c.Stdout = tty
	c.Stderr = tty

	return c.Run()
}
