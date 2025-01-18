package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

// TAPMetadata contains the information from the TAP header block
type TAPMetadata struct {
	Type         byte    // 0=program, 1=number array, 2=character array, 3=code/bytes
	Filename     string  // Up to 10 characters
	DataLength   uint16  // Expected length of the data
	StartAddress uint16  // For code blocks, where to load in memory
	AutoStart    *uint16 // For BASIC programs, which line to auto-start (nil if not applicable)
	VarName      *byte   // For arrays, the variable name (nil if not applicable)
}

// verifyChecksum performs XOR of all bytes including flag byte
func verifyChecksum(data []byte, checksum byte) bool {
	result := byte(0)
	for _, b := range data {
		result ^= b
	}
	return result == checksum
}

// parseHeader extracts metadata from a header block
func parseHeader(headerData []byte) (*TAPMetadata, error) {
	if len(headerData) != 19 { // flag + 17 bytes + checksum
		return nil, fmt.Errorf("invalid header length: %d", len(headerData))
	}

	if headerData[0] != 0x00 { // verify flag
		return nil, fmt.Errorf("invalid header flag: 0x%02x", headerData[0])
	}

	if !verifyChecksum(headerData[:len(headerData)-1], headerData[len(headerData)-1]) {
		return nil, fmt.Errorf("header checksum mismatch")
	}

	blockType := headerData[1]
	if blockType > 3 {
		return nil, fmt.Errorf("invalid block type: %d", blockType)
	}

	filename := strings.TrimRight(string(headerData[2:12]), " ")
	dataLength := binary.LittleEndian.Uint16(headerData[12:14])
	param1 := binary.LittleEndian.Uint16(headerData[14:16])

	metadata := &TAPMetadata{
		Type:       blockType,
		Filename:   filename,
		DataLength: dataLength,
	}

	// Handle parameters based on block type
	switch blockType {
	case 0: // Program
		metadata.AutoStart = &param1
	case 1, 2: // Number or character array
		varName := headerData[15] // param2 high byte contains variable name
		metadata.VarName = &varName
	case 3: // Code
		metadata.StartAddress = param1
	}

	return metadata, nil
}

// LoadTAPResult contains both the loaded data and its metadata
type LoadTAPResult struct {
	Data     []byte
	Metadata *TAPMetadata
}

func loadTap(filename string) (*LoadTAPResult, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	// First block must be a header
	var headerLength uint16
	if err := binary.Read(f, binary.LittleEndian, &headerLength); err != nil {
		return nil, fmt.Errorf("reading header length: %w", err)
	}

	if headerLength != 19 { // flag + 17 bytes + checksum
		return nil, fmt.Errorf("invalid header block length: %d", headerLength)
	}

	headerBlock := make([]byte, headerLength)
	if _, err := io.ReadFull(f, headerBlock); err != nil {
		return nil, fmt.Errorf("reading header block: %w", err)
	}

	metadata, err := parseHeader(headerBlock)
	if err != nil {
		return nil, fmt.Errorf("parsing header: %w", err)
	}

	// Read data block
	var dataBlockLength uint16
	if err := binary.Read(f, binary.LittleEndian, &dataBlockLength); err != nil {
		return nil, fmt.Errorf("reading data block length: %w", err)
	}

	// Data block should be: flag + data + checksum
	if dataBlockLength != metadata.DataLength+2 {
		return nil, fmt.Errorf("data block length mismatch: expected %d, got %d",
			metadata.DataLength+2, dataBlockLength)
	}

	dataBlock := make([]byte, dataBlockLength)
	if _, err := io.ReadFull(f, dataBlock); err != nil {
		return nil, fmt.Errorf("reading data block: %w", err)
	}

	// Verify data block flag and checksum
	if dataBlock[0] != 0xFF {
		return nil, fmt.Errorf("invalid data block flag: 0x%02x", dataBlock[0])
	}

	if !verifyChecksum(dataBlock[:len(dataBlock)-1], dataBlock[len(dataBlock)-1]) {
		return nil, fmt.Errorf("data block checksum mismatch")
	}

	// Extract just the data (without flag and checksum)
	data := dataBlock[1 : len(dataBlock)-1]

	// Verify no more blocks
	var extraLength uint16
	if err := binary.Read(f, binary.LittleEndian, &extraLength); err != io.EOF {
		return nil, fmt.Errorf("unexpected extra data in file")
	}

	return &LoadTAPResult{
		Data:     data,
		Metadata: metadata,
	}, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s file.tap\n", os.Args[0])
		os.Exit(1)
	}

	result, err := loadTap(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded '%s' (%d bytes)\n", result.Metadata.Filename, len(result.Data))
	fmt.Printf("Type: %d\n", result.Metadata.Type)
	
	switch result.Metadata.Type {
	case 0: // Program
		fmt.Printf("Autostart line: %d\n", *result.Metadata.AutoStart)
	case 1, 2: // Arrays
		fmt.Printf("Variable name: %c\n", *result.Metadata.VarName)
	case 3: // Code
		fmt.Printf("Load address: %d\n", result.Metadata.StartAddress)
	}
}