package logh_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/haroflow/logh"
)

func TestLogHighlight(t *testing.T) {
	color.NoColor = false // For testing, so we can compare colored strings

	t.Run("empty expression matches nothing and prints all lines with default color", func(t *testing.T) {
		input := strings.NewReader("line1\n" +
			"line2\n" +
			"partial match line2\n" +
			"line3\n" +
			"line2")

		want := "line1\n" +
			"line2\n" +
			"partial match line2\n" +
			"line3\n" +
			"line2\n" // An extra \n there

		got := &bytes.Buffer{}
		logh.Highlight(input, got, "")
		assertStringEqual(t, got.String(), want)
	})

	t.Run("highlight all lines containing the word", func(t *testing.T) {
		input := strings.NewReader("line1\n" +
			"line2\n" +
			"partial match line2\n" +
			"line3\n" +
			"line2")

		want := "line1\n" +
			color.RedString("line2") + "\n" +
			color.RedString("partial match line2") + "\n" +
			"line3\n" +
			color.RedString("line2") + "\n" // An extra \n there

		got := &bytes.Buffer{}
		logh.Highlight(input, got, "line2")
		assertStringEqual(t, got.String(), want)
	})

	t.Run("accepts multiple expressions, highlights with a color for each one", func(t *testing.T) {
		input := strings.NewReader("a\n" +
			"b\n" +
			"c\n" +
			"d\n" +
			"e")

		want := "a\n" +
			color.RedString("b") + "\n" +
			color.GreenString("c") + "\n" +
			color.BlueString("d") + "\n" +
			color.YellowString("e") + "\n" // An extra \n there

		got := &bytes.Buffer{}
		logh.Highlight(input, got, "b", "c", "d", "e")
		assertStringEqual(t, got.String(), want)
	})

	t.Run("starts reusing colors if number of matches is greater than number of colors", func(t *testing.T) {
		chars := "abcdefghijklmnopqrstuvwxyz"

		inputStr := ""
		want := ""
		for i := 0; i < len(chars); i++ {
			ch := string(chars[i])
			inputStr += ch + "\n"

			c := logh.Colors[i%len(logh.Colors)]
			want += color.New(c).Sprint(ch) + "\n"
		}
		input := strings.NewReader(inputStr)

		got := &bytes.Buffer{}
		logh.Highlight(input, got, strings.Split(chars, "")...)
		assertStringEqual(t, got.String(), want)
	})
}

func assertStringEqual(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("\nGOT:\n%s\n\nWANT:\n%s", got, want)
	}
}
