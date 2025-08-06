package errorsext

import (
	"context"
	"time"
)

// MaxAttemptsNonRetryableReset will apply the max attempts to all errors not determined to be retryable,
// but will reset the attempts if a retryable error is encountered after a non-retryable error.
const MaxAttemptsNonRetryableReset MaxAttemptsMode = iota

// MaxAttemptsMode is used to set the mode for the maximum number of attempts.
//
// E. g. Should the max attempts apply to all errors,
// just ones not determined to be retryable,
// reset on retryable errors, etc.
type MaxAttemptsMode uint8

// BackoffFn is a function used to apply a backoff strategy to the retryable function.
//
// It accepts `E` in cases where the amount of time to backoff is dynamic,
// for example when and http request fails with a 429 status code,
// the `Retry-After` header can be used to determine how long to backoff.
// It is not required to use or handle `E` and can be ignored if desired.
type BackoffFn[E any] func(ctx context.Context, attempt int, e E)

// EarlyReturnFn is the function that can be used to bypass all retry logic,
// no matter the MaxAttemptsMode,
// for when the type of `E` will never succeed and should not be retried.
//
// eg. If retrying an HTTP request and getting 400 Bad Request,
// it's unlikely to ever succeed and should not be retried.
type EarlyReturnFn[E any] func(ctx context.Context, e E) (earlyReturn bool)

// IsRetryableFn2 is called to determine if the type E is retryable.
type IsRetryableFn2[E any] func(ctx context.Context, e E) (isRetryable bool)

// Retryer is used to retry any fallible operation.
type Retryer[T, E any] struct {
	isRetryableFn   IsRetryableFn2[E]
	isEarlyReturnFn EarlyReturnFn[E]
	maxAttemptsMode MaxAttemptsMode
	maxAttempts     uint8
	bo              BackoffFn[E]
	timeout         time.Duration
}

// NewRetryer returns a new `Retryer` with sane default values.
//
// The default values are:
// - `MaxAttemptsMode` is `MaxAttemptsNonRetryableReset`.
// - `MaxAttempts` is 5.
// - `Timeout` is 0 no context timeout.
// - `IsRetryableFn` will always return false as `E` is unknown until defined.
// - `BackoffFn` will sleep for 200ms. It's recommended to use exponential backoff for production.
// - `EarlyReturnFn` will be None.
func NewRetryer[T, E any]() Retryer[T, E] {
	return Retryer[T, E]{
		isRetryableFn:   func(_ context.Context, _ E) bool { return false },
		maxAttemptsMode: MaxAttemptsNonRetryableReset,
		maxAttempts:     5,
		bo: func(ctx context.Context, attempt int, _ E) {
			t := time.NewTimer(time.Millisecond * 200)
			defer t.Stop()
			select {
			case <-ctx.Done():
			case <-t.C:
			}
		},
	}
}

// Backoff sets the backoff function for the `Retryer`.
func (r Retryer[T, E]) Backoff(fn BackoffFn[E]) Retryer[T, E] {
	if fn == nil {
		fn = func(_ context.Context, _ int, _ E) {}
	}

	r.bo = fn
	return r
}

// Timeout sets the timeout for the `Retryer`.
// This is the timeout per `RetyableFn` attempt and not the entirety of the `Retryer` execution.
// Timeout of 0 will disable the timeout and is the default.
func (r Retryer[T, E]) Timeout(timeout time.Duration) Retryer[T, E] {
	r.timeout = timeout
	return r
}
