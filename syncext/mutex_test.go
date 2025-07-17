package syncext

import (
	"sync"
	"testing"

	"github.com/pchchv/express/resultext"
	. "github.com/pchchv/go-assert"
)

func TestMutex(t *testing.T) {
	m := NewMutex(make(map[string]int))
	guard := m.Lock()
	guard.T["foo"] = 1
	guard.Unlock()
	m.PerformMut(func(m map[string]int) {
		m["boo"] = 1
	})
	guard = m.Lock()
	myMap := guard.T
	Equal(t, 2, len(myMap))
	Equal(t, myMap["foo"], 1)
	Equal(t, myMap["boo"], 1)
	Equal(t, m.TryLock(), resultext.Err[MutexGuard[map[string]int, *sync.Mutex]](struct{}{}))
	guard.Unlock()

	result := m.TryLock()
	Equal(t, result.IsOk(), true)
	result.Unwrap().Unlock()
}

func TestRWMutex(t *testing.T) {
	m := NewRWMutex(make(map[string]int))
	guard := m.Lock()
	guard.T["foo"] = 1
	Equal(t, m.TryLock().IsOk(), false)
	Equal(t, m.TryRLock().IsOk(), false)
	guard.Unlock()

	m.PerformMut(func(m map[string]int) {
		m["boo"] = 2
	})
	guard = m.Lock()
	mp := guard.T
	Equal(t, mp["foo"], 1)
	Equal(t, mp["boo"], 2)
	guard.Unlock()

	rguard := m.RLock()
	myMap := rguard.T
	Equal(t, len(myMap), 2)
	Equal(t, m.TryRLock().IsOk(), true)
	rguard.RUnlock()

	m.Perform(func(m map[string]int) {
		Equal(t, 1, m["foo"])
		Equal(t, 2, m["boo"])
	})
	rguard = m.RLock()
	myMap = rguard.T
	Equal(t, len(myMap), 2)
	Equal(t, myMap["foo"], 1)
	Equal(t, myMap["boo"], 2)
	rguard.RUnlock()
}
