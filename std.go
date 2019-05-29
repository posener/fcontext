package context

import stdcontext "context"

type Context = stdcontext.Context

var (
	Background  = stdcontext.Background
	TODO        = stdcontext.TODO
	WithCancel  = stdcontext.WithCancel
	WithTimeout = stdcontext.WithTimeout
)
