package main

import (
	"os"

	"github.com/haroflow/logh"
)

func main() {
	matches := os.Args[1:]
	logh.Highlight(os.Stdin, os.Stdout, matches...)
}
