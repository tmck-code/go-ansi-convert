package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pborman/getopt/v2"
	"github.com/tmck-code/go-ansi-convert/src/convert"
)

type Args struct {
	InputFile       string
	OutputFile      string
	Stdin           bool
	Stdout          bool
	FlipHorizontal  bool
	FlipVertical    bool
	Sanitise        bool
	Justify         bool
	Help			bool
}

func main() {
	getopt.SetProgram("ansi-flip")

	help := getopt.BoolLong("help", 'h', "display this help message")

	inputFile := getopt.StringLong("input", 'i', "", "Input file path (default: stdin)")
	outputFile := getopt.StringLong("output", 'o', "", "Output file path (default: stdout)")

	flip := getopt.EnumLong("flip", 'f', []string{"h", "v", "h,v", "v,h"}, "", "Flip horizontally (h), vertically (v), or both (h,v or v,h)")
	getopt.BoolLong("sanitise", 's', "Sanitise ANSI lines, ensuring that each line ends with a reset code")
	justify := getopt.BoolLong("justify", 'j', "Justify lines to the same length (sanitise mode only)")

	getopt.Lookup("flip").SetGroup("operation")
	getopt.Lookup("sanitise").SetGroup("operation")
	getopt.Lookup("help").SetGroup("operation")
	getopt.RequiredGroup("operation")

	getopt.Parse()

	args := Args{
		InputFile:      *inputFile,
		OutputFile:     *outputFile,
		Stdin:          !getopt.IsSet("input"),
		Stdout:         !getopt.IsSet("output"),
		FlipHorizontal: strings.Contains(*flip, "h"),
		FlipVertical:   strings.Contains(*flip, "v"),
		Sanitise:       getopt.IsSet("sanitise"),
		Justify:        *justify,
		Help:			*help,
	}

	if args.Help {
		getopt.Usage()
		return
	}

	writeOutput(args, process(args, readInput(args)))
}

func readInput(args Args) string {
	if args.Stdin {
		return readStdin()
	} else {
		return readFile(args.InputFile)
	}
}

func readStdin() string {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(1)
	}
	return string(data)
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", path, err)
		os.Exit(1)
	}
	return string(data)
}

func process(args Args, input string) string {
	if args.Sanitise {
		return convert.SanitiseUnicodeString(input, args.Justify)
	}
	return runFlip(input, args)
}

func runFlip(input string, args Args) string {
	tokenized := convert.TokeniseANSIString(input)
	if args.FlipHorizontal {
		tokenized = convert.ReverseANSIString(tokenized)
	}
	return convert.BuildANSIString(tokenized, 0)
}

func writeOutput(args Args, output string) {
	if args.Stdout {
		fmt.Print(output)
	} else {
		writeFile(args.OutputFile, output)
	}
}

func writeFile(path string, content string) {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file %s: %v\n", path, err)
		os.Exit(1)
	}
}
