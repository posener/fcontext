package fcontext

import (
	"testing"

	"github.com/posener/fcontext/contexttest"
)

func TestContext(t *testing.T) {
	imp := contexttest.Implementation{
		WithValue:  WithValue,
		WithCancel: WithCancel,
	}
	imp.Run(t)
}
