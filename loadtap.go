package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	headerFlag = 0x00
	dataFlag   = 0xFF
)

// TAPHeader represents a TAP file header block
type TAPHeader struct {
	Type        byte
	Filename    string
	DataLength  uint16
	Param1      uint16
	Param2      uint16
	Checksum    byte
}

// TAPBlock represents a block in a TAP file
type TAPBlock struct {
	Length   uint16
	Flag     byte
	Data     []byte
	Checksum byte
	Header   *TAPHeader
}

// parseHeader creates a TAPHeader from raw bytes
func parseHeader(headerData []byte, checksum byte) (*TAPHeader, error) {
	if len(headerData) != 17 { // Header should be exactly 17 bytes (without checksum)
		return nil, fmt.Errorf("invalid header data length: %d", len(headerData))
	}

	// Extract filename (10 bytes) and trim spaces
	filename := strings.TrimRight(string(headerData[1:11]), " ")

	return &TAPHeader{
		Type:       headerData[0],
		Filename:   filename,
		DataLength: binary.LittleEndian.Uint16(headerData[11:13]),
		Param1:     binary.LittleEndian.Uint16(headerData[13:15]),
		Param2:     binary.LittleEndian.Uint16(headerData[15:17]),
		Checksum:   checksum,
	}, nil
}

// readTAPFile reads a TAP file and returns a slice of TAPBlock
func readTAPFile(filename string) ([]TAPBlock, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	var blocks []TAPBlock

	for {
		// Read block length (2 bytes, little-endian)
		var length uint16
		if err := binary.Read(file, binary.LittleEndian, &length); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("reading block length: %w", err)
		}

		// Read flag byte
		var flag byte
		if err := binary.Read(file, binary.LittleEndian, &flag); err != nil {
			return nil, fmt.Errorf("reading flag: %w", err)
		}

		// Read block data (excluding checksum)
		data := make([]byte, length-2) // -2 for flag and checksum
		if _, err := io.ReadFull(file, data); err != nil {
			return nil, fmt.Errorf("reading data: %w", err)
		}

		// Read checksum
		var checksum byte
		if err := binary.Read(file, binary.LittleEndian, &checksum); err != nil {
			return nil, fmt.Errorf("reading checksum: %w", err)
		}

		block := TAPBlock{
			Length:   length,
			Flag:     flag,
			Data:     data,
			Checksum: checksum,
		}

		// If this is a header block (flag = 0x00 and length = 0x13)
		if flag == headerFlag && length == 0x13 {
			if header, err := parseHeader(data, checksum); err == nil {
				block.Header = header
			} else {
				fmt.Fprintf(os.Stderr, "Warning: Failed to parse header: %v\n", err)
			}
		}

		blocks = append(blocks, block)
	}

	return blocks, nil
}

// printBlockInfo prints information about a TAP block
func printBlockInfo(block TAPBlock, index int) {
	fmt.Printf("\nBlock %d:\n", index)
	fmt.Printf("  Length: %d\n", block.Length)
	fmt.Printf("  Flag: 0x%02X (%s)\n", block.Flag, flagTypeString(block.Flag))

	if block.Header != nil {
		fmt.Println("  Header Information:")
		fmt.Printf("    Type: %d\n", block.Header.Type)
		fmt.Printf("    Filename: %s\n", block.Header.Filename)
		fmt.Printf("    Data Length: %d\n", block.Header.DataLength)
		fmt.Printf("    Param1: %d\n", block.Header.Param1)
		fmt.Printf("    Param2: %d\n", block.Header.Param2)
	}
	fmt.Printf("  Checksum: 0x%02X\n", block.Checksum)
	fmt.Printf("  Data Length: %d bytes\n", len(block.Data))
}

// flagTypeString returns a string description of the flag type
func flagTypeString(flag byte) string {
	if flag == headerFlag {
		return "Header"
	}
	return "Data"
}

// dumpHex prints a hexadecimal dump of the data
func dumpHex(data []byte) {
	fmt.Println("  Data:")
	const bytesPerLine = 16
	for i := 0; i < len(data); i += bytesPerLine {
		end := i + bytesPerLine
		if end > len(data) {
			end = len(data)
		}
		line := data[i:end]
		hexBytes := make([]string, len(line))
		for j, b := range line {
			hexBytes[j] = fmt.Sprintf("%02X", b)
		}
		fmt.Printf("    %s\n", strings.Join(hexBytes, " "))
	}
}

func main() {
	dump := flag.Bool("d", false, "Dump block data as hex")
	raw := flag.Bool("r", false, "Output raw block data")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-d] [-r] <tap-file>\n", os.Args[0])
		os.Exit(1)
	}

	filename := flag.Arg(0)
	blocks, err := readTAPFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *raw {
		// Output just the raw data blocks (skip headers)
		for _, block := range blocks {
			if block.Flag != headerFlag {
				os.Stdout.Write(block.Data)
			}
		}
	} else {
		fmt.Printf("Found %d blocks in %s\n", len(blocks), filename)
		for i, block := range blocks {
			printBlockInfo(block, i)
			if *dump {
				dumpHex(block.Data)
			}
		}
	}
}