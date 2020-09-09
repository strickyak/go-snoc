package snoc

import (
	"log"
	"testing"
)

func Test2(t *testing.T) {
	env := NewEnv()
	log.Printf("ENV: %v", env)

	xs := ParseText("10 ( + 10 90) (add 100 900)", "Test1")
	for i, x := range xs {
		log.Printf("X[%d]: (%T)%v", i, x, x)
		z := x.Eval(env)
		log.Printf("Z[%d]: (%T)%v", i, z, z)
	}
}
