package main

// go build cli/logh.go
// set GOOS=linux; set GOARCH=amd64; go build cli/logh.go

import (
	"flag"
	"fmt"
	"os"

	"github.com/haroflow/logh"
)

func main() {
	ignoreCase := flag.Bool("i", false, "ignore case")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Usage:")
		fmt.Println("  cat file.txt | logh [flags] expr1 [expr2] [expr3]...")
		fmt.Println("")
		fmt.Println("Flags:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	config := logh.HighlightConfig{
		IgnoreCase: *ignoreCase,
	}
	matches := flag.Args()
	logh.Highlight(os.Stdin, os.Stdout, config, matches...)
}
