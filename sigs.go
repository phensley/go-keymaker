package keymaker

import "sync/atomic"

// State represents an initial state.
type State struct {
	v int32
}

// OK if state unchanged
func (c *State) OK() bool {
	return atomic.LoadInt32(&c.v) == 0
}

// Touch changes the state
func (c *State) Touch() {
	atomic.AddInt32(&c.v, 1)
}
