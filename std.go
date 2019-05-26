package fcontext

import "context"

type Context = context.Context

var (
	Background  = context.Background
	TODO        = context.TODO
	WithCancel  = context.WithCancel
	WithTimeout = context.WithTimeout
)