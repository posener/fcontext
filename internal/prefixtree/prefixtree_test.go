package prefixtree

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasics(t *testing.T) {
	t.Parallel()

	n := New()

	assert.True(t, n.HasAncestor(n))
	assert.True(t, n.NewChild().HasAncestor(n))
	assert.False(t, n.HasAncestor(n.NewChild()))
}

func TestGraph(t *testing.T) {
	t.Parallel()

	n := New()
	n1 := n.NewChild()
	n11 := n1.NewChild()
	n12 := n1.NewChild()
	n2 := n.NewChild()

	tests := []struct {
		// Two nodes.
		parent, child *Node
		// Whether parent is actually ancestor of child.
		ancestor bool
	}{
		{n, n, true},
		{n, n1, true},
		{n, n2, true},
		{n, n11, true},
		{n, n12, true},

		{n1, n, false},
		{n2, n, false},
		{n11, n, false},
		{n12, n, false},

		{n1, n1, true},
		{n1, n11, true},
		{n1, n12, true},

		{n11, n1, false},
		{n12, n1, false},

		{n1, n2, false},
		{n11, n2, false},
		{n12, n2, false},

		{n2, n1, false},
		{n2, n11, false},
		{n2, n12, false},

		{n11, n12, false},
		{n12, n11, false},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.ancestor, tt.child.HasAncestor(tt.parent))
	}
}

func TestDepth(t *testing.T) {
	t.Parallel()
	const size = idSize*8 + 1

	nodes := []*Node{New()}

	for i := 0; i < size; i++ {
		nodes = append(nodes, nodes[len(nodes)-1].NewChild())
	}

	// Make sure that we at least passed one rank.
	require.True(t, nodes[len(nodes)-1].rank > 0)

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			assert.Equalf(t, j <= i, nodes[i].HasAncestor(nodes[j]), "nodes[%d].HasAncestor(nodes[%d])", i, j)
		}
	}
}
