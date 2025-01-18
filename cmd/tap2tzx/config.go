// config.go
package main

import (
	"fmt"
	"io"
	"os"
)

type blockPosition struct {
	id       string
	position int64
	size     int64
}

// calculateBlockSize calculates the total size of a block in bytes
func calculateBlockSize(block blockConfig) (int64, error) {
	var size int64

	// Add size for group block if present
	if block.Group != "" {
		size += 2 + int64(len(block.Group)) // ID + length byte + string
	}

	// Add size for description if present
	if block.Desc != "" {
		size += 2 + int64(len(block.Desc)) // ID + length byte + string
	}

	// Add size for jump block if present
	if block.JumpTo != "" {
		size += 3 // ID + 2 bytes for offset
	}

	// Add size for loop start if present
	if block.LoopStart > 0 {
		size += 3 // ID + 2 bytes for repetitions
	}

	// Add size for loop end if present
	if block.LoopEnd {
		size += 1 // Just ID byte
	}

	// Add size for TAP data if present
	if block.File != "" {
		info, err := os.Stat(block.File)
		if err != nil {
			return 0, fmt.Errorf("getting file size: %w", err)
		}
		// For each TAP block we need:
		// - ID byte (1)
		// - Pause duration (2)
		// - Length of TAP data (2)
		// - The TAP data itself
		size += 5 + info.Size()
	}

	return size, nil
}

// calculateBlockPositions performs a dry run to calculate positions of all blocks
func calculateBlockPositions(config *tzxConfig) (map[string]int64, error) {
	positions := make(map[string]int64)
	currentPos := int64(0)

	// Skip past TZX header
	currentPos += 10 // Signature(7) + EOF(1) + Version(2)

	// Skip past archive info if present
	if config.Metadata.Title != "" || config.Metadata.Author != "" || config.Metadata.Year != "" {
		archiveSize := int64(3) // Block ID + number of strings
		if config.Metadata.Title != "" {
			archiveSize += 2 + int64(len(config.Metadata.Title))
		}
		if config.Metadata.Author != "" {
			archiveSize += 2 + int64(len(config.Metadata.Author))
		}
		if config.Metadata.Year != "" {
			archiveSize += 2 + int64(len(config.Metadata.Year))
		}
		currentPos += archiveSize
	}

	// Skip past hardware info if present
	if config.Hardware.K128Only || config.Hardware.UseAY || config.Hardware.Model != "" {
		hardwareSize := int64(2) // Block ID + number of entries
		if config.Hardware.K128Only {
			hardwareSize += 6 // Two entries of 3 bytes each
		}
		if config.Hardware.UseAY {
			hardwareSize += 3
		}
		if config.Hardware.Model != "" {
			hardwareSize += 3
		}
		currentPos += hardwareSize
	}

	// Calculate positions for all blocks
	for _, block := range config.Blocks {
		if block.ID != "" {
			positions[block.ID] = currentPos
		}

		size, err := calculateBlockSize(block)
		if err != nil {
			return nil, fmt.Errorf("calculating block size: %w", err)
		}
		currentPos += size
	}

	return positions, nil
}

// processConfigBlock processes a single block from the config
func processConfigBlock(w io.Writer, block blockConfig, positions map[string]int64, currentPos *int64) error {
	if block.Group != "" {
		if err := writeGroupStart(w, block.Group); err != nil {
			return fmt.Errorf("writing group start: %w", err)
		}
		*currentPos += 2 + int64(len(block.Group))
	}

	if block.Desc != "" {
		if err := writeTextDescription(w, block.Desc); err != nil {
			return fmt.Errorf("writing description: %w", err)
		}
		*currentPos += 2 + int64(len(block.Desc))
	}

	if block.JumpTo != "" {
		targetPos, ok := positions[block.JumpTo]
		if !ok {
			return fmt.Errorf("jump target '%s' not found", block.JumpTo)
		}

		// Calculate relative jump
		offset := targetPos - (*currentPos + 3) // +3 for jump block size
		if offset > 32767 || offset < -32768 {
			return fmt.Errorf("jump offset %d out of range [-32768, 32767]", offset)
		}

		if err := writeJumpBlock(w, int16(offset)); err != nil {
			return fmt.Errorf("writing jump block: %w", err)
		}
		*currentPos += 3
	}

	if block.LoopStart > 0 {
		if err := writeLoopStart(w, uint16(block.LoopStart)); err != nil {
			return fmt.Errorf("writing loop start: %w", err)
		}
		*currentPos += 3
	}

	if block.File != "" {
		inFile, err := os.Open(block.File)
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

			if err := writeStandardSpeedBlock(w, data, 1000); err != nil {
				return fmt.Errorf("writing TAP block: %w", err)
			}
			*currentPos += 5 + int64(len(data))
		}
	}

	if block.LoopEnd {
		if err := writeLoopEnd(w); err != nil {
			return fmt.Errorf("writing loop end: %w", err)
		}
		*currentPos += 1
	}

	return nil
}

// processConfiguredBlocks processes all blocks from the config file
func processConfiguredBlocks(w io.Writer, config *tzxConfig) error {
	// First pass: calculate all positions
	positions, err := calculateBlockPositions(config)
	if err != nil {
		return fmt.Errorf("calculating positions: %w", err)
	}

	// Second pass: write all blocks
	currentPos := int64(10) // Start after TZX header

	// Skip past metadata and hardware info blocks
	if config.Metadata.Title != "" || config.Metadata.Author != "" || config.Metadata.Year != "" {
		archiveSize := int64(3)
		if config.Metadata.Title != "" {
			archiveSize += 2 + int64(len(config.Metadata.Title))
		}
		if config.Metadata.Author != "" {
			archiveSize += 2 + int64(len(config.Metadata.Author))
		}
		if config.Metadata.Year != "" {
			archiveSize += 2 + int64(len(config.Metadata.Year))
		}
		currentPos += archiveSize
	}

	if config.Hardware.K128Only || config.Hardware.UseAY || config.Hardware.Model != "" {
		hardwareSize := int64(2)
		if config.Hardware.K128Only {
			hardwareSize += 6
		}
		if config.Hardware.UseAY {
			hardwareSize += 3
		}
		if config.Hardware.Model != "" {
			hardwareSize += 3
		}
		currentPos += hardwareSize
	}

	// Write all blocks
	for _, block := range config.Blocks {
		if err := processConfigBlock(w, block, positions, &currentPos); err != nil {
			return fmt.Errorf("processing block: %w", err)
		}
	}

	return nil
}