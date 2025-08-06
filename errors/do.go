package errorsext

import (
	"context"

	optionext "github.com/pchchv/extender/values/option"
	resultext "github.com/pchchv/extender/values/result"
)

// IsRetryableFn is called to determine if the error is retryable and optionally returns the reason for logging and metrics.
type IsRetryableFn[E any] func(err E) (reason string, isRetryable bool)

// RetryableFn is a function that can be retried.
type RetryableFn[T, E any] func(ctx context.Context) resultext.Result[T, E]

// OnRetryFn is called after IsRetryableFn returns true and before the retry is attempted.
//
// this allows for interception, short-circuiting and adding of backoff strategies.
type OnRetryFn[E any] func(ctx context.Context, originalErr E, reason string, attempt int) optionext.Option[E]
