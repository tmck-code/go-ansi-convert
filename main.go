package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pborman/getopt/v2"
	"github.com/tmck-code/go-ansi-convert/src/convert"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
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

	displaySep := getopt.StringLong("display-separator", 0, " ", "Separator string between original and flipped when displaying")
	displaySepWidth := getopt.IntLong("display-separator-width", 0, 1, "Width of separator between original and flipped when displaying")
	displaySwapped := getopt.BoolLong("display-swapped", 'x', "When displaying, reverse the order of original and flipped")

	getopt.Lookup("convert-ans").SetGroup("operation")
	getopt.Lookup("flip").SetGroup("operation")
	getopt.Lookup("sanitise").SetGroup("operation")
	getopt.Lookup("help").SetGroup("operation")
	getopt.Lookup("optimise").SetGroup("operation")
	getopt.Lookup("display-sauce").SetGroup("operation")
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
	}

	if args.Help {
		getopt.Usage()
		return
	}

	if args.DisplaySAUCEInfo {
		sauceInfo, _, err := convert.ParseSAUCEFromFile(args.InputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing SAUCE record: %v\n", err)
			os.Exit(1)
		}
		if args.DisplaySAUCEInfoJSON {
			jsonStr, err := sauceInfo.ToJSON()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error converting SAUCE info to JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(jsonStr)
		} else {
			fmt.Println(sauceInfo.ToString())
		}
		return
	}

	input := readInput(args)
	result := process(args, input)

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

func readInput(args Args) string {
	if args.Stdin {
		return readStdin()
	} else {
		return readFile(args.InputFile, args.ConvertAns)
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

func readFile(path string, decodeCP437 bool) string {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", path, err)
		os.Exit(1)
	}

	if decodeCP437 {
		// Decode from CP437 to UTF-8
		decoder := charmap.CodePage437.NewDecoder()
		reader := transform.NewReader(bytes.NewReader(data), decoder)
		decoded, err := io.ReadAll(reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding CP437: %v\n", err)
			os.Exit(1)
		}
		return string(decoded)
	}

	return string(data)
}

func process(args Args, input string) string {
	if args.ConvertAns {
		sauce, fileData, err := convert.ParseSAUCEFromFile(args.InputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing SAUCE record: %v\n", err)
			os.Exit(1)
		}
		return convert.ConvertAns(fileData, *sauce)
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
		fmt.Fprintf(os.Stderr, "Error writing file %s: %v\n", path, err)
		os.Exit(1)
	}
}
