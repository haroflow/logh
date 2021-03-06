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
		logh.Highlight(input, got, logh.HighlightConfig{}, "")
		assertStringEqual(t, got.String(), want)
	})

	// TODO: Add config for this? The user can pass .*word.*, which makes sense...
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
		logh.Highlight(input, got, logh.HighlightConfig{}, ".*line2.*")
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
		logh.Highlight(input, got, logh.HighlightConfig{}, "b", "c", "d", "e")
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
		logh.Highlight(input, got, logh.HighlightConfig{}, strings.Split(chars, "")...)
		assertStringEqual(t, got.String(), want)
	})

	t.Run("ignore case", func(t *testing.T) {
		input := strings.NewReader("line1\nLINE2\nLiNe3\n")
		want := color.RedString("line") + "1\n" +
			color.RedString("LINE") + "2\n" +
			color.RedString("LiNe") + "3\n"

		got := &bytes.Buffer{}
		config := logh.HighlightConfig{
			IgnoreCase: true,
		}
		logh.Highlight(input, got, config, "line")
		assertStringEqual(t, got.String(), want)
	})

	t.Run("allow highlighting matches only, multiple expressions", func(t *testing.T) {
		input := strings.NewReader("line1\nline2\nline3\nfourth line\n")
		want := "line1\n" +
			"line" + color.RedString("2") + "\n" +
			color.GreenString("line3") + "\n" +
			color.BlueString("fourth") + " line\n"

		got := &bytes.Buffer{}
		logh.Highlight(input, got, logh.HighlightConfig{}, "2", "line3", "fourth")
		assertStringEqual(t, got.String(), want)
	})

	t.Run("allow highlighting matches only, multiple expressions on the same line", func(t *testing.T) {
		input := strings.NewReader("col1 col2 col3\n")
		want := "col1 col" + color.RedString("2") + " " + color.GreenString("col3") + "\n"

		got := &bytes.Buffer{}
		logh.Highlight(input, got, logh.HighlightConfig{}, "2", "col3")
		assertStringEqual(t, got.String(), want)
	})

	t.Run("should error when passed an invalid regular expression", func(t *testing.T) {
		input := strings.NewReader("abc\n")

		got := &bytes.Buffer{}
		err := logh.Highlight(input, got, logh.HighlightConfig{}, "*test")

		if err == nil {
			t.Errorf("expected an error, got %s", err)
		}
	})

	// t.Run("show debug output", func(t *testing.T) {
	// 	t.Error("TODO")
	// })

	// t.Run("should warn if there are no expressions", func(t *testing.T) {
	// 	t.Error("TODO")
	// })
}

func assertStringEqual(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("\nGOT:\n%s\n\nWANT:\n%s", got, want)
	}
}
