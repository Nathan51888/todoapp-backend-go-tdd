package interactions_test

import (
	"mytodoapp/domain/interactions"
	"mytodoapp/specifications"
	"testing"
)

func TestGreet(t *testing.T) {
	specifications.GreetSpecification(
		t,
		specifications.GreetAdapter(interactions.Greet),
	)
}
