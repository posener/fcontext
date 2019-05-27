package fcontext_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/posener/fcontext"
)

type withValue func(context.Context, interface{}, interface{}) context.Context

var values = []int{
	10,
	100,
	1000,
}

func Benchmark(b *testing.B) {
	b.Run("fcontext", func(b *testing.B) { runBenchmarks(b, fcontext.WithValue) })
	b.Run("stdctx", func(b *testing.B) { runBenchmarks(b, context.WithValue) })
}

func runBenchmarks(b *testing.B, wv withValue) {
	benchmarks := []struct {
		name string
		fn   func(*testing.B, withValue)
	}{
		{"WithValue/Depth", withValueNested},
		{"WithValue/Breadth", withValueShallow},
		{"Value/DifferentKeys", valueDifferentKeys},
		{"Value/SameKey", valueSameKey},
	}
	for _, bench := range benchmarks {
		b.Run(bench.name, func(b *testing.B) {
			bench.fn(b, wv)
		})
	}
}

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

func valueDifferentKeys(b *testing.B, withValue withValue) {
	ctx := context.Background()
	for _, value := range values {
		b.Run(fmt.Sprintf("%dKeys", value), func(b *testing.B) {
			for i := 0; i < value; i++ {
				ctx = withValue(ctx, i, i)
			}

			b.Run("Average", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					for j := 0; j < value; j++ {
						_ = ctx.Value(j)
					}
				}
			})

			b.Run("Latest", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_ = ctx.Value(value - 1)
				}
			})

			b.Run("Earliest", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_ = ctx.Value(0)
				}
			})

			b.Run("NotFound", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_ = ctx.Value(value)
				}
			})
		})
	}
}

func valueSameKey(b *testing.B, withValue withValue) {
	ctx := context.Background()
	for _, value := range values {
		for i := 0; i < value; i++ {
			ctx = withValue(ctx, 0, 0)
		}
		b.Run(fmt.Sprintf("%dKeys", value), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = ctx.Value(0)
			}
		})
	}
}
