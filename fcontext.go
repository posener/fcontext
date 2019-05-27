// Package fcontext provides a fully compatible (pseudo) constant
// value access-time alternative to the standard library context
// package.
//
// The standard library context provides values access-time which
// is linear with the amount of values that are stored in the
// it. This implementation provides a constant access time in the
// most common context use case (This is why the term 'pseudo' is
// used). Please see the benchmarks below for details.
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
package fcontext

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
	// fork indicates if a child was already created from this node.
	fork bool
}

// Holds value's data and the rank of the context that this data was created in.
type value struct {
	data interface{}
	rank int
}

// WithValue returns a copy of parent in which the value associated with key is val.
// The exposed behavior is identical to the standard library function:
// https://golang.org/pkg/context/#WithValue.
func WithValue(ctx Context, key, val interface{}) Context {
	child := newChild(ctx)
	child.values[key] = append(child.values[key], value{data: val, rank: child.rank})
	return child
}

// WithValues is similar to WithValue, but for multiple key-value paris. The pairs
// argument must be of even number of elements.
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
		child.values[key] = append(child.values[key], value{data: val, rank: child.rank})
	}
	return child
}

// newChild prepares a new child for a given parent context.
func newChild(ctx Context) *node {
	parent, ok := ctx.(*node)
	// If the given context is not a node, or if the given context is
	// a node that already had a child, return a new node with new map
	// fro the given context.
	if !ok || parent.fork {
		return new(ctx)
	}
	// Set the fork flag to the parent since it is no longer allowed to
	// have children that share the same values map.
	parent.fork = true

	// Return the node first child that share the same Context and
	// values map.
	return &node{
		Context: parent.Context,
		values:  parent.values,
		rank:    parent.rank + 1,
	}
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key.
// For more info see the interface documentation:
// https://golang.org/pkg/context/#Context
func (n *node) Value(key interface{}) interface{} {
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
	// Fallback to lookup in parent context.
	return n.Context.Value(key)
}

// new converts any context to a new node with a new values map.
func new(ctx Context) *node {
	return &node{
		Context: ctx,
		values:  make(map[interface{}][]value),
	}
}
