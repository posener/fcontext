package fcontext

import (
	"context"
	"testing"
)

type provider struct {
	Background func() context.Context
	WithValue  func(context.Context, interface{}, interface{}) context.Context
}

var standard = provider{
	Background: context.Background,
	WithValue:  context.WithValue,
}

var fcontext = provider{
	Background: Background,
	WithValue:  WithValue,
}

func BenchmarkFContext_withValue_shallow(b *testing.B) {
	benchmarkWithValueShallow(b, fcontext)
}

func BenchmarkStdContext_withValue_shallow(b *testing.B) {
	benchmarkWithValueShallow(b, standard)
}

func BenchmarkFContext_withValue_nested(b *testing.B) {
	benchmarkWithValueNested(b, fcontext)
}

func BenchmarkStdContext_withValue_nested(b *testing.B) {
	benchmarkWithValueNested(b, standard)
}

func BenchmarkFContext_differentValues_first(b *testing.B) {
	benchmarkDifferentValues(b, fcontext, 0)
}
func BenchmarkStdContext_differentValues_first(b *testing.B) {
	benchmarkDifferentValues(b, standard, 0)
}

func BenchmarkFContext_differentValues_last(b *testing.B) {
	benchmarkDifferentValues(b, fcontext, values-1)
}
func BenchmarkStdContext_differentValues_last(b *testing.B) {
	benchmarkDifferentValues(b, standard, values-1)
}

func BenchmarkFContext_differentValues_notFound(b *testing.B) {
	benchmarkDifferentValues(b, fcontext, values)
}
func BenchmarkStdContext_differentValues_notFound(b *testing.B) {
	benchmarkDifferentValues(b, standard, values)
}

func BenchmarkFContext_sameValue(b *testing.B) {
	benchmarkSameValue(b, fcontext)
}

func BenchmarkStdContext_sameValue(b *testing.B) {
	benchmarkSameValue(b, standard)
}

const values = 100

func benchmarkWithValueShallow(b *testing.B, provider provider) {
	ctx := provider.Background()
	for i := 0; i < b.N; i++ {
		provider.WithValue(ctx, i, i)
	}
}

func benchmarkWithValueNested(b *testing.B, provider provider) {
	ctx := provider.Background()
	for i := 0; i < b.N; i++ {
		ctx = provider.WithValue(ctx, i, i)
	}
}

func benchmarkDifferentValues(b *testing.B, provider provider, valueToGet int) {
	ctx := provider.Background()
	for i := 0; i < values; i++ {
		ctx = provider.WithValue(ctx, i, i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx.Value(valueToGet)
	}
}

func benchmarkSameValue(b *testing.B, provider provider) {
	ctx := provider.Background()
	for i := 0; i < values; i++ {
		ctx = provider.WithValue(ctx, 0, 0)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx.Value(0)
	}
}
