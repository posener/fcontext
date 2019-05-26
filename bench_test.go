package fcontext_test

import (
	"context"
	"testing"

	"github.com/posener/fcontext"
)

type withValue func(context.Context, interface{}, interface{}) context.Context

func Benchmark(b *testing.B) {
	benchmarks := []struct{
		name string
		fn func(*testing.B, withValue)
	}{
		{"WithValue/Deep", withValueNested},
		{"WithValue/Shallow", withValueShallow},
		{"Value/DifferentKeys/Average", valueDifferentKeysAvg},
		{"Value/DifferentKeys/LastestValue", func(b *testing.B, wv withValue) { valueDifferentKeys(b, wv, 0) }},
		{"Value/DifferentKeys/EarliestValue", func(b *testing.B, wv withValue) { valueDifferentKeys(b, wv, values-1) }},
		{"Value/DifferentKeys/NotFound", func(b *testing.B, wv withValue) { valueDifferentKeys(b, wv, values) }},
		{"Value/SameKey", valueSameKey},
	}

	for _, bench := range benchmarks {
		b.Run(bench.name, func(b *testing.B) {
			b.Run("fcontext", func(b *testing.B) {bench.fn(b, fcontext.WithValue)})
			b.Run("stdctx", func(b *testing.B) {bench.fn(b, context.WithValue)})
		})
	}
}

const values = 1000

func withValueShallow(b *testing.B, withValue withValue) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		withValue(ctx, i, i)
	}
}

func withValueNested(b *testing.B, withValue withValue) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		ctx = withValue(ctx, i, i)
	}
}

func valueDifferentKeysAvg(b *testing.B, withValue withValue) {
	ctx := context.Background()
	for i := 0; i < values; i++ {
		ctx = withValue(ctx, i, i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < values; j++ {
		_ = ctx.Value(j)
		}
	}
}

func valueDifferentKeys(b *testing.B, withValue withValue, valueToGet int) {
	ctx := context.Background()
	for i := 0; i < values; i++ {
		ctx = withValue(ctx, i, i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx.Value(valueToGet)
	}
}

func valueSameKey(b *testing.B, withValue withValue) {
	ctx := context.Background()
	for i := 0; i < values; i++ {
		ctx = withValue(ctx, 0, 0)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx.Value(0)
	}
}
