package fcontext

import (
	"github.com/posener/fcontext/contexttest"
	"testing"
)

func TestContext(t *testing.T) {
	imp := contexttest.Implementation {
		WithValue: WithValue,
		WithCancel: WithCancel,
	}
	imp.Run(t)
}