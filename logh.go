package logh

import (
	"bufio"
	"fmt"
	"io"
	"regexp"

	"github.com/fatih/color"
)

var Colors = []color.Attribute{
	color.FgRed,
	color.FgGreen,
	color.FgBlue,
	color.FgYellow,
	color.FgCyan,
	color.FgHiRed,
	color.FgHiGreen,
	color.FgHiBlue,
	color.FgHiYellow,
	color.FgHiCyan,
}

func Highlight(in io.Reader, out io.Writer, matches ...string) {
	// Instantiates all the colors
	colors := make([]*color.Color, len(matches))
	for i := 0; i < len(matches); i++ {
		c := Colors[i%len(Colors)]
		colors[i] = color.New(c)
	}

	// Create a regexp for each matches
	regexes := make([]*regexp.Regexp, len(matches))
	for i, match := range matches {
		if match == "" { // if match is empty, we shouldn't really match with anything
			match = "^$"
		}

		rg, _ := regexp.Compile(match) // TODO handle error
		regexes[i] = rg
	}

	// Read each line from in
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()

		// See if matches with any regexes
		matched := false
		for i, rg := range regexes {
			if rg.MatchString(line) {
				// Matches! Highlight and output
				fmt.Fprintf(out, "%s\n", colors[i].Sprint(line))

				matched = true
				break
			}
		}

		// No matches, print with default color
		if !matched {
			fmt.Fprintf(out, "%s\n", scanner.Text())
		}
	}
}
