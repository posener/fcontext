package contexttest

import (
	"sync"
	"time"
	"testing"
	"context"

	"github.com/stretchr/testify/assert"
)

type Implementation struct {
	WithCancel func(ctx context.Context) (context.Context, context.CancelFunc)
	WithValue func(ctx context.Context, key, val interface{}) context.Context
}

func (imp *Implementation) Run(t *testing.T) {
	t.Run("Value/Tree", imp.testValue_tree)
	t.Run("Value/Override", imp.testValue_valueOverride)
	t.Run("Value/Concurrency", imp.testValue_concurrency)
	t.Run("Cancel", imp.testCancel)
	t.Run("Cancel/Tree", imp.testCancel_tree)
	t.Run("Cancel/ParentAlreadyCanceled", imp.testCancel_parentAlreadyCanceled)
	t.Run("Cancel/Concurrency", imp.testCancel_concurrency)
}

func (imp *Implementation) testValue_tree(t *testing.T) {
	t.Parallel()

	ctx0 := context.Background()

	ctx01 := imp.WithValue(ctx0, 1, 1)
	assert.Equal(t, 1, ctx01.Value(1))
	assert.Equal(t, nil, ctx0.Value(1))

	ctx02 := imp.WithValue(ctx0, 2, 2)
	assert.Equal(t, nil, ctx02.Value(1))
	assert.Equal(t, 2, ctx02.Value(2))
	assert.Equal(t, nil, ctx0.Value(2))
	assert.Equal(t, 1, ctx01.Value(1))
	assert.Equal(t, nil, ctx01.Value(2))

	ctx021 := imp.WithValue(ctx02, 3, 3)
	assert.Equal(t, nil, ctx021.Value(1))
	assert.Equal(t, 2, ctx021.Value(2))
	assert.Equal(t, 3, ctx021.Value(3))
}

func (imp *Implementation) testValue_valueOverride(t *testing.T) {
	t.Parallel()

	ctx0 := imp.WithValue(context.Background(), 0, 0)
	ctx1 := imp.WithValue(ctx0, 0, 1)

	assert.Equal(t, 0, ctx0.Value(0))
	assert.Equal(t, 1, ctx1.Value(0))
}

func (imp *Implementation)testValue_concurrency(t *testing.T) {
	t.Parallel()

	const n = 1000

	ctx := imp.WithValue(context.Background(), 0, 0)

	var wg sync.WaitGroup
	wg.Add(3)

	go func(ctx context.Context) {
		for i := 0; i < n; i++ {
			ctx.Value(0)
		}
		wg.Done()
	}(ctx)

	go func(ctx context.Context) {
		for i := 0; i < n; i++ {
			ctx.Value(0)
		}
		wg.Done()
	}(ctx)

	go func(ctx context.Context) {
		for i := 0; i < n; i++ {
			imp.WithValue(ctx, i, i)
		}
		wg.Done()
	}(ctx)

	for i := 0; i < n; i++ {
		ctx = imp.WithValue(ctx, i, i)
	}

	wg.Wait()
}

func (imp *Implementation) testCancel(t *testing.T) {
	t.Parallel()
	ctx, cancel := imp.WithCancel(context.Background())
	assertNotCanceled(t, ctx)

	cancel()
	assertCanceled(t, ctx)
}

func (imp *Implementation) testCancel_tree(t *testing.T) {
	t.Parallel()

	ctx0, cancel0 := imp.WithCancel(context.Background())
	ctx00, cancel00 := imp.WithCancel(ctx0)
	ctx01, cancel01 := imp.WithCancel(ctx0)
	ctx000, cancel001 := imp.WithCancel(ctx00)

	cancel00()
	assertNotCanceled(t, ctx0)
	assertCanceled(t, ctx00)
	assertNotCanceled(t, ctx01)
	assertCanceled(t, ctx000)

	cancel0()
	assertCanceled(t, ctx0)
	assertCanceled(t, ctx00)
	assertCanceled(t, ctx01)
	assertCanceled(t, ctx000)

	cancel01()
	cancel001()
	assertCanceled(t, ctx0)
	assertCanceled(t, ctx00)
	assertCanceled(t, ctx01)
	assertCanceled(t, ctx000)
}

func  (imp *Implementation) testCancel_parentAlreadyCanceled(t *testing.T) {
	t.Parallel()

	ctx, cancel := imp.WithCancel(context.Background())
	cancel()
	assertCanceled(t, ctx)

	ctx, cancel = imp.WithCancel(ctx)
	assertCanceled(t, ctx)
	cancel()
}

func (imp *Implementation) testCancel_concurrency(t *testing.T) {
	t.Parallel()

	const n = 1000

	ctx, cancel := imp.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(3)

	go func(ctx context.Context) {
		for i := 0; i < n; i++ {
			ctx.Err()
		}
		wg.Done()
	}(ctx)

	go func(ctx context.Context) {
		for i := 0; i < n; i++ {
			ctx.Err()
		}
		wg.Done()
	}(ctx)

	go func(ctx context.Context) {
		for i := 0; i < n; i++ {
			imp.WithCancel(ctx)
		}
		wg.Done()
	}(ctx)

	for i := 0; i < n; i++ {
		ctx, cancel = imp.WithCancel(ctx)
		cancel()
	}

	wg.Wait()
}

func assertCanceled(t *testing.T, ctx context.Context) {
	assert.Equal(t, context.Canceled, ctx.Err())
	select {
	case <-ctx.Done():
	case <-time.After(100 * time.Millisecond):
		t.Error("context not done")
	}
}

func assertNotCanceled(t *testing.T, ctx context.Context) {
	assert.Nil(t, ctx.Err())
	select {
	case <-ctx.Done():
		t.Error("context done")
	case <-time.After(100 * time.Millisecond):
	}
}
