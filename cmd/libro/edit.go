package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// RunEditSubcmd executes the "edit" sub-command.
func (app *App) RunEditSubcmd(args []string) error {
	fs := flag.NewFlagSet("edit", flag.ExitOnError)

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: %s %s [option...] DATA\n", app.Name(), fs.Name())
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	suffix := fs.String("suffix", "json", "sets the file extension to identify the type of data to be edited. Its main use is to use a useful syntax highlighting during edition")
	fs.StringVar(&app.Editor, "editor", os.Getenv("EDITOR"), "sets editor name to use to edit data")

	if err := fs.Parse(args); err != nil {
		return err
	}

	var data io.Reader
	switch fs.NArg() {
	case 0:
		if fi, _ := os.Stdin.Stat(); (fi.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
			return fmt.Errorf("invalid number of argument(s)\nRun %s %s -help", app.Name(), fs.Name())
		}
		data = os.Stdin
	case 1:
		data = strings.NewReader(fs.Arg(0))
	default:
		return fmt.Errorf("invalid number of argument(s)\nRun %s %s -help", app.Name(), fs.Name())
	}

	if _, err := app.edit(data, app.Stdout, *suffix); err != nil {
		return err
	}

	return nil
}

func (app *App) edit(r io.Reader, w io.Writer, suffix string) (int64, error) {
	if app.Editor == "" {
		return io.Copy(w, r)
	}

	tmpfile, err := os.CreateTemp("", "*."+suffix)
	if err != nil {
		return 0, err
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.ReadFrom(r); err != nil {
		return 0, err
	}

	if err := runInTTY(app.Editor, tmpfile.Name()); err != nil {
		return 0, err
	}

	if _, err := tmpfile.Seek(0, io.SeekStart); err != nil {
		return 0, err
	}

	return io.Copy(w, tmpfile)
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
