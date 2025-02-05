package tap

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	HeaderLength = 0x13
	HeaderFlag   = 0x00
	DataFlag     = 0xFF

	// Block types
	Program = 0x00
	Data    = 0x01
	Chars   = 0x02
	Bytes   = 0x03
)

// Header represents a TAP file header block
type Header struct {
	BlockLength uint16
	Flag        byte
	Type        byte
	Filename    [10]byte
	DataLength  uint16
	Param1      uint16 // Start address for bytes, autostart line for program
	Param2      uint16 // Program length for program, 32768 for bytes
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

// createHeaderBlock creates a TAP header block
func createHeaderBlock(blockType byte, filename string, dataLength uint16, param1, param2 uint16) []byte {
	header := Header{
		BlockLength: HeaderLength,
		Flag:        HeaderFlag,
		Type:        blockType,
		DataLength:  dataLength,
		Param1:      param1,
		Param2:      param2,
	}

	// Pad filename with spaces
	copy(header.Filename[:], []byte(filename))
	for i := len(filename); i < 10; i++ {
		header.Filename[i] = ' '
	}

	// Create buffer for header data
	headerData := make([]byte, 0, HeaderLength+2) // +2 for block length

	// Write block length
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, header.BlockLength)
	headerData = append(headerData, buf...)

	// Build header data for checksum calculation
	checksumData := make([]byte, 0, HeaderLength-1) // -1 as checksum isn't included
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

// createDataBlock creates a TAP data block
func createDataBlock(data []byte) []byte {
	blockLength := uint16(len(data) + 2) // +2 for flag and checksum

	// Create buffer for data block
	dataBlock := make([]byte, 0, int(blockLength)+2) // +2 for block length

	// Write block length
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, blockLength)
	dataBlock = append(dataBlock, buf...)

	// Write flag
	dataBlock = append(dataBlock, DataFlag)

	// Write data
	dataBlock = append(dataBlock, data...)

	// Calculate and append checksum
	checksumData := dataBlock[2:] // Skip block length
	checksum := calculateChecksum(checksumData)
	dataBlock = append(dataBlock, checksum)

	return dataBlock
}

// BinaryToTAP converts a binary file to TAP format
func BinaryToTAP(inputPath, outputPath, name string, startAddress uint16) error {
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

	// Create and write header block for bytes
	headerBlock := createHeaderBlock(Bytes, name, uint16(len(inputData)), startAddress, 0)
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

// WriteBasicToTAP writes a BASIC program to TAP format
func WriteBasicToTAP(w io.Writer, name string, data []byte, autostart uint16) error {
	// Create and write header block for BASIC program
	headerBlock := createHeaderBlock(Program, name, uint16(len(data)), autostart, uint16(len(data)))
	if _, err := w.Write(headerBlock); err != nil {
		return fmt.Errorf("writing header block: %w", err)
	}

	// Create and write data block
	dataBlock := createDataBlock(data)
	if _, err := w.Write(dataBlock); err != nil {
		return fmt.Errorf("writing data block: %w", err)
	}

	return nil
}