package prefix

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNil(t *testing.T) {
	t.Parallel()

	assert.True(t, New().HasPrefix(New()))
	assert.True(t, New().Append(1).HasPrefix(New()))
	assert.False(t, New().HasPrefix(New().Append(1)))
}

func TestAppend(t *testing.T) {
	a := New().Append(1)
	assert.Equal(t, byte(1), a.content[0])
	a = a.Append(1)
	assert.Equal(t, byte(3), a.content[0])
	a = a.Append(31)
	assert.Equal(t, byte(127), a.content[0])
}

func TestBasic(t *testing.T) {
	t.Parallel()

	a := New().Append(1)
	b := New().Append(2)
	ab := a.Append(2)
	aa := a.Append(1)

	assert.True(t, ab.HasPrefix(a))
	assert.False(t, ab.HasPrefix(b))
	assert.False(t, a.HasPrefix(ab))
	assert.False(t, a.HasPrefix(aa))
}

func TestAppendMoreThanSize(t *testing.T) {
	t.Parallel()

	p := New()
	for i := 0; i < chunkSize+1; i++ {
		p = p.Append(1)
	}

	assert.True(t, p.HasPrefix(p))
	assert.True(t, p.Append(1).HasPrefix(p))
	assert.False(t, p.HasPrefix(p.Append(1)))

	p2 := p
	for i := 0; i < chunkSize+1; i++ {
		p2 = p2.Append(1)
	}
	assert.True(t, p2.HasPrefix(p))
	assert.False(t, p.HasPrefix(p2))
}

func TestFirstPrefixDiffer(t *testing.T) {
	t.Parallel()

	p1 := New()
	for i := 0; i < chunkSize*8; i++ {
		p1 = p1.Append(1)
	}

	p2 := New()
	for i := 0; i < chunkSize*8; i++ {
		p2 = p2.Append(2)
	}

	// Make sure we still in the first level.
	require.Equal(t, 0, p1.rank)
	require.Equal(t, 0, p2.rank)

	assert.False(t, p1.HasPrefix(p2))

	p1 = p1.Append(1)
	p2 = p2.Append(1)

	// Make sure we passed to the next level.
	require.Equal(t, 1, p1.rank)
	require.Equal(t, 1, p2.rank)

	assert.False(t, p1.HasPrefix(p2))
}
