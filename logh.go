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

		// TODO refactor this new implementation
		lineColors := make([]int, len(line))
		for i := range lineColors {
			lineColors[i] = -1
		}

		for i, rg := range regexes {
			indexes := rg.FindAllStringIndex(line, -1)

			for _, match := range indexes {
				start, end := match[0], match[1]
				for j := start; j < end; j++ {
					lineColors[j] = i
				}

				matched = true
			}

			// for _, match := range indexes {
			// 	start, end := match[0], match[1]
			// 	fmt.Fprintf(out, "%s%s%s\n", line[:start], colors[i].Sprint(line[start:end]), line[end:])

			// 	matched = true
			// 	break
			// }

			// TODO: Add config to highlight the whole line?
			// if rg.MatchString(line) {
			// 	// Matches! Highlight and output
			// 	fmt.Fprintf(out, "%s\n", colors[i].Sprint(line))

			// 	matched = true
			// 	break
			// }
		}

		if matched {
			lastColorIdx := -1
			newLine := ""
			buff := ""
			for i, ch := range line {
				if lineColors[i] != lastColorIdx {
					if lastColorIdx != -1 {
						lastColor := colors[lastColorIdx]
						newLine += lastColor.Sprint(buff)
					} else {
						newLine += buff
					}
					buff = ""

					lastColorIdx = lineColors[i]
				}
				buff += string(ch)
			}
			if buff != "" {
				if lastColorIdx != -1 {
					lastColor := colors[lastColorIdx]
					newLine += lastColor.Sprint(buff)
				} else {
					newLine += buff
				}
			}
			fmt.Fprintf(out, "%s\n", newLine)
		}

		// No matches, print with default color
		if !matched {
			fmt.Fprintf(out, "%s\n", scanner.Text())
		}
	}
}
