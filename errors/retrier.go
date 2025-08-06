package errorsext

// MaxAttemptsMode is used to set the mode for the maximum number of attempts.
//
// E. g. Should the max attempts apply to all errors, just ones not determined to be retryable, reset on retryable errors, etc.
type MaxAttemptsMode uint8
