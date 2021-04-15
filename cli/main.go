package main

// go build -o logh cli/main.go

import (
	"os"

	"github.com/haroflow/logh"
)

func main() {
	matches := os.Args[1:]
	logh.Highlight(os.Stdin, os.Stdout, matches...)
}
