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

// ExitFn sets the exit function to use.
// Defaults to `os.Exit`.
//
// This is used in the unit tests but can be used to intercept the
// exit call and do something else as needed also.
func (c *contextBuilder) ExitFn(exitFn func(int)) *contextBuilder {
	c.exitFn = exitFn
	return c
}

// ForceExit sets whether to force terminate ungracefully upon receipt of a second signal.
// Defaults to true.
func (c *contextBuilder) ForceExit(forceExit bool) *contextBuilder {
	c.forceExit = forceExit
	return c
}

// Timeout sets the timeout for graceful shutdown before forcing the issue exiting with exit code 1.
// Defaults to 30 seconds.
//
// A timeout of <= 0, not recommended,
// disables the timeout and will wait forever for a seconds signal or application shuts down.
func (c *contextBuilder) Timeout(timeout time.Duration) *contextBuilder {
	c.timeout = timeout
	return c
}

// Signals sets the signals to listen for. Defaults to `os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT`.
func (c *contextBuilder) Signals(signals ...os.Signal) *contextBuilder {
	c.signals = signals
	return c
}
