// +build main

package main

import (
	"os"

	. "github.com/strickyak/go-snoc"
	. "github.com/strickyak/yak"
)

func main() {
	results, _ := Repl(NewEnv(), os.Stdin)
	for i, result := range results {
		L("==> result[%d] = %v", i, result)
	}
}
