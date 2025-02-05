package main

import (
	"flag"
	"fmt"
	"os"

	"zxgotools/pkg/tap"
)

func main() {
	var (
		name    = flag.String("name", "", "Name for code block (max 10 chars)")
		address = flag.Uint("address", 32768, "Start address (default: 32768)")
	)

	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--name NAME] [--address ADDR] input.bin output.tap\n", os.Args[0])
		os.Exit(1)
	}

	inputFile, outputFile := args[0], args[1]

	if err := tap.BinaryToTAP(inputFile, outputFile, *name, uint16(*address)); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %s to %s\n", inputFile, outputFile)
}