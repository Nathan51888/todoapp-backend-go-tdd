package specifications

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Greeter interface {
	Greet() (string, error)
}

func GreetSpecification(t testing.TB, greeter Greeter) {
	got, err := greeter.Greet()
	assert.NoError(t, err)
	assert.Equal(t, "Hello, world", got)
}
