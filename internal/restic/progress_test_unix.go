package restic

import (
	"io/ioutil"
	"os"
	"syscall"
	"testing"
	"time"

	rtest "github.com/restic/restic/internal/test"

	"golang.org/x/sys/unix"
)

func TestProgressSignal(t *testing.T) {
	var updated bool

	tmp, err := ioutil.TempFile("", "restic-fake-terminal")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	p := &Progress{}

	p.OnUpdate = func(Stat, time.Duration, bool) {
		updated = true
	}

	p.Start()
	unix.Kill(unix.Getpid(), syscall.SIGUSR1)
	//p.Report(Stat{})
	//<-time.After(20 * time.Millisecond) // Wait for ticker to fire.
	p.Done()

	rtest.Assert(t, updated, "OnUpdate not called for SIGUSR1")
}
