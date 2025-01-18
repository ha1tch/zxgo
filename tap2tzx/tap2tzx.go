package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	tzxSignature = "ZXTape!"
	tzxEOF       = 0x1A
	tzxMajorVer  = 1
	tzxMinorVer  = 20
)

type blockConfig struct {
	File      string `yaml:"file,omitempty"`
	Group     string `yaml:"group,omitempty"`
	Desc      string `yaml:"desc,omitempty"`
	JumpTo    string `yaml:"jump_to,omitempty"`
	ID        string `yaml:"id,omitempty"`
	LoopStart int    `yaml:"loop_start,omitempty"`
	LoopEnd   bool   `yaml:"loop_end,omitempty"`
}

type tzxConfig struct {
	Metadata struct {
		Title  string `yaml:"title"`
		Author string `yaml:"author"`
		Year   string `yaml:"year"`
	} `yaml:"metadata"`
	Hardware struct {
		K128Only  bool   `yaml:"128k_only"`
		UseAY     bool   `yaml:"use_ay"`
		UsePaging bool   `yaml:"use_paging"`
		Model     string `yaml:"model"`
	} `yaml:"hardware"`
	Blocks []blockConfig `yaml:"blocks"`
}

type options struct {
	// File options
	output     string
	configFile string

	// Basic options
	pauseDuration uint16
	addArchive    bool
	title         string
	author        string
	year          string

	// 128K options
	k128Only  bool
	useAY     bool
	usePaging bool
	modelType string
	multiload bool

	// Simple grouping
	currentGroup string
}

// writeTzxHeader writes the initial TZX file signature and version
func writeTzxHeader(w io.Writer) error {
	if _, err := w.Write([]byte(tzxSignature)); err != nil {
		return fmt.Errorf("writing signature: %w", err)
	}
	if err := binary.Write(w, binary.LittleEndian, uint8(tzxEOF)); err != nil {
		return fmt.Errorf("writing EOF marker: %w", err)
	}
	if err := binary.Write(w, binary.LittleEndian, uint8(tzxMajorVer)); err != nil {
		return fmt.Errorf("writing major version: %w", err)
	}
	if err := binary.Write(w, binary.LittleEndian, uint8(tzxMinorVer)); err != nil {
		return fmt.Errorf("writing minor version: %w", err)
	}
	return nil
}

// writeArchiveInfo writes an archive info block (0x32)
func writeArchiveInfo(w io.Writer, opts options) error {
	if !opts.addArchive {
		return nil
	}

	// Count how many fields we'll write
	numFields := 0
	totalLength := 0
	if opts.title != "" {
		numFields++
		totalLength += len(opts.title) + 2 // +2 for ID and length bytes
	}
	if opts.author != "" {
		numFields++
		totalLength += len(opts.author) + 2
	}
	if opts.year != "" {
		numFields++
		totalLength += len(opts.year) + 2
	}

	if numFields == 0 {
		return nil
	}

	// Write block ID
	if err := binary.Write(w, binary.LittleEndian, uint8(0x32)); err != nil {
		return fmt.Errorf("writing archive block ID: %w", err)
	}

	// Write total length (2 bytes for length + 1 byte for number of strings + text data)
	if err := binary.Write(w, binary.LittleEndian, uint16(totalLength+1)); err != nil {
		return fmt.Errorf("writing archive block length: %w", err)
	}

	// Write number of text strings
	if err := binary.Write(w, binary.LittleEndian, uint8(numFields)); err != nil {
		return fmt.Errorf("writing number of text strings: %w", err)
	}

	// Write text fields
	if opts.title != "" {
		if err := binary.Write(w, binary.LittleEndian, uint8(0)); err != nil { // Title ID
			return fmt.Errorf("writing title ID: %w", err)
		}
		if err := binary.Write(w, binary.LittleEndian, uint8(len(opts.title))); err != nil {
			return fmt.Errorf("writing title length: %w", err)
		}
		if _, err := w.Write([]byte(opts.title)); err != nil {
			return fmt.Errorf("writing title: %w", err)
		}
	}

	if opts.author != "" {
		if err := binary.Write(w, binary.LittleEndian, uint8(2)); err != nil { // Author ID
			return fmt.Errorf("writing author ID: %w", err)
		}
		if err := binary.Write(w, binary.LittleEndian, uint8(len(opts.author))); err != nil {
			return fmt.Errorf("writing author length: %w", err)
		}
		if _, err := w.Write([]byte(opts.author)); err != nil {
			return fmt.Errorf("writing author: %w", err)
		}
	}

	if opts.year != "" {
		if err := binary.Write(w, binary.LittleEndian, uint8(3)); err != nil { // Year ID
			return fmt.Errorf("writing year ID: %w", err)
		}
		if err := binary.Write(w, binary.LittleEndian, uint8(len(opts.year))); err != nil {
			return fmt.Errorf("writing year length: %w", err)
		}
		if _, err := w.Write([]byte(opts.year)); err != nil {
			return fmt.Errorf("writing year: %w", err)
		}
	}

	return nil
}

// writeHardwareInfo writes a hardware type block (0x33)
func writeHardwareInfo(w io.Writer, opts options) error {
	if !opts.k128Only && !opts.useAY && !opts.usePaging && opts.modelType == "" {
		return nil
	}

	// Write block ID (0x33)
	if err := binary.Write(w, binary.LittleEndian, uint8(0x33)); err != nil {
		return fmt.Errorf("writing hardware block ID: %w", err)
	}

	// Count hardware entries
	entries := 0
	if opts.k128Only {
		entries += 2 // One for 128K required, one for 48K incompatible
	}
	if opts.useAY {
		entries++
	}
	if opts.modelType != "" {
		entries++
	}

	// Write number of hardware entries
	if err := binary.Write(w, binary.LittleEndian, uint8(entries)); err != nil {
		return fmt.Errorf("writing hardware entry count: %w", err)
	}

	if opts.k128Only {
		// 128K required
		if err := binary.Write(w, binary.LittleEndian, []byte{
			0x00, // Computer type
			0x03, // ZX Spectrum 128k
			0x01, // Uses this hardware
		}); err != nil {
			return fmt.Errorf("writing 128K requirement: %w", err)
		}

		// 48K incompatible
		if err := binary.Write(w, binary.LittleEndian, []byte{
			0x00, // Computer type
			0x01, // ZX Spectrum 48k
			0x03, // Doesn't work on this hardware
		}); err != nil {
			return fmt.Errorf("writing 48K incompatibility: %w", err)
		}
	}

	if opts.useAY {
		// Uses AY sound chip
		if err := binary.Write(w, binary.LittleEndian, []byte{
			0x03, // Sound device type
			0x00, // Classic AY
			0x01, // Uses this hardware
		}); err != nil {
			return fmt.Errorf("writing AY usage: %w", err)
		}
	}

	if opts.modelType != "" {
		var modelID uint8
		switch opts.modelType {
		case "+2":
			modelID = 0x04 // ZX Spectrum 128k +2 (grey case)
		case "+2A":
			modelID = 0x05 // ZX Spectrum 128k +2A
		case "+3":
			modelID = 0x05 // Also 0x05 as +3 is same hardware as +2A
		default:
			modelID = 0x03 // Default to basic 128K
		}

		if err := binary.Write(w, binary.LittleEndian, []byte{
			0x00,    // Computer type
			modelID, // Specific model
			0x01,    // Uses this hardware
		}); err != nil {
			return fmt.Errorf("writing model requirement: %w", err)
		}
	}

	return nil
} // writeStandardSpeedBlock writes a TAP block as a standard speed data block (0x10)
func writeStandardSpeedBlock(w io.Writer, data []byte, pause uint16) error {
	// Write block ID
	if err := binary.Write(w, binary.LittleEndian, uint8(0x10)); err != nil {
		return fmt.Errorf("writing block ID: %w", err)
	}

	// Write pause duration
	if err := binary.Write(w, binary.LittleEndian, pause); err != nil {
		return fmt.Errorf("writing pause duration: %w", err)
	}

	// Write data length
	if err := binary.Write(w, binary.LittleEndian, uint16(len(data))); err != nil {
		return fmt.Errorf("writing data length: %w", err)
	}

	// Write data
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("writing data: %w", err)
	}

	return nil
}

// write48KStopBlock writes a "Stop the tape if in 48K mode" block (0x2A)
func write48KStopBlock(w io.Writer) error {
	// Write block ID
	if err := binary.Write(w, binary.LittleEndian, uint8(0x2A)); err != nil {
		return fmt.Errorf("writing 48K stop block ID: %w", err)
	}

	// Write block length (always 4 bytes of 0)
	return binary.Write(w, binary.LittleEndian, uint32(0))
}

// writeTextDescription writes a text description block (0x30)
func writeTextDescription(w io.Writer, desc string) error {
	if desc == "" {
		return nil
	}

	// Write block ID
	if err := binary.Write(w, binary.LittleEndian, uint8(0x30)); err != nil {
		return fmt.Errorf("writing description block ID: %w", err)
	}

	// Write length of text
	length := uint8(len(desc))
	if err := binary.Write(w, binary.LittleEndian, length); err != nil {
		return fmt.Errorf("writing description length: %w", err)
	}

	// Write description text
	if _, err := w.Write([]byte(desc)); err != nil {
		return fmt.Errorf("writing description: %w", err)
	}

	return nil
}

// writeGroupStart writes a group start block (0x21)
func writeGroupStart(w io.Writer, name string) error {
	if name == "" {
		return nil
	}

	// Write block ID
	if err := binary.Write(w, binary.LittleEndian, uint8(0x21)); err != nil {
		return fmt.Errorf("writing group start block ID: %w", err)
	}

	// Write length of name
	length := uint8(len(name))
	if err := binary.Write(w, binary.LittleEndian, length); err != nil {
		return fmt.Errorf("writing group name length: %w", err)
	}

	// Write group name
	if _, err := w.Write([]byte(name)); err != nil {
		return fmt.Errorf("writing group name: %w", err)
	}

	return nil
}

// writeGroupEnd writes a group end block (0x22)
func writeGroupEnd(w io.Writer) error {
	// Write block ID
	return binary.Write(w, binary.LittleEndian, uint8(0x22))
}

// writeJumpBlock writes a jump block (0x23)
func writeJumpBlock(w io.Writer, offset int16) error {
	// Write block ID
	if err := binary.Write(w, binary.LittleEndian, uint8(0x23)); err != nil {
		return fmt.Errorf("writing jump block ID: %w", err)
	}

	// Write relative jump value
	return binary.Write(w, binary.LittleEndian, offset)
}

// writeLoopStart writes a loop start block (0x24)
func writeLoopStart(w io.Writer, repetitions uint16) error {
	// Write block ID
	if err := binary.Write(w, binary.LittleEndian, uint8(0x24)); err != nil {
		return fmt.Errorf("writing loop start block ID: %w", err)
	}

	// Write number of repetitions
	return binary.Write(w, binary.LittleEndian, repetitions)
}

// writeLoopEnd writes a loop end block (0x25)
func writeLoopEnd(w io.Writer) error {
	// Write block ID
	return binary.Write(w, binary.LittleEndian, uint8(0x25))
}

// readTapBlock reads a single TAP block from the input file
func readTapBlock(r io.Reader) ([]byte, error) {
	// Read block length
	var length uint16
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, fmt.Errorf("reading block length: %w", err)
	}

	// Read block data
	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, fmt.Errorf("reading block data: %w", err)
	}

	return data, nil
}

// processConfig reads and processes a YAML configuration file
func processConfig(filename string) (*tzxConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var config tzxConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	return &config, nil
}
