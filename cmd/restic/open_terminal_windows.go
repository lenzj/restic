// +build windows

package main

import "os"

// openTerminal opens the console input and screen buffers.
func openTerminal() (t *controllingTerminal, err error) {
	// https://docs.microsoft.com/en-us/windows/console/console-handles
	conin, err := os.Open("CONIN$")
	if err != nil {
		return nil, err
	}

	conout, err := os.OpenFile("CONOUT$", os.O_WRONLY, 0)
	if err != nil {
		conin.Close()
		return nil, err
	}

	return &controllingTerminal{conin, conout}, nil
}

type controllingTerminal struct {
	conin, conout *os.File
}

// Returns the input buffer's Handle for reading passwords.
func (t *controllingTerminal) Fd() uintptr {
	return t.conin.Fd()
}

func (t *controllingTerminal) Write(p []byte) (int, error) {
	return t.conout.Write(p)
}

func (t *controllingTerminal) Close() error {
	err1 := t.conin.Close()
	err2 := t.conout.Close()
	if err1 != nil {
		return err1
	}
	return err2
}
