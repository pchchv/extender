package errorsext

import "context"

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
