// r.go: repl

package snoc

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
	//"strings"

	. "github.com/strickyak/yak"
)

func PreprocessFunc(name string, params []*Sym, body Any, outer *ProtoFunc) (pf *ProtoFunc) {
	Log("PreprocessFunc: %q %v <<< %v <<< %v", name, params, body, outer)

	defer func() {
		r := recover()
		if r != nil {
			pf = nil
			debug.PrintStack()
			log.Panicf("ddt: PreprocessFunc error: %v", r)
		}
	}()

	pf = &ProtoFunc{
		Outer:  outer,
		Params: params,
		Body:   nil,
		Name:   name,
	}
	var preprocess func(a Any) Any
	preprocess = func(a Any) Any {
		switch t := a.(type) {
		case *Sym:
			Log("preprocess *Sym: <<< %v", t)
			for p := pf; p != nil; p = p.Outer {
				Log("preprocess *Sym: === p=%v", p)
				for i, prm := range p.Params {
					Log("preprocess *Sym: === [i=%d]prm=%v", i, prm)
					if t == prm {
						v := &Var{Proto: p, Slot: i, Sym: t}
						Log("preprocess *Sym: >>> CHANGED TO %v", v)
						return v
					}
				}
			}
			Log("preprocess *Sym: >>> SAME")
			return t // Default: dont change sym.
		case *Pair:
			if t == NIL {
				return NIL
			}
			log.Printf("ddt: case *Pair: H <<< %v >>> T <<< %v >>>", t.H, t.T)
			switch t.H {
			case FN:
				return PreprocessFunc(Serial("FN_"), ListToVecOfSym(t.T.H), t.T.T.H, pf)
			case Intern("let"):
				id2 := Serial("LET_")
				var params2 []*Sym
				var values2 []Any
				var body2 Any
				for p := t.T; p.T != NIL; p = p.T.T {
					params2 = append(params2, p.H.(*Sym))
					values2 = append(values2, p.T.H)
					body2 = p.T.T.H
				}

				// Now we have params2, values2, & body2.
				pf2 := &ProtoFunc{
					Outer:  pf,
					Params: params2,
					Values: make([]Any, len(values2)),
					Body:   nil,
					Name:   id2 + "_LET_",
					IsLet:  true,
				}

				for i, e := range values2 {
					// pf2.Values[i] = PreprocessFunc(id2+params2[i].S, params2, e, pf2)
					pf2.Values[i] = PreprocessFunc(id2+params2[i].S, nil, e, pf2)
				}
				// pf2.Body = PreprocessFunc(id2+"_RESULT_", params2, body2, pf2)
				pf2.Body = PreprocessFunc(id2+"_RESULT_", nil, body2, pf2)
				return Snoc(NIL, pf2)

			default:
				return &Pair{
					H: preprocess(t.H),
					T: preprocess(t.T).(*Pair),
				}
			}
		}
		return a
	}

	Log("PreprocessFunc: %q %v ==== body_in: %v", name, params, body)
	pf.Body = preprocess(body)
	Log("PreprocessFunc: %q %v ==== body_out: %v", name, params, pf.Body)
	Log("PreprocessFunc: %q %v >>> %v", name, params, pf)
	return pf
}

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

func TryReplEval(terp *Terp, xs []Any) (result Any, err interface{}) {
	defer func() {
		r := recover()
		if r != nil {
			err = recover()
			result = err
		}
	}()

	result = NIL
	env := &Env{
		Terp: terp,
	}
	for _, x := range xs {
		log.Printf("ddt: for")
		if p, ok := x.(*Pair); ok {
			if p.H == DEF {
				vec := ListToVec(p.T)
				MustEq(len(vec), 2)
				sym, ok := vec[0].(*Sym)
				if !ok {
					Throw(vec[0], "DEF needs symbol at first")
				}
				terp.Globals[sym] = vec[1]
				result = NIL
				continue
			} else if p.H == DEFUN {
				log.Printf("ddt: DEFUN")
				println(p.H, DEFUN)
				vec := ListToVec(p.T)
				log.Printf("ddt: vec: %#v", vec)
				MustEq(len(vec), 3)
				sym, ok := vec[0].(*Sym)
				log.Printf("ddt: sym: %v %v", sym, ok)
				if !ok {
					Throw(vec[0], "DEFUN needs symbol at first")
				}
				// func PreprocessFunc(name string, params []*Sym, body Any, outer *ProtoFunc) *ProtoFunc
				proto := PreprocessFunc(sym.S, ListToVecOfSym(vec[1]), vec[2], nil)
				log.Printf("ddt: DEFUN sym %v proto %v", sym, proto)
				// defun := Snoc(Snoc(Snoc(NIL, vec[2]), vec[1]), FN)
				terp.Globals[sym] = proto
				result = NIL
				continue
			}
		}
		result = Eval(x, env)
	}
	return result, nil
}

func Repl(terp *Terp, r io.Reader) []Any {
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
		result, err := TryReplEval(terp, xs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			errStr := fmt.Sprintf("*ERROR* %v", err)
			// results = append(results, Snoc(Snoc(NIL, errStr), Intern("*ERROR*")))
			results = append(results, errStr)
		} else {
			fmt.Fprintf(os.Stderr, "---->   %v\n", result)
			results = append(results, result)
		}
	}
	return results
}
