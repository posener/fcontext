# fcontext

[![Build Status](https://travis-ci.org/posener/fcontext.svg?branch=master)](https://travis-ci.org/posener/fcontext)
[![GoDoc](https://godoc.org/github.com/posener/fcontext?status.svg)](http://godoc.org/github.com/posener/fcontext)
[![goreadme](https://goreadme.herokuapp.com/badge/posener/fcontext.svg)](https://goreadme.herokuapp.com)

Package fcontext provides a fully compatible (pseudo) constant
value access-time alternative to the standard library context
package.

The standard library context provides values access-time which
is linear with the amount of values that are stored in the
it. This implementation provides a constant access time in the
most common context use case (This is why the term 'pseudo' is
used). Please see the benchmarks below for details.
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

This implementation will work either way, but will improve the
performance of the first pattern significantly.

#### Benchmarks

Run the benchmarks with `make bench`. Results (On personal machine):

**Access**: Constant access time regardless to the number of stored
values. Compared to the standard library, on the average case, it
performs 40%!(NOVERB)better for 10 values, 9 times better for 100 values
and 71 times better for 1000 values.

**Store**: About 5 times slower and takes about 5 more memory than
the standard library context. (Can take up to 10 times if the
context is only grown shallowly).


---

Created by [goreadme](https://github.com/apps/goreadme)
