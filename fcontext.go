// Package fcontext provides an alternative to the standard library
// context package.
//
// This implementation is fully compatible with the standard library
// context implementation and interface. 
// It compromises memory consumption and store time for improving
// the average access time for a values stored in the context object.
// Please see the benchmarks below for details.
// Other parts of the context implementation left untouched.
//
// Benchmarks
//
// Run the benchmarks with `make bench`.
// 
// * Store: About 5 times slower and takes about 5 more
//     memory than the standard library context. (Can take up to 10
//     times for an edge case usage usage).
// * Access: About 7 times faster on average access. Can be up to 23
//     times faster, or 5 times slower in the edge cases.
package fcontext

import "github.com/posener/fcontext/internal/prefixtree"

// node implements the Context interface.
type node struct {
	// Context is stored For non-values operations, fallback scenarios,
	// and conversion from other Context implementations.
	Context
	// Node in the tree. It is used to identify visibility of values.
	*prefixtree.Node
	// values store data shared by all nodes. The key is the value key,
	// and the value is a list of all assigned values to this key.
	// A node can access a value only if the node is an ancestor of the
	// value's node. This basically means that the value was assigned by an
	// ancestor of this node.
	values map[interface{}][]value
}

// Holds value's data and a pointer to a prefix that provides information
// who can access this data.
type value struct {
	data interface{}
	node *node
}

// WithValue returns a copy of parent in which the value associated with key is val.
// The exposed behavior is identical to the standard library function:
// https://golang.org/pkg/context/#WithValue.
func WithValue(ctx Context, key, val interface{}) Context {
	child := newChild(ctx)
	child.values[key] = append(child.values[key], value{data: val, node: child})
	return child
}

// WithValues is similar to WithValue, but for multiple values. The pairs is composed
// of key-value pairs and must be of even number of elements.
func WithValues(ctx Context, pairs ...interface{}) Context {
	if len(pairs) == 0 {
		return ctx
	}
	if len(pairs)%2 != 0 {
		panic("pairs length must be even")
	}
	child := newChild(ctx)
	for i := 0; i < len(pairs); i += 2 {
		key := pairs[i]
		val := pairs[i+1]
		child.values[key] = append(child.values[key], value{data: val, node: child})
	}
	return child
}

// newChild prepares a new child for a given parent context.
func newChild(ctx Context) *node {
	parent, ok := ctx.(*node)
	// Convert an existing implementation of context to node.
	if !ok {
		parent = new(ctx)
	}

	// Create a new prefix for the child, and increase the next id for
	// the parent.
	child := parent.NewChild()
	if child == nil {
		parent = new(ctx)
		child = parent.NewChild()
	}

	return &node{
		Context: parent.Context,
		values:  parent.values,
		Node:    child,
	}
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key.
// For more info see the interface documentation:
// https://golang.org/pkg/context/#Context
func (n *node) Value(key interface{}) interface{} {
	values := n.values[key]
	for i := len(values) - 1; i >= 0; i-- {
		v := values[i]
		if n.HasAncestor(v.node.Node) {
			return v.data
		}
	}
	// Fallback to lookup in parent context.
	return n.Context.Value(key)
}

// new converts any context to a new node.
func new(ctx Context) *node {
	return &node{
		Context: ctx,
		values:  make(map[interface{}][]value),
		Node:    prefixtree.New(),
	}
}