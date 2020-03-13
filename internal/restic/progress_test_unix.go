// +build !windows

package restic

import (
	"syscall"
	"testing"
	"time"

	rtest "github.com/restic/restic/internal/test"

	"golang.org/x/sys/unix"
)

func TestProgressSignal(t *testing.T) {
	var updated bool

	p := &Progress{}

	p.OnUpdate = func(Stat, time.Duration, bool) {
		updated = true
	}

	p.Start()
	unix.Kill(unix.Getpid(), syscall.SIGUSR1)
	p.Done()

	rtest.Assert(t, updated, "OnUpdate not called for SIGUSR1")
}
