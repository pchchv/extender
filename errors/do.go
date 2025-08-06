package errorsext

import (
	"context"

	resultext "github.com/pchchv/extender/values/result"
)

// IsRetryableFn is called to determine if the error is retryable and optionally returns the reason for logging and metrics.
type IsRetryableFn[E any] func(err E) (reason string, isRetryable bool)

// RetryableFn is a function that can be retried.
type RetryableFn[T, E any] func(ctx context.Context) resultext.Result[T, E]
