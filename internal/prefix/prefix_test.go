package prefix

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNil(t *testing.T) {
	t.Parallel()

	assert.True(t, New().HasPrefix(New()))
	assert.True(t, New().Append('a').HasPrefix(New()))
	assert.False(t, New().HasPrefix(New().Append('a')))
}

func TestBasic(t *testing.T) {
	t.Parallel()

	a := New().Append('a')
	b := New().Append('b')
	ab := a.Append('b')
	aa := a.Append('a')

	assert.True(t, ab.HasPrefix(a))
	assert.False(t, ab.HasPrefix(b))
	assert.False(t, a.HasPrefix(ab))
	assert.False(t, a.HasPrefix(aa))
}

func TestAppendMoreThanSize(t *testing.T) {
	t.Parallel()

	p := New()
	for i := 0; i < size+1; i++ {
		p = p.Append('a')
	}

	assert.True(t, p.HasPrefix(p))
	assert.True(t, p.Append('a').HasPrefix(p))
	assert.False(t, p.HasPrefix(p.Append('a')))

	p2 := p
	for i := 0; i < size+1; i++ {
		p2 = p2.Append('a')
	}
	assert.True(t, p2.HasPrefix(p))
	assert.False(t, p.HasPrefix(p2))
}

func TestFirstPrefixDiffer(t *testing.T) {
	t.Parallel()

	p1 := New()
	for i := 0; i < size; i++ {
		p1 = p1.Append('a')
	}

	p2 := New()
	for i := 0; i < size; i++ {
		p2 = p2.Append('b')
	}

	// Make sure we still in the first level.
	require.Equal(t, 0, p1.rank)
	require.Equal(t, 0, p2.rank)

	assert.False(t, p1.HasPrefix(p2))

	p1 = p1.Append('a')
	p2 = p2.Append('a')

	// Make sure we passed to the next level.
	require.Equal(t, 1, p1.rank)
	require.Equal(t, 1, p2.rank)

	assert.False(t, p1.HasPrefix(p2))
}
