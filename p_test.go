package snoc

import (
	"log"
	"testing"
)

func Test1(t *testing.T) {
	xs := ParseText("alpha ( beta gamma ) delta", "Test1")
	log.Printf("XS: %v", xs)
}
