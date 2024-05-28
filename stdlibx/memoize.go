package stdlibx

import "sync"

// Memoize is a simple struct that wraps a `sync.Once`
// to provide thread safe memoized results for costly computation.
type Memoize[T any] struct {
	// Fn is called once and its result/error is cached.
	Fn func() (T, error)
	// res stores result of the computation.
	res T
	// err stores error if the computation failed.
	err error
	// once is used to ensure computation is only performed one time.
	once sync.Once
}

// Get returns the value + error from the
func (m *Memoize[T]) Get() (T, error) {
	m.once.Do(func() {
		m.res, m.err = m.Fn()
	})
	return m.res, m.err
}
