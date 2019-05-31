package fcontext

import "context"

type Context = context.Context
type CancelFunc = context.CancelFunc

var (
	Background       = context.Background
	TODO             = context.TODO
	WithTimeout      = context.WithTimeout
	Canceled         = context.Canceled
	WithDeadline     = context.WithDeadline
	DeadlineExceeded = context.DeadlineExceeded
)
