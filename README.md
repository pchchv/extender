# extender [![CI](https://github.com/pchchv/extender/workflows/CI/badge.svg)](https://github.com/pchchv/extender/actions?query=workflow%3ACI+event%3Apush) [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/extender)](https://pkg.go.dev/github.com/pchchv/extender) [![Go Report Card](https://goreportcard.com/badge/github.com/pchchv/extender)](https://goreportcard.com/report/github.com/pchchv/extender)

The `extender` package extends the basic go packages by adding missing or additional functionality to them.
All packages match the std go package name with the optional *ext* suffix to avoid name conflicts.

## Contents:
- Detachable context.
- Result & Option types.
- Generic Mutex and RWMutex.
- Generic Doubly Linked List.
- Proper RFC3339Nano definition.
- Retrier for helping with any fallible operation.
- HTTP helper functions and constant placeholders.
- Bytes helper placeholders units eg. MB, MiB, GB...
- unsafe []byte->string & string->[]byte helper functions.