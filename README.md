# fcontext

[![Build Status](https://travis-ci.org/posener/fcontext.svg?branch=master)](https://travis-ci.org/posener/fcontext)
[![codecov](https://codecov.io/gh/posener/fcontext/branch/master/graph/badge.svg)](https://codecov.io/gh/posener/fcontext)
[![golangci](https://golangci.com/badges/github.com/posener/fcontext.svg)](https://golangci.com/r/github.com/posener/fcontext)
[![GoDoc](https://godoc.org/github.com/posener/fcontext?status.svg)](http://godoc.org/github.com/posener/fcontext)
[![goreadme](https://goreadme.herokuapp.com/badge/posener/fcontext.svg)](https://goreadme.herokuapp.com)

Package fcontext provides a fully compatible (pseudo) constant
value access-time alternative to the standard library context
package.

The standard library context provides values access-time which
is linear with the amount of values that are stored in the
it. This implementation provides a constant access time in the
most common context use case and linear access time for the
less common use cases (This is why the term 'pseudo' is used).
Please see the benchmarks below for details.
Other parts of the context implementation left untouched.

#### Concepts

The main assumption that is made in this implementation is that
context values tree is mostly grows tall and barely grows wide.
This means that the way that the context will mostly be used is
by adding more values to the existing context:

```go
ctx = context.WithValue(ctx, 1, 1)
ctx = context.WithValue(ctx, 2, 2)
ctx = context.WithValue(ctx, 3, 3)
```

And not creating new branches of the existing context:

```go
ctx1 := context.WithValue(ctx, 1, 1)
ctx2 := context.WithValue(ctx, 2, 2)
ctx3 := context.WithValue(ctx, 3, 3)
```

The last form might be more familiar in the following code:

```go
func main() {
	ctx := context.WithValue(context.Background(), 2, 2)
	f(ctx)
	f(ctx)
	// ...
}

func f(ctx context.Context) {
	ctx = context.WithValue(2, 2)
	// ...
}
```

This implementation will work either way, but will improve the
performance of the first pattern significantly.

#### Benchmarks

Run the benchmarks with `make bench`. Results (On personal machine):

**Access**: Constant access time regardless to the number of stored
values. Compared to the standard library, on the average case, it
performs 10%!(NOVERB)better for 10 values, 4 times better for 100 values
and 35 times better for 1000 values.

**Store**: About 6 times slower and takes about 4 more memory than
the standard library context. (Can take up to 8 times if the
context is only grown shallowly).

#### Usage

This library is fully compatible with the standard library context.

```diff
 import (
-	"context"
+ 	"github.com/posener/fcontext"
 )
```


---

Created by [goreadme](https://github.com/apps/goreadme)
