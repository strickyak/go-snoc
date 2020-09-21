# go-snoc
Simple LISP written in Go.

Currently uses a simple environment list for local variables,
no fancy lexical bindng.   NO I'M CHANGING THAT....

## work in progress...

Right now `go test` should work:

```
[0]<---- (def pos 1)
---->   nil
[0]<---- (def neg -1)
---->   nil
[0]<---- (def zero 0)
---->   nil
[0]<---- (defun signum (x) (if (< x 0) neg (> x 0) pos zero))
---->   nil
[0]<---- (list (signum -888) (signum 0) (signum 123))
---->   (-1 0 1)
[0]<---- (list (list 1 2 3) (list 4 5 6))
---->   ((1 2 3) (4 5 6))
[0]<---- (let A (list 1 2 3) B (list 4 5 6) C (list A B) (list A B C))
---->   ((1 2 3) (4 5 6) ((1 2 3) (4 5 6)))
[0]<---- (defun my-triangle (x) (if (< x 1) 0 (+ x (my-triangle (- x 1)))))
---->   nil
[0]<---- (my-triangle 6)
---->   21
[0]<---- (defun my-length (x) (if (null? x) 0 (+ 1 (my-length (tail x)))))
---->   nil
[0]<---- (my-length (list 9 7 5 3 1))
---->   5
[0]<---- (defun my-descending (n) (if (<= n 0) (list) (cons n (my-descending (- n 1)))))
---->   nil
[0]<---- (my-descending 7)
---->   (7 6 5 4 3 2 1)
[0]<---- (defun my-descending (n) (if (<= n 0) (list) (cons n (my-descending (- n 1)))))
---->   nil
[0]<---- (defun my-sum (aList) (if (null? aList) 0 (+ (head aList) (my-sum (tail aList)))))
---->   nil
[0]<---- 111
[1]<---- 222
[2]<---- 333
---->   333
[0]<---- (my-sum (my-descending 100))
---->   5050
```

You can also use snoc.go to evaluate stdin:

```
$ echo '(defun !(x) (if (< x 1) 1 (* x (! (- x 1))))) (! 10)' | go run snoc.go 
[0]<---- (defun ! (x) (if (< x 1) 1 (* x (! (- x 1)))))
[1]<---- (! 10)
---->   3.6288e+06
2020/09/13 17:33:18 ==> result[0] = 3.6288e+06

```
