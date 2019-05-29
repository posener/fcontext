// Package fcontext provides a fully compatible (pseudo) constant
// value access-time alternative to the standard library context
// package.
//
// The standard library context provides values access-time which
// is linear with the amount of values that are stored in the
// it. This implementation provides a constant access time in the
// most common context use case and linear access time for the
// less common use cases (This is why the term 'pseudo' is used).
// Please see the benchmarks below for details.
// Other parts of the context implementation left untouched.
//
// Concepts
//
// The main assumption that is made in this implementation is that
// context values tree is mostly grows tall and barely grows wide.
// This means that the way that the context will mostly be used is
// by adding more values to the existing context:
//
// 	ctx = context.WithValue(ctx, 1, 1)
// 	ctx = context.WithValue(ctx, 2, 2)
// 	ctx = context.WithValue(ctx, 3, 3)
//
// And not creating new branches of the existing context:
//
// 	ctx1 := context.WithValue(ctx, 1, 1)
// 	ctx2 := context.WithValue(ctx, 2, 2)
// 	ctx3 := context.WithValue(ctx, 3, 3)
//
// The last form might be more familiar in the form:
//
//	func main() {
// 		ctx := context.WithValue(context.Background(), 2, 2)
// 		f(ctx)
//		f(ctx)
// 	}
//
//	func f(ctx context.Context) {
// 		ctx = context.WithValue(2, 2)
// 	}
//
// This implementation will work either way, but will improve the
// performance of the first pattern significantly.
//
// Benchmarks
//
// Run the benchmarks with `make bench`. Results (On personal machine):
//
// **Access**: Constant access time regardless to the number of stored
// values. Compared to the standard library, on the average case, it
// performs 40% better for 10 values, 9 times better for 100 values
// and 71 times better for 1000 values.
//
// **Store**: About 5 times slower and takes about 5 more memory than
// the standard library context. (Can take up to 10 times if the
// context is only grown shallowly).
//
// Usage
//
// This library is fully compatible with the standard library context.
//
// 	 import (
// 	-	"context"
// 	+ 	context "github.com/posener/fcontext"
// 	 )
package fcontext

import (
	"sync"
)

// node implements the Context interface.
// A node can be created from any context, or from another node.
// Only one node can be created from another node - when the tree grows
// tall, once that was done, the fork flag is set to the parent node.
// A node that was created from another node will have an increased rank,
// and they will share the values and Context fields.
// If a child is needed to be created from a node that has the fork
// flag turned on, it will be created as a new node with rank 0 and an
// empty map, and will point to the parent through the Context field.
// Values in the values map are related to a specific node, and have
// an identical rank. Nodes can only access values in the map with a
// lower or equal rank (values that were known when the node was
// created).
type node struct {
	// Context is stored For non-values operations, breadth grow,
	// and conversion from other Context implementations.
	Context
	// values stores data shared by all nodes for fast lookup. The key
	// is the value key given to the WithValue function. The value is
	// a list of all assigned values to this key, with additional access
	// information for each key.
	values map[interface{}][]value
	// rank is the the depth of the current context from the root context.
	rank int
	// hasChild indicates if a child was already created from this node.
	hasChild bool
	// mu to allow concurrent usage.
	mu *sync.RWMutex
}

// Holds value's data and the rank of the context that this data was created in.
type value struct {
	data interface{}
	rank int
}

// WithValue returns a copy of parent in which the value associated with
// key is val. The exposed behavior is identical to the standard library
// function: https://golang.org/pkg/context/#WithValue.
func WithValue(ctx Context, key, val interface{}) Context {
	child := newChild(ctx)
	child.mu.Lock()
	defer child.mu.Unlock()
	child.put(key, val)
	return child
}

// WithValues is similar to WithValue, but for multiple key-values paris.
// The pairs argument and must be of even number of elements.
// It is more efficient to use this function than call the WithValue
// function multiple times.
func WithValues(ctx Context, pairs ...interface{}) Context {
	if len(pairs) == 0 {
		return ctx
	}
	if len(pairs)%2 != 0 {
		panic("pairs length must be even")
	}
	child := newChild(ctx)
	child.mu.Lock()
	defer child.mu.Unlock()
	for i := 0; i < len(pairs); i += 2 {
		child.put(pairs[i], pairs[i+1])
	}
	return child
}

func (n *node) put(key, val interface{}) {
	n.values[key] = append(n.values[key], value{data: val, rank: n.rank})
}

// newChild prepares a new child for a given parent context.
func newChild(ctx Context) *node {
	parent, ok := ctx.(*node)
	// Convert a non-node context.
	if !ok {
		return new(ctx)
	}

	parent.mu.Lock()
	defer parent.mu.Unlock()

	// Nodes are allowed to have only one child, return an empty new node
	// for the given context.
	if parent.hasChild {
		return new(ctx)
	}

	// Create a new child from this parent node.
	child := *parent
	child.rank++

	parent.hasChild = true

	return &child
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key.
// For more info see the interface documentation:
// https://golang.org/pkg/context/#Context
func (n *node) Value(key interface{}) interface{} {
	if val := n.lookupValues(key); val != nil {
		return val
	}
	// Fallback to lookup in parent context.
	return n.Context.Value(key)
}

func (n *node) lookupValues(key interface{}) interface{} {
	n.mu.RLock()
	defer n.mu.RUnlock()
	values := n.values[key]
	// Iterate the values in a reverse order to get the
	// most updated value that matches n's rank. (Most recent
	// values are added last).
	for i := len(values) - 1; i >= 0; i-- {
		v := values[i]
		// n is allowed to see v only if its rank is greater or equal
		// to v's rank.
		if v.rank <= n.rank {
			return v.data
		}
	}
	return nil
}

// new converts any context to a new node with a new values map.
func new(ctx Context) *node {
	return &node{
		Context: ctx,
		values:  make(map[interface{}][]value),
		mu:      &sync.RWMutex{},
	}
}
