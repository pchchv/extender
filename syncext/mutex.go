package syncext

import (
	"sync"

	"github.com/pchchv/express/resultext"
)

// MutexGuard protects the inner contents of a Mutex for safety and unlocking.
type MutexGuard[T any, M interface{ Unlock() }] struct {
	m M
	T T // is the inner generic type of the Mutex
}

// Unlock unlocks the Mutex value.
func (g MutexGuard[T, M]) Unlock() {
	g.m.Unlock()
}

// Mutex creates a type safe mutex wrapper ensuring one cannot access the
// values of a locked values without first gaining a lock.
type Mutex[T any] struct {
	m     *sync.Mutex
	value T
}

// NewMutex creates a new Mutex for use.
func NewMutex[T any](value T) Mutex[T] {
	return Mutex[T]{
		m:     new(sync.Mutex),
		value: value,
	}
}

// Lock locks the Mutex and returns value for use,
// safe for mutation if the lock is already in use,
// the calling goroutine blocks until the mutex is available.
func (m Mutex[T]) Lock() MutexGuard[T, *sync.Mutex] {
	m.m.Lock()
	return MutexGuard[T, *sync.Mutex]{
		m: m.m,
		T: m.value,
	}
}

// PerformMut safely locks and unlocks the Mutex values and performs the provided function returning its error if one
// otherwise setting the returned value as the new mutex value.
func (m Mutex[T]) PerformMut(f func(T)) {
	guard := m.Lock()
	f(guard.T)
	guard.Unlock()
}

// TryLock tries to lock Mutex and reports whether it succeeded.
// If it does the value is returned for use in the Ok result otherwise Err with empty value.
func (m Mutex[T]) TryLock() resultext.Result[MutexGuard[T, *sync.Mutex], struct{}] {
	if m.m.TryLock() {
		return resultext.Ok[MutexGuard[T, *sync.Mutex], struct{}](MutexGuard[T, *sync.Mutex]{
			m: m.m,
			T: m.value,
		})
	} else {
		return resultext.Err[MutexGuard[T, *sync.Mutex], struct{}](struct{}{})
	}
}

// RMutexGuard protects the inner contents of a RWMutex for safety and unlocking.
type RMutexGuard[T any] struct {
	rw *sync.RWMutex
	// T is the inner generic type of the Mutex
	T T
}

// RUnlock unlocks the RWMutex value.
func (g RMutexGuard[T]) RUnlock() {
	g.rw.RUnlock()
}

// RWMutex creates a type safe RWMutex wrapper ensuring one cannot access the
// values of a locked values without first gaining a lock.
type RWMutex[T any] struct {
	rw    *sync.RWMutex
	value T
}

// NewRWMutex creates a new RWMutex for use.
func NewRWMutex[T any](value T) RWMutex[T] {
	return RWMutex[T]{
		rw:    new(sync.RWMutex),
		value: value,
	}
}

// TryLock tries to lock RWMutex and returns the value in the Ok result if successful.
// If it does the value is returned for use in the Ok result otherwise Err with empty value.
func (m RWMutex[T]) TryLock() resultext.Result[MutexGuard[T, *sync.RWMutex], struct{}] {
	if m.rw.TryLock() {
		return resultext.Ok[MutexGuard[T, *sync.RWMutex], struct{}](
			MutexGuard[T, *sync.RWMutex]{
				m: m.rw,
				T: m.value,
			})
	} else {
		return resultext.Err[MutexGuard[T, *sync.RWMutex]](struct{}{})
	}
}

// Lock locks the Mutex and returns value for use,
// safe for mutation if the lock is already in use,
// the calling goroutine blocks until the mutex is available.
func (m RWMutex[T]) Lock() MutexGuard[T, *sync.RWMutex] {
	m.rw.Lock()
	return MutexGuard[T, *sync.RWMutex]{
		m: m.rw,
		T: m.value,
	}
}

// TryRLock tries to lock RWMutex for reading and returns the value in the Ok result if successful.
// If it does the value is returned for use in the Ok result otherwise Err with empty value.
func (m RWMutex[T]) TryRLock() resultext.Result[RMutexGuard[T], struct{}] {
	if m.rw.TryRLock() {
		return resultext.Ok[RMutexGuard[T], struct{}](
			RMutexGuard[T]{
				rw: m.rw,
				T:  m.value,
			},
		)
	} else {
		return resultext.Err[RMutexGuard[T]](struct{}{})
	}
}

// RLock locks the RWMutex for reading and returns the value for read-only use.
// It should not be used for recursive read locking,
// because a blocked Lock call excludes new readers from acquiring the lock.
func (m RWMutex[T]) RLock() RMutexGuard[T] {
	m.rw.RLock()
	return RMutexGuard[T]{
		rw: m.rw,
		T:  m.value,
	}
}

// Perform safely locks and unlocks the RWMutex read-only values and performs the provided function.
func (m RWMutex[T]) Perform(f func(T)) {
	guard := m.RLock()
	f(guard.T)
	guard.RUnlock()
}

// PerformMut safely locks and unlocks the RWMutex mutable values and performs the provided function.
func (m RWMutex[T]) PerformMut(f func(T)) {
	guard := m.Lock()
	f(guard.T)
	guard.Unlock()
}
