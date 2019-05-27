package fcontext

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue_parentingGraph(t *testing.T) {
	t.Parallel()

	ctx0 := Background()
	assert.Equal(t, nil, ctx0.Value(1))

	ctx01 := WithValue(ctx0, 1, 1)
	assert.Equal(t, 1, ctx01.Value(1))
	assert.Equal(t, nil, ctx0.Value(1))

	ctx02 := WithValue(ctx0, 2, 2)
	assert.Equal(t, nil, ctx02.Value(1))
	assert.Equal(t, 2, ctx02.Value(2))
	assert.Equal(t, nil, ctx0.Value(2))
	assert.Equal(t, 1, ctx01.Value(1))
	assert.Equal(t, nil, ctx01.Value(2))

	ctx021 := WithValue(ctx02, 3, 3)
	assert.Equal(t, nil, ctx021.Value(1))
	assert.Equal(t, 2, ctx021.Value(2))
	assert.Equal(t, 3, ctx021.Value(3))
}

func TestValue_valueOverride(t *testing.T) {
	t.Parallel()

	ctx0 := WithValue(Background(), 0, 0)
	ctx1 := WithValue(ctx0, 0, 1)

	assert.Equal(t, 0, ctx0.Value(0))
	assert.Equal(t, 1, ctx1.Value(0))
}

func TestConcurrency(t *testing.T) {
	t.Parallel()

	const n = 1000

	ctx1 := WithValue(Background(), 0, 0)
	ctx2 := WithValue(ctx1, 1, 1)

	go func() {
		for i := 0; i < n; i++ {
			ctx1.Value(0)
		}
	}()

	for i := 0; i < n; i++ {
		ctx2 = WithValue(ctx2, i, i)
	}
}
