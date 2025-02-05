package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ha1tch/zxgotools/pkg/basic"
	"github.com/ha1tch/zxgotools/pkg/tap"
)

func main() {
	var (
		basicMode = flag.Bool("basic", false, "Convert BASIC text file")
		binMode = flag.Bool("binary", false, "Convert binary file")
		name = flag.String("name", "", "Name for TAP block (max 10 chars)")
		address = flag.Uint("address", 32768, "Start address for binary files (default: 32768)")
		autostart = flag.Uint("autostart", 0, "Auto-start line for BASIC programs")
		caseIndependent = flag.Bool("c", false, "Case independent token matching")
	)

	flag.Parse()

	if !*basicMode && !*binMode {
		fmt.Fprintf(os.Stderr, "Error: Must specify either --basic or --binary mode\n")
		flag.Usage()
		os.Exit(1)
	}

	if *basicMode && *binMode {
		fmt.Fprintf(os.Stderr, "Error: Cannot specify both --basic and --binary\n")
		flag.Usage()
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--basic|--binary] [options] input output.tap\n", os.Args[0])
		flag.Usage()
		os.Exit(1)
	}

	inputFile, outputFile := args[0], args[1]

	if *basicMode {
		opts := []basic.Option{}
		if *caseIndependent {
			opts = append(opts, basic.WithCaseIndependent(true))
		}
		
		if err := convertBasic(inputFile, outputFile, *name, uint16(*autostart), opts...); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := convertBinary(inputFile, outputFile, *name, uint16(*address)); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("Successfully created %s\n", outputFile)
}

func convertBinary(inputFile, outputFile, name string, startAddress uint16) error {
	return tap.BinaryToTAP(inputFile, outputFile, name, startAddress)
}

func convertBasic(inputFile, outputFile, name string, autostart uint16, opts ...basic.Option) error {
	// Read and parse BASIC
	input, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("opening input file: %w", err)
	}
	defer input.Close()

	parser := basic.NewParser(opts...)
	data, err := parser.Parse(input)
	if err != nil {
		return fmt.Errorf("parsing BASIC: %w", err)
	}

	// Create output file
	out, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer out.Close()

	// Use filename if no name provided
	if name == "" {
		name = filepath.Base(inputFile)
		name = strings.TrimSuffix(name, filepath.Ext(name))
		if len(name) > 10 {
			name = name[:10]
		}
	}

	// Write as TAP
	if err := tap.WriteBasicToTAP(out, name, data, autostart); err != nil {
		return fmt.Errorf("writing TAP file: %w", err)
	}

	// Report 128K requirement if detected
	if parser.Is128K() {
		fmt.Println("Note: Program requires 128K")
	}

	return nil
}