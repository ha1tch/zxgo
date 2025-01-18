// main.go
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"unsafe"
)

func parseFlags() (*options, error) {
	opts := &options{}

	// Basic options
	flag.StringVar(&opts.output, "o", "", "output TZX file (required)")
	flag.StringVar(&opts.configFile, "c", "", "YAML configuration file")
	flag.UintVar((*uint)(unsafe.Pointer(&opts.pauseDuration)), "p", 1000, "pause duration between blocks in ms")

	// Metadata options
	flag.BoolVar(&opts.addArchive, "m", false, "add metadata block")
	flag.StringVar(&opts.title, "title", "", "program title (requires -m)")
	flag.StringVar(&opts.author, "author", "", "program author (requires -m)")
	flag.StringVar(&opts.year, "year", "", "publication year (requires -m)")

	// 128K options
	flag.BoolVar(&opts.k128Only, "128", false, "program requires 128K")
	flag.BoolVar(&opts.useAY, "ay", false, "program uses AY sound chip")
	flag.BoolVar(&opts.usePaging, "paging", false, "program uses memory paging")
	flag.StringVar(&opts.modelType, "model", "", "required model: +2, +2A, or +3")
	flag.BoolVar(&opts.multiload, "multiload", false, "program is multiload (adds 48K stop blocks)")

	// Group option
	flag.StringVar(&opts.currentGroup, "group", "", "group name for following files")

	flag.Parse()

	if opts.output == "" {
		return nil, fmt.Errorf("output file (-o) is required")
	}

	if flag.NArg() == 0 && opts.configFile == "" {
		return nil, fmt.Errorf("no input files or config file specified")
	}

	return opts, nil
}

func writeInitialBlocks(w io.Writer, opts *options) error {
	// Write TZX header
	if err := writeTzxHeader(w); err != nil {
		return fmt.Errorf("writing TZX header: %w", err)
	}

	// Write archive info if requested
	if err := writeArchiveInfo(w, *opts); err != nil {
		return fmt.Errorf("writing archive info: %w", err)
	}

	// Write hardware info if any 128K options are set
	if err := writeHardwareInfo(w, *opts); err != nil {
		return fmt.Errorf("writing hardware info: %w", err)
	}

	return nil
}

func processInputFile(w io.Writer, filename string, opts *options) error {
	// Write description block if file has extension .desc
	if filepath.Ext(filename) == ".desc" {
		desc, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("reading description file: %w", err)
		}
		if err := writeTextDescription(w, string(desc)); err != nil {
			return fmt.Errorf("writing description: %w", err)
		}
		return nil
	}

	// Process TAP file
	inFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer inFile.Close()

	for {
		data, err := readTapBlock(inFile)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading TAP block: %w", err)
		}

		if err := writeStandardSpeedBlock(w, data, opts.pauseDuration); err != nil {
			return fmt.Errorf("writing TZX block: %w", err)
		}
	}

	return nil
}

func processCommandLineInputs(w io.Writer, opts *options) error {
	var currentGroup bool

	for i := 0; i < flag.NArg(); i++ {
		filename := flag.Arg(i)

		// Handle group start if group name is set and group not yet started
		if opts.currentGroup != "" && !currentGroup {
			if err := writeGroupStart(w, opts.currentGroup); err != nil {
				return fmt.Errorf("writing group start: %w", err)
			}
			currentGroup = true
		}

		if err := processInputFile(w, filename, opts); err != nil {
			return fmt.Errorf("processing %s: %w", filename, err)
		}

		// Add 48K stop block if multiload (except after last file)
		if opts.multiload && i < flag.NArg()-1 {
			if err := write48KStopBlock(w); err != nil {
				return fmt.Errorf("writing 48K stop block: %w", err)
			}
		}
	}

	// Close current group if open
	if currentGroup {
		if err := writeGroupEnd(w); err != nil {
			return fmt.Errorf("writing group end: %w", err)
		}
	}

	return nil
}

func main() {
	opts, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		flag.Usage()
		os.Exit(1)
	}

	outFile, err := os.Create(opts.output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	if err := writeInitialBlocks(outFile, opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Handle config file if specified
	if opts.configFile != "" {
		config, err := processConfig(opts.configFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing config file: %v\n", err)
			os.Exit(1)
		}

		if err := processConfiguredBlocks(outFile, config); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing blocks from config: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Process command line inputs
	if err := processCommandLineInputs(outFile, opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}