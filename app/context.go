package appext

import (
	"os"
	"syscall"
	"time"
)

type contextBuilder struct {
	signals   []os.Signal
	timeout   time.Duration
	exitFn    func(int)
	forceExit bool
}

// Context returns a new context builder, with sane defaults,
// that can be overridden. Calling `Build()` finalizes
// the new desired context and returns the configured `context.Context`.
func Context() *contextBuilder {
	return &contextBuilder{
		signals: []os.Signal{
			os.Interrupt,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		},
		timeout:   30 * time.Second,
		forceExit: true,
		exitFn:    os.Exit,
	}
}
