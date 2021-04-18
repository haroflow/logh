package main

// go build cli/logh.go
// set GOOS=linux; set GOARCH=amd64; go build cli/logh.go

import (
	"os"

	"github.com/haroflow/logh"
)

func main() {
	matches := os.Args[1:]
	logh.Highlight(os.Stdin, os.Stdout, matches...)
}
