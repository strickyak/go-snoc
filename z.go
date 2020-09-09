// +build main

package main

import (
	"fmt"
	"strings"
	"text/scanner"
	//"unicode"
)

func main() {
	const src = `%var1 /* var2% */ "str" (def f(x-ray y'th z#^$#&*!^!^&!) (+/add/+ x y z))'(quote'd)'z  // final`

	var s scanner.Scanner
	s.Init(strings.NewReader(src))
	s.Filename = "default"
	s.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanChars | scanner.ScanStrings | scanner.ScanRawStrings | scanner.ScanComments | scanner.SkipComments

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		fmt.Printf("%s: %s\n", s.Position, s.TokenText())
	}

	fmt.Println()
	s.Init(strings.NewReader(src))
	s.Filename = "percent"
	s.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanStrings | scanner.ScanRawStrings | scanner.ScanComments | scanner.SkipComments

	// treat leading '%' as part of an identifier
	s.IsIdentRune = func(ch rune, i int) bool {
		return ch != '(' && ch != ')' && ch != '"' && (ch != '/' || i > 0) && ch > ' '

		// return ch == '%' && i == 0 || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0
	}

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		fmt.Printf("%s: %s\n", s.Position, s.TokenText())
	}
}
