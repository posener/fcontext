package fcontext

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCancel(t *testing.T) {
	t.Parallel()
	ctx, cancel := WithCancel(Background())
	assertNotCanceled(t, ctx)

	cancel()
	assertCanceled(t, ctx)
}

func TestCancelTree(t *testing.T) {
	t.Parallel()

	ctx0, cancel0 := WithCancel(Background())
	ctx00, cancel00 := WithCancel(ctx0)
	ctx01, cancel01 := WithCancel(ctx0)
	ctx000, cancel001 := WithCancel(ctx00)

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

func TestParentAlreadyCanceled(t *testing.T) {
	t.Parallel()

	ctx, cancel := WithCancel(Background())
	cancel()
	assertCanceled(t, ctx)

	ctx, cancel = WithCancel(ctx)
	assertCanceled(t, ctx)
	cancel()
}

func TestCancelConcurrency(t *testing.T) {
	t.Parallel()

	const n = 1000

	ctx, cancel := WithCancel(Background())

	var wg sync.WaitGroup
	wg.Add(3)

	go func(ctx Context) {
		for i := 0; i < n; i++ {
			ctx.Err()
		}
		wg.Done()
	}(ctx)

	go func(ctx Context) {
		for i := 0; i < n; i++ {
			ctx.Err()
		}
		wg.Done()
	}(ctx)

	go func(ctx Context) {
		for i := 0; i < n; i++ {
			WithCancel(ctx)
		}
		wg.Done()
	}(ctx)

	for i := 0; i < n; i++ {
		ctx, cancel = WithCancel(ctx)
		cancel()
	}

	wg.Wait()
}

func assertCanceled(t *testing.T, ctx Context) {
	assert.Equal(t, Canceled, ctx.Err())
	select {
	case <-ctx.Done():
	case <-time.After(100 * time.Millisecond):
		t.Error("context not done")
	}
}

func assertNotCanceled(t *testing.T, ctx Context) {
	assert.Nil(t, ctx.Err())
	select {
	case <-ctx.Done():
		t.Error("context done")
	case <-time.After(100 * time.Millisecond):
	}
}
