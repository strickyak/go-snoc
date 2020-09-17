// +build main

package main

import (
	"flag"
	"os"

	. "github.com/strickyak/go-snoc"
	. "github.com/strickyak/yak"
)

func main() {
	flag.Parse()

	results, _ := Repl(NewEnv(), os.Stdin)
	for i, result := range results {
		L("==> result[%d] = %v", i, result)
	}
}
