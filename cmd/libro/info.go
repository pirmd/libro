package main

import (
	"flag"
	"fmt"
)

// RunInfoSubcmd executes the "info" sub-command.
func (app *App) RunInfoSubcmd(args []string) error {
	fs := flag.NewFlagSet("info", flag.ExitOnError)

	// re-use generic flag from main command
	// Note: generic flag are not post-processed by the sub-command therefore
	// the needed logic to process them should be handled by the flag Parse()
	// [for example using the flag.Value interface].
	app.FlagSet.VisitAll(func(f *flag.Flag) {
		fs.Var(f.Value, f.Name, f.Usage)
	})

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: %s %s [option...] FILENAME\n", app.Name(), fs.Name())
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	fs.BoolVar(&app.Library.UseGuesser, "use-guesser", false, "completes book's metadata by guessing lacking information from book's filename and title")
	fs.BoolVar(&app.Library.UseGooglebooks, "use-googlebooks", false, "completes book's metadata by searching lacking information from Googlebooks")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return fmt.Errorf("wrong number of arguments\nRun %s %s -help", app.Name(), fs.Name())
	}
	path := fs.Arg(0)

	b, err := app.Library.Read(path)
	if err != nil {
		return fmt.Errorf("cannot retrieve information about '%s': %w", path, err)
	}

	if err := app.Formatter.Execute(app.Stdout, b); err != nil {
		return fmt.Errorf("fail to display book information: %v", err)
	}
	fmt.Fprintln(app.Stdout)

	return nil
}
