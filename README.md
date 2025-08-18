# extender [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/extender)](https://pkg.go.dev/github.com/pchchv/extender)

Package that extends the basic go packages by adding missing or additional functionality to them. All packages match the std go package name with the optional `ext` suffix to avoid name conflicts.

## Contents:
- Detachable context.
- Result & Option types.
- Generic Mutex and RWMutex.
- Generic Doubly Linked List.
- Proper RFC3339Nano definition.
- Retrier for helping with any fallible operation.
- HTTP helper functions and constant placeholders.
- Bytes helper placeholders units eg. MB, MiB, GB, ...
- unsafe []byte->string & string->[]byte helper functions.