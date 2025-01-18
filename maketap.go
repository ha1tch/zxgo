package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	headerLength = 0x13
	headerFlag   = 0x00
	dataFlag     = 0xFF
	bytesType    = 0x03
)

// TAPHeader represents a TAP file header block
type TAPHeader struct {
	BlockLength uint16
	Flag        byte
	Type        byte
	Filename    [10]byte
	DataLength  uint16
	Param1      uint16 // Start address for bytes
	Param2      uint16 // Reserved for bytes
	Checksum    byte
}

// calculateChecksum performs XOR of all bytes
func calculateChecksum(data []byte) byte {
	var checksum byte
	for _, b := range data {
		checksum ^= b
	}
	return checksum
}

// createHeaderBlock creates a TAP header block for CODE format
func createHeaderBlock(filename string, dataLength uint16, startAddress uint16) []byte {
	header := TAPHeader{
		BlockLength: headerLength,
		Flag:        headerFlag,
		Type:        bytesType,
		DataLength:  dataLength,
		Param1:      startAddress,
		Param2:      0,
	}

	// Pad filename with spaces
	copy(header.Filename[:], []byte(filename))
	for i := len(filename); i < 10; i++ {
		header.Filename[i] = ' '
	}

	// Create buffer for header data
	headerData := make([]byte, 0, headerLength+2) // +2 for block length

	// Write block length
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, header.BlockLength)
	headerData = append(headerData, buf...)

	// Build header data for checksum calculation
	checksumData := make([]byte, 0, headerLength-1) // -1 as checksum isn't included
	checksumData = append(checksumData, header.Flag)
	checksumData = append(checksumData, header.Type)
	checksumData = append(checksumData, header.Filename[:]...)

	buf = make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, header.DataLength)
	checksumData = append(checksumData, buf...)

	binary.LittleEndian.PutUint16(buf, header.Param1)
	checksumData = append(checksumData, buf...)

	binary.LittleEndian.PutUint16(buf, header.Param2)
	checksumData = append(checksumData, buf...)

	// Calculate and append checksum
	header.Checksum = calculateChecksum(checksumData)

	// Build final header block
	headerData = append(headerData, checksumData...)
	headerData = append(headerData, header.Checksum)

	return headerData
}

// createDataBlock creates a TAP data block containing the actual code/bytes
func createDataBlock(data []byte) []byte {
	blockLength := uint16(len(data) + 2) // +2 for flag and checksum

	// Create buffer for data block
	dataBlock := make([]byte, 0, int(blockLength)+2) // +2 for block length

	// Write block length
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, blockLength)
	dataBlock = append(dataBlock, buf...)

	// Write flag
	dataBlock = append(dataBlock, dataFlag)

	// Write data
	dataBlock = append(dataBlock, data...)

	// Calculate and append checksum
	checksumData := dataBlock[2:] // Skip block length
	checksum := calculateChecksum(checksumData)
	dataBlock = append(dataBlock, checksum)

	return dataBlock
}

func binaryToTAP(inputPath, outputPath, name string, startAddress uint16) error {
	// Read input file
	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	// Use input filename if no name provided
	if name == "" {
		name = filepath.Base(inputPath)
		name = strings.TrimSuffix(name, filepath.Ext(name))
		if len(name) > 10 {
			name = name[:10]
		}
	}

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer outFile.Close()

	// Create and write header block
	headerBlock := createHeaderBlock(name, uint16(len(inputData)), startAddress)
	if _, err := outFile.Write(headerBlock); err != nil {
		return fmt.Errorf("writing header block: %w", err)
	}

	// Create and write data block
	dataBlock := createDataBlock(inputData)
	if _, err := outFile.Write(dataBlock); err != nil {
		return fmt.Errorf("writing data block: %w", err)
	}

	return nil
}

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

	if err := binaryToTAP(inputFile, outputFile, *name, uint16(*address)); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %s to %s\n", inputFile, outputFile)
}
