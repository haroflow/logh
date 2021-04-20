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

type HighlightConfig struct {
	IgnoreCase bool
}

func Highlight(in io.Reader, out io.Writer, config HighlightConfig, expressions ...string) error {
	// Instantiates all the colors
	colors := make([]*color.Color, len(expressions))
	// Create a regexp for each expression
	regexes := make([]*regexp.Regexp, len(expressions))
	for i, expr := range expressions {
		if expr == "" { // if expr is empty, we shouldn't really match with anything
			expr = "^$"
		}
		if config.IgnoreCase {
			expr = "(?i)" + expr
		}

		rg, err := regexp.Compile(expr)
		if err != nil {
			return fmt.Errorf("error compiling regular expression %s: %s", expr, err)
		}
		regexes[i] = rg

		c := Colors[i%len(Colors)]
		colors[i] = color.New(c)
	}

	// Read each line from in
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()

		// Maps a color for each char in line.
		// -1 is the default color
		// each matching regex will "paint" this slice with the color index
		lineColors := make([]int, len(line))
		for i := range lineColors {
			lineColors[i] = -1
		}

		// See if matches with any regexes
		matched := false
		for i, rg := range regexes {
			indexes := rg.FindAllStringIndex(line, -1)

			for _, match := range indexes {
				start, end := match[0], match[1]

				for j := start; j < end; j++ {
					lineColors[j] = i
				}

				matched = true
			}
		}

		if matched {
			buff := ""
			lastColorIdx := -1

			for i, ch := range line {
				if lineColors[i] == lastColorIdx {
					buff += string(ch)
					continue
				}

				if lastColorIdx == -1 {
					fmt.Fprintf(out, "%s", buff)
				} else {
					lastColor := colors[lastColorIdx]
					fmt.Fprintf(out, "%s", lastColor.Sprint(buff))
				}

				lastColorIdx = lineColors[i]
				buff = string(ch)
			}

			if lastColorIdx == -1 {
				fmt.Fprintf(out, "%s\n", buff)
			} else {
				lastColor := colors[lastColorIdx]
				fmt.Fprintf(out, "%s\n", lastColor.Sprint(buff))
			}
		} else {
			// No matches, print with default color
			fmt.Fprintf(out, "%s\n", scanner.Text())
		}
	}

	return nil
}
