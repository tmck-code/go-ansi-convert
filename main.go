package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tmck-code/go-ansi-convert/src/convert"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	mode := os.Args[1]

	switch mode {
	case "flip":
		runFlip()
	case "sanitise":
		runSanitise()
	default:
		fmt.Fprintf(os.Stderr, "Unknown mode: %s\n\n", mode)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  ansi-flip flip <input> <output>")
	fmt.Fprintln(os.Stderr, "  ansi-flip sanitise [-justify] <input> <output>")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Modes:")
	fmt.Fprintln(os.Stderr, "  flip      Reverse/flip ANSI text")
	fmt.Fprintln(os.Stderr, "  sanitise  Clean up ANSI codes")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Options for sanitise mode:")
	fmt.Fprintln(os.Stderr, "  -justify  Justify lines to the same length")
}

func runFlip() {
	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "Error: flip mode requires <input> and <output> file paths")
		fmt.Fprintln(os.Stderr, "Usage: ansi-flip flip <input> <output>")
		os.Exit(1)
	}

	inputFile := os.Args[2]
	outputFile := os.Args[3]

	input := readFile(inputFile)
	tokenized := convert.TokeniseANSIString(input)
	reversed := convert.ReverseANSIString(tokenized)
	output := convert.BuildANSIString(reversed, 0)
	writeFile(outputFile, output)
}

func runSanitise() {
	// Parse sanitise-specific flags
	sanitiseFlags := flag.NewFlagSet("sanitise", flag.ExitOnError)
	justify := sanitiseFlags.Bool("justify", false, "Justify lines to the same length")

	// Parse flags starting from os.Args[2:]
	sanitiseFlags.Parse(os.Args[2:])

	args := sanitiseFlags.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: sanitise mode requires <input> and <output> file paths")
		fmt.Fprintln(os.Stderr, "Usage: ansi-flip sanitise [-justify] <input> <output>")
		os.Exit(1)
	}

	inputFile := args[0]
	outputFile := args[1]

	input := readFile(inputFile)
	output := convert.SanitiseUnicodeString(input, *justify)
	writeFile(outputFile, output)
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", path, err)
		os.Exit(1)
	}
	content := string(data)
	// Trim trailing newline if present
	if len(content) > 0 && content[len(content)-1] == '\n' {
		content = content[:len(content)-1]
	}
	return content
}

func writeFile(path string, content string) {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file %s: %v\n", path, err)
		os.Exit(1)
	}
}
