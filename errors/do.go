package errorsext

// IsRetryableFn is called to determine if the error is retryable and optionally returns the reason for logging and metrics.
type IsRetryableFn[E any] func(err E) (reason string, isRetryable bool)
