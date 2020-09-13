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

func TryReplParse(s string) (xs []X, ok bool) {
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

func TryReplEval(env Env, xs []X) (result X, newenv Env, err interface{}) {
	defer func() {
		err = recover()
	}()
	result = NIL
	for _, x := range xs {
		if p, ok := x.(*Pair); ok {
			if p.H == DEF {
				vec := ListToVec(p.T)
				MustEq(len(vec), 2)
				result, env = NIL, env.Snoc(vec[1]).Snoc(vec[0])
				continue
			} else if p.H == DEFUN {
				vec := ListToVec(p.T)
				MustEq(len(vec), 3)
				defun := NIL.Snoc(vec[2]).Snoc(vec[1]).Snoc(FN)
				result, env = NIL, env.Snoc(defun).Snoc(vec[0])
				continue
			}
		}
		result = x.Eval(env)
	}
	newenv = env
	return
}

func Repl(env Env, r io.Reader) ([]X, Env) {
	sc := bufio.NewScanner(r)
	var results []X
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
			results = append(results, NIL.Snoc(&Str{S: errStr}).Snoc(Intern("*ERROR*")))
		} else {
			fmt.Fprintf(os.Stderr, "---->   %v\n", result)
			env = newenv
			results = append(results, result)
		}
	}
	return results, env
}
