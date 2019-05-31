package fcontext

import "sync"

type Cancelable interface {
	Context
	CancelWithErr(err error)
}

type Registerer interface {
	Register(Cancelable)
	Unregister(Cancelable)
}

type Child interface {
	Context
	Parent() Context
}

func WithCancel(ctx Context) (Context, CancelFunc) {
	registerer := lookupRegisterer(ctx)
	child := &cancelNode{
		Context:    ctx,
		registerer: registerer,
		children:   make(map[Cancelable]struct{}),
		done:       make(chan struct{}),
	}
	registerer.Register(child)
	return child, child.cancel
}

type cancelNode struct {
	Context
	done       chan struct{}
	registerer Registerer
	err        error
	children   map[Cancelable]struct{}
	mu         sync.RWMutex
}

func (n *cancelNode) Done() <-chan struct{} {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.done
}

func (n *cancelNode) Err() error {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.err
}

func (n *cancelNode) cancel() {
	n.CancelWithErr(Canceled)
	n.registerer.Unregister(n)
}

func (n *cancelNode) CancelWithErr(err error) {
	if err == nil {
		panic("context: internal error: missing cancel error")
	}

	n.mu.Lock()
	defer n.mu.Unlock()
	if n.err != nil {
		return // already canceled
	}
	n.err = err
	close(n.done)
	for child := range n.children {
		child.CancelWithErr(err)
	}
	n.children = nil
}

func (n *cancelNode) Register(child Cancelable) {
	if err := n.Err(); err != nil {
		// Already cancelled, cancel immediately.
		child.CancelWithErr(err)
		return
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	n.children[child] = struct{}{}
}

func (n *cancelNode) Unregister(child Cancelable) {
	n.mu.Lock()
	defer n.mu.Unlock()
	delete(n.children, child)
}

func lookupRegisterer(ctx Context) Registerer {
	for {
		switch c := ctx.(type) {
		case Registerer:
			return c
		case Child:
			ctx = c.Parent()
		default:
			return &wrapper{Context: ctx}
		}
	}
}

type wrapper struct {
	Context
}

func (w *wrapper) Register(child Cancelable) {
	go func() {
		select {
		case <-w.Done():
			child.CancelWithErr(w.Err())
		case <-child.Done():
		}
	}()
}
func (w *wrapper) Unregister(child Cancelable) {

}
