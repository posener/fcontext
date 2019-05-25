package fcontext

import (
	"context"

	"github.com/posener/fcontext/internal/prefix"
)

// node implements the Context interface.
type node struct {
	// Context is stored For non-values operations, fallback scenarios,
	// and conversion from other Context implementations.
	context.Context
	// values store data shared by all nodes. The key is the value key,
	// and the value is a list of all assigned values to this key.
	// A node can access a value only if the node id has the prefix of this
	// value's id. This basically means that the value was assigned by an
	// ancestor of this node.
	values map[interface{}][]value
	// id for the node. A node has the ID that is equal to its parent ID
	// with an added unique byte. All children of a single node have unique
	// IDs.
	id *prefix.Prefix
	// nextID is the unique ID for the next child of this node.
	nextID byte
}

// Holds value's data and a pointer to a prefix that provides information
// who can access this data.
type value struct {
	data interface{}
	id   *prefix.Prefix
}

func Background() context.Context {
	return new(context.Background())
}

// WithValue returns a copy of parent in which the value associated with key is val.
// The exposed behavior is identical to the standard library function:
// https://golang.org/pkg/context/#WithValue.
func WithValue(ctx context.Context, key, val interface{}) context.Context {
	child := newChild(ctx)
	child.values[key] = append(child.values[key], value{data: val, id: child.id})
	return child
}

// WithValues is similar to WithValue, but for multiple values. The pairs is composed
// of key-value pairs and must be of even number of elements.
func WithValues(ctx context.Context, pairs ...interface{}) context.Context {
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
		child.values[key] = append(child.values[key], value{data: val, id: child.id})
	}
	return child
}

// newChild prepares a new child for a given parent context.
func newChild(ctx context.Context) *node {
	parent, ok := ctx.(*node)
	// Convert an existing implementation of context to node.
	if !ok {
		parent = new(ctx)
	}
	// If the parent have exceeded its number of allowed children, create
	// a parent with the fallback logic.
	if parent.nextID == byte(255) {
		parent = new(ctx)
	}

	// Create a new prefix for the child, and increase the next id for
	// the parent.
	id := parent.id.Append(parent.nextID)
	parent.nextID++

	return &node{
		Context: parent.Context,
		values:  parent.values,
		id:      id,
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
		if n.id.HasPrefix(v.id) {
			return v.data
		}
	}
	// Fallback to lookup in parent context.
	return n.Context.Value(key)
}

// new converts any context to a new node.
func new(ctx context.Context) *node {
	return &node{
		Context: ctx,
		values:  make(map[interface{}][]value),
	}
}
