package snoc

import (
	"fmt"
	"github.com/strickyak/yak"
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

func ListToVec(a X) []X {
	var z []X
	for {
		p, ok := a.(*Pair)
		if !ok {
			yak.MustEq(a, NIL)
			break
		}
		z = append(z, p.H)
		a = p.T
	}
	return z
}

func VecToList(a []X) X {
	z := X(NIL)
	for i := len(a) - 1; i >= 0; i-- {
		z = &Pair{H: a[i], T: z}
	}
	return z
}

func ParseExprs(toks []Tok) (string, []Tok, []X) {
	//log.Printf("<<<<<<<<<<<<<<< %v", toks)
	var z []X
	last := ""
LOOP:
	for len(toks) > 0 {
		//log.Printf("XXXXXXXXXXXXXXX %v", toks)
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
				z = append(z, &Float{F: f})
			} else {
				z = append(z, Intern(t.Text))
			}
			toks = rest
		}
	}
	//log.Printf(">>>>>>>>>>>>>>> %v ::: %v ::: %v", z, last, toks)
	return last, toks, z
}

func ParseText(text, filename string) []X {
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

func ParseFile(filename string) []X {
	bb, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panicf("Cannot ReadFile %q: %v", filename, err)
	}
	return ParseText(string(bb), filename)
}
