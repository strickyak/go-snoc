// r.go: repl

package snoc

import (
	"bufio"
	"fmt"
	"io"
	"os"
	//"strings"

	. "github.com/strickyak/yak"
)

func TryReplParse(s string) (xs []Any, ok bool) {
	defer func() {
		r := recover()
		if r != nil {
			ok = false
		}
	}()
	ok = true
	xs = ParseText(s, "*repl*")
	return
}

func TryReplEval(env Env, xs []Any) (result Any, newenv Env, err interface{}) {
	defer func() {
		err = recover()
	}()
	result = NIL
	for _, x := range xs {
		if p, ok := x.(*Pair); ok {
			if p.H == DEF {
				vec := ListToVec(p.T)
				MustEq(len(vec), 2)
				sym, ok := vec[0].(*Sym)
				if !ok {
					Throw(vec[0], "DEF needs symbol at first")
				}
				result, env = NIL, env.SnocSnoc(vec[1], sym)
				continue
			} else if p.H == DEFUN {
				vec := ListToVec(p.T)
				MustEq(len(vec), 3)
				defun := Snoc(Snoc(Snoc(NIL, vec[2]), vec[1]), FN)
				sym, ok := vec[0].(*Sym)
				if !ok {
					Throw(vec[0], "DEFUN needs symbol at first")
				}
				result, env = NIL, env.SnocSnoc(defun, sym)
				continue
			}
		}
		result = Eval(x, env)
	}
	newenv = env
	return
}

func Repl(env Env, r io.Reader) ([]Any, Env) {
	sc := bufio.NewScanner(r)
	var results []Any
	buf := ""
	for sc.Scan() {
		//L("TEXT: %q", sc.Text())
		buf += sc.Text() + "\n"
		//L("TRY: %q", buf)
		xs, ok := TryReplParse(buf)
		if !ok {
			continue
		}
		buf = ""
		if len(xs) == 0 {
			continue
		}

		for i, x := range xs {
			fmt.Fprintf(os.Stderr, "[%d]<---- %v\n", i, x)
		}
		result, newenv, err := TryReplEval(env, xs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			errStr := fmt.Sprintf("%v", err)
			results = append(results, Snoc(Snoc(NIL, errStr), Intern("*ERROR*")))
		} else {
			fmt.Fprintf(os.Stderr, "---->   %v\n", result)
			env = newenv
			results = append(results, result)
		}
	}
	return results, env
}
