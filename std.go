package fcontext

import "context"

type Context = context.Context
type CancelFunc = context.CancelFunc

var (
	Background  = context.Background
	TODO        = context.TODO
	WithCancel  = context.WithCancel
	WithTimeout = context.WithTimeout
)
