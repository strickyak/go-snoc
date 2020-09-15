package snoc

import (
	"testing"
)

func Test1(t *testing.T) {
	xs := ParseText("alpha ( beta gamma ) delta", "Test1")
	t.Logf("XS: %v", xs)
}
