package util

import (
	"os"
	"os/exec"
)

// ExecInTTY executes a command in a TTY.
func ExecInTTY(name string, arg ...string) error {
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
