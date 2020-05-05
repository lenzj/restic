// +build !windows

package main

import (
	"os"
	"os/exec"
)

var devtty = "/dev/tty" // Variable so that the test can reset it.

// openTerminal opens the controlling terminal.
func openTerminal() (*controllingTerminal, error) {
	return os.OpenFile(devtty, os.O_RDWR, 0)
}

type controllingTerminal = os.File

// passTerminal passes tty as stdin and stdout to cmd.
func passTerminal(cmd *exec.Cmd, tty *controllingTerminal) {
	cmd.Stdin = tty
	cmd.Stdout = tty
}
