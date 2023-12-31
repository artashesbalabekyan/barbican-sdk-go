package client

import "context"

type Iterator interface {
	Next() (string, bool)
	Close() error
}

type ListIterator struct {
	ch     <-chan string
	ctx    context.Context
	cancel context.CancelCauseFunc
}

// Next moves the iterator to the next key, if any.
// This key is available until Next is called again.
//
// It returns true if and only if there is a new key
// available. If there are no more keys or an error
// has been encountered, Next returns false.
func (i *ListIterator) Next() (string, bool) {
	select {
	case v, ok := <-i.ch:
		return v, ok
	case <-i.ctx.Done():
		return "", false
	}
}

// Err returns the first error, if any, encountered
// while iterating over the set of keys.
func (i *ListIterator) Close() error {
	// i.cancel(context.Canceled)
	return context.Cause(i.ctx)
}
