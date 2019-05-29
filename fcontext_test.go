package context

import (
	"sync"
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

	ctx := WithValue(Background(), 0, 0)

	var wg sync.WaitGroup
	wg.Add(3)

	go func(ctx Context) {
		for i := 0; i < n; i++ {
			ctx.Value(0)
		}
		wg.Done()
	}(ctx)

	go func(ctx Context) {
		for i := 0; i < n; i++ {
			ctx.Value(0)
		}
		wg.Done()
	}(ctx)

	go func(ctx Context) {
		for i := 0; i < n; i++ {
			WithValue(ctx, i, i)
		}
		wg.Done()
	}(ctx)

	for i := 0; i < n; i++ {
		ctx = WithValue(ctx, i, i)
	}

	wg.Wait()
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
