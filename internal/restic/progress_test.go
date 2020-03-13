package restic

import (
	"testing"
	"time"

	rtest "github.com/restic/restic/internal/test"
)

func TestProgress(t *testing.T) {
	var started, updated, ticker, done bool

	p := &Progress{}
	p.d = 100 * time.Microsecond

	p.OnStart = func() {
		started = true
	}
	p.OnUpdate = func(_ Stat, _ time.Duration, t bool) {
		updated = true
		ticker = ticker || t
	}
	p.OnDone = func(Stat, time.Duration, bool) {
		done = true
	}

	p.Start()
	p.Report(Stat{})
	<-time.After(20 * time.Millisecond) // Wait for ticker to fire.
	p.Done()

	rtest.Assert(t, started, "OnStart not called")
	rtest.Assert(t, updated, "OnUpdate not called")
	rtest.Assert(t, ticker, "OnUpdate not called for ticker")
	rtest.Assert(t, done, "OnDone not called")
}
