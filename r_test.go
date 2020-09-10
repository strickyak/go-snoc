package snoc

import (
	"strings"
	"testing"

	. "github.com/strickyak/yak"
)

func TestSignum(t *testing.T) {
	r := strings.NewReader(`
    (def pos 1)
    (def neg -1)
    (def zero 0)
    (defun signum (x) (if
      (< x 0) neg
      (> x 0) pos
      zero))
    (signum -888)
    (signum 0)
    (signum 123)
  `)
	results, _ := Repl(NewEnv(), r)
	for i, result := range results {
		L("==> result[%d] = %v", i, result)
	}

	z := results[len(results)-3]
	if z.(*Float) == nil || z.(*Float).F != -1 {
		t.Errorf("Expected results[-3] to be -1, got %v", z)
	}

	z = results[len(results)-2]
	if z.(*Float) == nil || z.(*Float).F != 0 {
		t.Errorf("Expected results[-3] to be 0, got %v", z)
	}

	z = results[len(results)-1]
	if z.(*Float) == nil || z.(*Float).F != 1 {
		t.Errorf("Expected results[-3] to be 1, got %v", z)
	}
}

func TestTriangle(t *testing.T) {
	r := strings.NewReader(`
    (defun triangle (x) (
      if (< x 1)
         0
         (+ x (triangle (- x 1)))
    ))
    (triangle 6)
  `)
	results, _ := Repl(NewEnv(), r)
	for i, result := range results {
		L("==> result[%d] = %v", i, result)
	}

	z := results[len(results)-1]
	if z.(*Float) == nil || z.(*Float).F != 21 {
		t.Errorf("Expected (triangle 6) to be 21, got %v", z)
	}
}
