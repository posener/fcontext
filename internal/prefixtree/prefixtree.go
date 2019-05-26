// Package prefixtree is implementation of an N-tree. This tree
// provides a compromise between fast ancestor check and memory
// footprint.
// Its limitations are that each node can have up to 32 direct
// children.
package prefixtree

import "bytes"

const (
	idSize      = 32
	maxChildren = 32
)

// Node in the tree.
type Node struct {
	// id identifies the current node.
	id [idSize]byte
	// i is the position of this node in the current ID.
	i int
	// bi is index in the i'th byte of this node ID.
	bi uint
	// rank indicates the distance from the root (nil) node.
	// Starts from zero.
	rank int
	// parent points to a parent node.
	parent *Node
	// nextChildID remembers the ID of the next child to set.
	nextChildID byte
}

// New creates new root node.
func New() *Node {
	return &Node{}
}

// NewChild returns a new child for a given node. If the number
// of children exceeded the allowed limit, the returned node will
// be nil.
func (n *Node) NewChild() *Node {
	if n.nextChildID >= maxChildren {
		// TODO: maybe there is a nicer way to solve this.
		return nil
	}

	// A helper function that either copies the parent node to
	// be used as a child node, or create a new linked child node.
	copyNode := func() *Node {
		// If node id is full id size, create a new empty id node
		// and link it to this one as a parent.
		if n.i == idSize {
			return &Node{parent: n, rank: n.rank + 1}
		}
		// Copy current node.
		ret := *n
		return &ret
	}
	child := copyNode()

	// Set to the child, the current ID in the right location.
	child.id[child.i] |= n.nextChildID << child.bi
	n.nextChildID++

	// Increase child indices.
	child.bi++
	if child.bi == 8 {
		child.bi = 0
		child.i++
	}
	child.nextChildID = 0
	return child
}

// HasAncestor returns true of anc is n or an ancestor of n.
func (n *Node) HasAncestor(anc *Node) bool {
	// If anc is longer than n, it can't be a anc of n.
	if n.rank < anc.rank || (n.rank == anc.rank && (n.i < anc.i || (n.i == anc.i && n.bi < anc.bi))) {
		return false
	}
	for n.rank > anc.rank {
		n = n.parent
	}
	// Optimization for comparing the same ancestor.
	if n == anc {
		return true
	}

	// If the n and anc do not share the same parent, they for sure
	// different.
	if n.parent != anc.parent {
		return false
	}
	// Here n.rank == anc.rank.

	// Check if the first i-1 bytes match.
	if anc.i > 0 && !bytes.Equal(n.id[:anc.i-1], anc.id[:anc.i-1]) {
		return false
	}
	// Compare last byte, which might be partially filled.
	mask := byte((1 << anc.bi) - 1)
	return (n.id[anc.i] & mask) == (anc.id[anc.i] & mask)
}
