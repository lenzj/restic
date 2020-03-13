// +build darwin freebsd netbsd openbsd dragonfly

package progress

import (
	"os"
	"os/signal"
	"syscall"
)

func (p *Progress) initSignals() {
	p.signal = make(chan os.Signal, 1)
	signal.Notify(p.signal, syscall.SIGINFO, syscall.SIGUSR1)
}
