package snoc

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"text/scanner"
)

type Tok struct {
	Pos  scanner.Position
	Text string
}

func Lex(text, filename string) (z []Tok) {
	var s scanner.Scanner
	s.Init(strings.NewReader(text))
	s.Filename = filename
	s.Mode = scanner.ScanIdents | scanner.ScanStrings | scanner.ScanRawStrings | scanner.ScanComments | scanner.SkipComments
	s.IsIdentRune = func(ch rune, i int) bool {
		return ch != '(' && ch != ')' && (ch != '/' || i > 0) && ch > ' '
	}
	s.Error = func(_ *scanner.Scanner, msg string) {
		log.Printf("Lex error at %v: %s", s.Position, msg)
	}
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		//log.Printf("%s: %s\n", s.Position, s.TokenText())
		z = append(z, Tok{s.Position, s.TokenText()})
	}
	if s.ErrorCount > 0 {
		log.Panicf("Lex found %d errors in %q", s.ErrorCount, filename)
	}
	return
}

func ListToVec(a Any) []Any {
	b, ok := a.(*Pair)
	if !ok {
		Throw(a, "ListToVec expected *Pair")
	}
	var z []Any
	for p := b; p != NIL; p = p.T {
		z = append(z, p.H)
	}
	return z
}

func VecToList(a []Any) Any {
	z := NIL
	for i := len(a) - 1; i >= 0; i-- {
		z = &Pair{H: a[i], T: z}
	}
	return z
}

func ParseExprs(toks []Tok) (string, []Tok, []Any) {
	var z []Any
	last := ""
LOOP:
	for len(toks) > 0 {
		t, rest := toks[0], toks[1:]
		switch t.Text {
		case "(":
			last2, rest2, vec := ParseExprs(rest)
			if last2 != ")" {
				panic(fmt.Errorf("Parens not terminated: last=%q rest=%v", last2, rest2))
			}
			toks = rest2
			z = append(z, VecToList(vec))
		case ")":
			toks = rest
			last = ")"
			break LOOP
		default:
			f, err := strconv.ParseFloat(t.Text, 64)
			if err == nil {
				z = append(z, f)
			} else if t.Text == "nil" {
				z = append(z, NIL)
			} else {
				z = append(z, Intern(t.Text))
			}
			toks = rest
		}
	}
	return last, toks, z
}

func ParseText(text, filename string) []Any {
	toks := Lex(text, filename)
	last, rest, xs := ParseExprs(toks)
	if last != "" {
		log.Panicf("Did not expect nonempty last: last=%q rest=%v", last, rest)
	}
	if len(rest) > 0 {
		log.Panicf("Unused tokens from %q: %v", filename, rest)
	}
	return xs
}

func ParseFile(filename string) []Any {
	bb, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panicf("Cannot ReadFile %q: %v", filename, err)
	}
	return ParseText(string(bb), filename)
}
