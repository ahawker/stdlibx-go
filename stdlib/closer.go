package stdlib

import (
	"io"
)

var _ io.Closer = CloserFn(nil)

// CloserFn is a function that can be used to close resources.
type CloserFn func() error

func (fn CloserFn) Close() error {
	return fn()
}

// NewCloserGroup creates a new *CloserGroup with sane defaults.
func NewCloserGroup(closers ...io.Closer) *CloserGroup {
	cg := &CloserGroup{
		Closers: make([]io.Closer, 0, len(closers)),
	}
	cg.Append(closers...)
	return cg
}

// CloserGroup is a collection of io.Closer instances that can be closed together.
type CloserGroup struct {
	Closers []io.Closer
}

// Append adds the given closers to the group.
func (g *CloserGroup) Append(closers ...io.Closer) {
	for _, closer := range closers {
		if closer == nil {
			continue
		}

		// When given a closer that's a group, we want to flatten & merge
		// the items.
		if cg, ok := closer.(*CloserGroup); ok {
			cg.Append(cg.Closers...)
			continue
		}

		g.Closers = append(g.Closers, closer)
	}
}

// Close closes all closers in the group.
func (g *CloserGroup) Close() error {
	eg := NewErrorGroup()
	for _, closer := range g.Closers {
		eg.Append(closer.Close())
	}
	return eg.ErrorOrNil()
}

// CloserJoin is a helper function that will append more closers
// onto an CloserGroup.
//
// If err is not already an ErrorGroup, then it will be turned into
// one. If any of the errs are ErrorGroup, they will be flattened
// one level into err.
// Any nil errors within errs will be ignored. If err is nil, a new
// *ErrorGroup will be returned containing the given errs.
func CloserJoin(closer io.Closer, closers ...io.Closer) *CloserGroup {
	if cg, ok := closer.(*CloserGroup); ok {
		cg.Append(closers...)
		return cg
	}
	cg := NewCloserGroup(closer)
	cg.Append(closers...)
	return cg
}
