package fcontext

import (
	"context"
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

func TestValue_convertStandardContext(t *testing.T) {
	t.Parallel()

	ctx0 := context.WithValue(context.Background(), 0, 0)
	ctx1 := WithValue(ctx0, 0, 1)

	assert.Equal(t, 0, ctx0.Value(0))
	assert.Equal(t, 1, ctx1.Value(0))
}

func TestValue_byteOverflow(t *testing.T) {
	t.Parallel()

	ctx0 := Background()
	var ctxs []context.Context

	for i := 0; i < 1024; i++ {
		ctxs = append(ctxs, WithValue(ctx0, i, i))
	}

	for i, ctx := range ctxs {
		assert.Equal(t, i, ctx.Value(i))
	}
}

func TestWithValues(t *testing.T) {
	t.Parallel()

	ctx := WithValues(Background(), 0, 0, 1, 1)
	assert.Equal(t, 0, ctx.Value(0))
	assert.Equal(t, 1, ctx.Value(1))
}

func TestWithValues_panic(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		WithValues(Background(), 0)
	})
}
