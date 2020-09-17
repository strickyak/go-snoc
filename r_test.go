package snoc

import (
	"strings"
	"testing"
)

func TestPrograms(t *testing.T) {
	scenarios := []struct {
		program string
		want    string
	}{
		{`
			(list (list 1 2 3) (list 4 5 6))
		`, "((1 2 3) (4 5 6))"},

		{`
			(list (+ 3 4) (quote (+ 3 9)) (quote xyzzy))
		`, "(7 (+ 3 9) xyzzy)"},

		{`(let
			    A (list 1 2 3)
					B (list 4 5 6)
					C (list A B)
					(list A B C)
			)
		`, "((1 2 3) (4 5 6) ((1 2 3) (4 5 6)))"},

		{`
    (def pos 1)
    (def neg -1)
    (def zero 0)
    (defun signum (x) (if
      (< x 0) neg
      (> x 0) pos
      zero))
    (list (signum -888) (signum 0) (signum 123) )
    `, "(-1 0 1)"},

		{`
			(defun my-triangle (x) (
				if (< x 1)
					 0
					 (+ x (my-triangle (- x 1)))
			))
			(my-triangle 6)
		`, "21"},

		{`
			(defun my-length (x) (
				if (null? x)
					 0
					 (+ 1 (my-length (tail x)))
			))
			(my-length (list 9 7 5 3 1))
		`, "5"},

		{`
			(defun my-descending (n) (
				if (<= n 0)
					 (list)
					 (cons n (my-descending (- n 1)))
			))
			(my-descending 7)
		`, "(7 6 5 4 3 2 1)"},

		{`
			(defun my-descending (n) (
				if (<= n 0)
					 (list)
					 (cons n (my-descending (- n 1)))
			))
			(defun my-sum (aList) (
				if (null? aList)
					 0
					 (+ (head aList) (my-sum (tail aList)))
			))
			111 222 333
			(my-sum (my-descending 7))
		`, "28"},

		{`(let x (quote (+ 20 3))
		       (eval x))
		`, "23"},

		{`(let x (fn (aaa bbb) (- aaa bbb))
		       y (list 100 4)
		       (apply x y))
		`, "96"},
	}

	for j, sc := range scenarios {
		t.Logf("<== program[%d] = %v", j, sc.program)
		r := strings.NewReader(sc.program)

		results, _ := Repl(NewEnv(), r)
		for i, result := range results {
			t.Logf("==> result[%d.%d] = %v", j, i, result)
		}

		z := results[len(results)-1]
		got := Stringify(z)
		if got != sc.want {
			t.Errorf("Got %q, wanted %q, for program <<< %s >>>", got, sc.want, sc.program)
		}

	}
}
