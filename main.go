package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pborman/getopt/v2"
	"github.com/tmck-code/go-ansi-convert/src/ansi-convert/convert"
	"github.com/tmck-code/go-ansi-convert/src/ansi-convert/log"
	"github.com/tmck-code/go-ansi-convert/src/ansi-convert/parse"
)

type Args struct {
	InputFile             string
	OutputFile            string
	Stdin                 bool
	Stdout                bool
	FlipHorizontal        bool
	FlipVertical          bool
	Sanitise              bool
	Optimise              bool
	Justify               bool
	Help                  bool
	Display               bool
	DisplaySeparator      string
	DisplaySeparatorWidth int
	DisplaySwapped        bool
	ConvertAns            bool
	DisplaySAUCEInfo      bool
	DisplaySAUCEInfoJSON  bool
	DetectEncoding        bool
}

// Check DEBUG mode, enables debug logging
func Debug() bool {
	debugValue := os.Getenv("DEBUG")
	return debugValue == "true" || debugValue == "1"
}

func main() {
	getopt.SetProgram("ansi-flip")

	help := getopt.BoolLong("help", 'h', "display this help message")

	inputFile := getopt.StringLong("input", 'i', "", "Input file path (default: stdin)")
	outputFile := getopt.StringLong("output", 'o', "", "Output file path (default: stdout)")

	flip := getopt.EnumLong("flip", 'f', []string{"h", "v", "h,v", "v,h"}, "", "Flip horizontally (h), vertically (v), or both (h,v or v,h)")
	getopt.BoolLong("sanitise", 's', "Sanitise ANSI lines, ensuring that each line ends with a reset code")
	justify := getopt.BoolLong("justify", 'j', "Justify lines to the same length (sanitise mode only)")
	optimise := getopt.BoolLong("optimise", 'O', "Optimise ANSI tokens to merge redundant color codes")
	display := getopt.BoolLong("display", 'd', "Display original and flipped side-by-side in terminal")

	convertAns := getopt.BoolLong("convert-ans", 'c', "Convert an ANSI .ans file (CP437 encoded) to UTF-8 ANSI")
	displaySAUCE := getopt.BoolLong("display-sauce", 'S', "Display SAUCE metadata from input file (if present)")
	displaySAUCEInfoJSON := getopt.BoolLong("display-sauce-json", 0, "Display SAUCE metadata from input file in JSON format (if present)")
	detectEncoding := getopt.BoolLong("detect-encoding", 'e', "Detect if input file is CP437 or ISO-8859-1 encoded")

	displaySep := getopt.StringLong("display-separator", 0, " ", "Separator string between original and flipped when displaying")
	displaySepWidth := getopt.IntLong("display-separator-width", 0, 1, "Width of separator between original and flipped when displaying")
	displaySwapped := getopt.BoolLong("display-swapped", 'x', "When displaying, reverse the order of original and flipped")

	getopt.Lookup("convert-ans").SetGroup("operation")
	getopt.Lookup("flip").SetGroup("operation")
	getopt.Lookup("sanitise").SetGroup("operation")
	getopt.Lookup("help").SetGroup("operation")
	getopt.Lookup("optimise").SetGroup("operation")
	getopt.Lookup("display-sauce").SetGroup("operation")
	getopt.Lookup("detect-encoding").SetGroup("operation")
	getopt.RequiredGroup("operation")

	getopt.Parse()

	args := Args{
		InputFile:             *inputFile,
		OutputFile:            *outputFile,
		Stdin:                 !getopt.IsSet("input"),
		Stdout:                !getopt.IsSet("output"),
		FlipHorizontal:        strings.Contains(*flip, "h"),
		FlipVertical:          strings.Contains(*flip, "v"),
		Sanitise:              getopt.IsSet("sanitise"),
		Optimise:              *optimise,
		Justify:               *justify,
		Help:                  *help,
		Display:               *display,
		DisplaySeparator:      *displaySep,
		DisplaySeparatorWidth: *displaySepWidth,
		DisplaySwapped:        *displaySwapped,
		ConvertAns:            *convertAns,
		DisplaySAUCEInfo:      *displaySAUCE,
		DisplaySAUCEInfoJSON:  *displaySAUCEInfoJSON,
		DetectEncoding:        *detectEncoding,
	}

	if args.Help {
		getopt.Usage()
		return
	}

	// 1. always detect the encoding
	// 2. always read & separate the SAUCE record (if present)

	encoding, input := readInput(args)
	if args.DetectEncoding {
		fmt.Printf("Detected encoding: \x1b[93m%s\x1b[0m\n", encoding)
		return
	}

	var sauce *convert.SAUCE
	var fileData string

	sauce, fileData, err := convert.ParseSAUCE(input)
	if err != nil {
		log.DebugFprintf("Error parsing SAUCE record: %v\n", err)
		sauce, fileData, err = convert.CreateSAUCERecord(args.InputFile)
		if err != nil {
			log.DebugFprintf("Error creating SAUCE record: %v\n", err)
			os.Exit(1)
		}
	}
	if args.DisplaySAUCEInfo {
		fmt.Println(sauce.ToString())
		return
	}

	result := process(args, fileData, sauce)

	if args.Display {
		if args.FlipHorizontal {
			displaySideBySide(input, result, args)
		} else if args.FlipVertical {
			displayAboveBelow(input, result, args)
		}
	} else {
		writeOutput(args, result)
	}
}

// displaySideBySide prints the original and flipped result side-by-side, separated by a space
func displaySideBySide(original, flipped string, args Args) {
	origLines := strings.Split(convert.SanitiseUnicodeString(original, true), "\n")
	flippedLines := strings.Split(flipped, "\n")
	for i := 0; i < len(origLines) && i < len(flippedLines); i++ {
		if !args.DisplaySwapped {
			fmt.Printf(
				"%s%s%s\n",
				origLines[i],
				strings.Repeat(args.DisplaySeparator, args.DisplaySeparatorWidth),
				flippedLines[i],
			)
		} else {
			fmt.Printf(
				"%s%s%s\n",
				flippedLines[i],
				strings.Repeat(args.DisplaySeparator, args.DisplaySeparatorWidth),
				origLines[i],
			)
		}
	}
}

func displayAboveBelow(original, flipped string, args Args) {
	sep := strings.Repeat(args.DisplaySeparator, args.DisplaySeparatorWidth)
	if sep != "" {
		sep += "\n"
	}
	if !args.DisplaySwapped {
		fmt.Printf(
			"%s%s%s",
			convert.SanitiseUnicodeString(original, true),
			sep,
			flipped,
		)
	} else {
		fmt.Printf(
			"%s%s%s",
			flipped,
			sep,
			convert.SanitiseUnicodeString(original, true),
		)
	}
}

func readInput(args Args) (string, string) {
	var raw []byte
	var err error

	if args.Stdin {
		raw, err = io.ReadAll(os.Stdin)
	} else {
		raw, err = os.ReadFile(args.InputFile)
	}
	if err != nil {
		log.DebugFprintf("Error reading stdin: %v\n", err)
		os.Exit(1)
	}

	encoding := parse.DetectEncoding(raw)
	data, err := parse.DecodeFileContents(raw, encoding)
	if err != nil {
		log.DebugFprintf("Error decoding file contents: %v\n", err)
		os.Exit(1)
	}

	return encoding, data
}

func process(args Args, input string, sauce *convert.SAUCE) string {
	if args.ConvertAns {
		return convert.ConvertAns(input, *sauce)
	}
	if args.Optimise {
		tokenized := convert.TokeniseANSIString(input)
		optimised := convert.OptimiseANSITokens(tokenized)
		return convert.BuildANSIString(optimised, 0)
	}
	if args.Sanitise {
		return convert.SanitiseUnicodeString(input, args.Justify)
	}
	return runFlip(input, args)
}

func runFlip(input string, args Args) string {
	tokenized := convert.TokeniseANSIString(input)
	if args.FlipHorizontal {
		tokenized = convert.FlipHorizontal(tokenized)
	}
	if args.FlipVertical {
		tokenized = convert.FlipVertical(tokenized)
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
		log.DebugFprintf("Error writing file %s: %v\n", path, err)
		os.Exit(1)
	}
}
