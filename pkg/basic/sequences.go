package basic

import (
	"fmt"
	"strconv"
	"strings"
)

// SequenceMatch represents an expanded special sequence
type SequenceMatch struct {
	Bytes  []byte  // The encoded bytes
	Length int     // Length of the matched text
}

// expandSequence tries to expand special sequences like {AT}, {INK}, etc.
// Returns nil if no sequence was found at the current position.
func (p *Parser) expandSequence(text string, stripSpaces bool) (*SequenceMatch, error) {
	if !strings.HasPrefix(text, "{") {
		return nil, nil
	}

	// Find closing brace
	end := strings.IndexByte(text, '}')
	if end == -1 {
		if len(text) > 10 { // Only warn if it's not just the start of a longer sequence
			return nil, fmt.Errorf("line %d: unclosed sequence", p.lineCount)
		}
		return nil, nil
	}

	seq := text[1:end] // Remove braces
	if seq == "" {
		return nil, nil
	}

	// Try each sequence type in order
	match, err := p.trySpecialCharacter(seq, end)
	if err != nil {
		return nil, err
	}
	if match != nil {
		return p.adjustSpaces(match, text, stripSpaces)
	}

	match, err = p.tryUDG(seq, end)
	if err != nil {
		return nil, err
	}
	if match != nil {
		return p.adjustSpaces(match, text, stripSpaces)
	}

	match, err = p.tryBlockGraphics(seq, end)
	if err != nil {
		return nil, err
	}
	if match != nil {
		return p.adjustSpaces(match, text, stripSpaces)
	}

	match, err = p.tryHexValue(seq, end)
	if err != nil {
		return nil, err
	}
	if match != nil {
		return p.adjustSpaces(match, text, stripSpaces)
	}

	match, err = p.tryControlSequence(seq, end)
	if err != nil {
		return nil, err
	}
	if match != nil {
		return p.adjustSpaces(match, text, stripSpaces)
	}

	// Not a valid sequence, but don't treat as error
	// Let the parser handle it as normal text
	return nil, nil
}

// trySpecialCharacter handles special character sequences
func (p *Parser) trySpecialCharacter(seq string, end int) (*SequenceMatch, error) {
	switch strings.ToUpper(seq) {
	case "(C)": // Copyright symbol
		return &SequenceMatch{Bytes: []byte{0x7F}, Length: end + 1}, nil
	case "CODE": // Special for OPEN #
		return &SequenceMatch{Bytes: []byte{0xAF}, Length: end + 1}, nil
	case "CAT": // Special for OPEN #
		return &SequenceMatch{Bytes: []byte{0xCF}, Length: end + 1}, nil
	}
	return nil, nil
}

// tryUDG handles User Defined Graphics sequences
func (p *Parser) tryUDG(seq string, end int) (*SequenceMatch, error) {
	if len(seq) != 1 || !isAlpha(seq[0]) {
		return nil, nil
	}

	udg := byte(strings.ToUpper(seq)[0])
	if udg < 'A' || udg > 'U' {
		return nil, nil
	}

	// Handle T and U which are 48K specific
	if udg == 'T' || udg == 'U' {
		switch p.is48K {
		case -1: // unknown
			p.is48K = 1 // mark as 48K
		case 0: // already marked as 128K
			return nil, fmt.Errorf("line %d: UDG '%c' not available in 128K mode", p.lineCount, udg)
		}
	}

	return &SequenceMatch{
		Bytes:  []byte{0x90 + (udg - 'A')},
		Length: end + 1,
	}, nil
}

// tryBlockGraphics handles block graphics sequences
func (p *Parser) tryBlockGraphics(seq string, end int) (*SequenceMatch, error) {
	if len(seq) != 2 {
		return nil, nil
	}

	var prefix byte
	switch seq[0] {
	case '+':
		prefix = 0x88
	case '-':
		prefix = 0x80
	default:
		return nil, nil
	}

	n, err := strconv.Atoi(seq[1:])
	if err != nil || n < 1 || n > 8 {
		return nil, nil
	}

	var value byte
	if prefix == 0x88 {
		value = byte((n % 8) ^ 7)
	} else {
		value = byte(n % 8)
	}

	return &SequenceMatch{
		Bytes:  []byte{prefix + value},
		Length: end + 1,
	}, nil
}

// tryHexValue handles hex value sequences
func (p *Parser) tryHexValue(seq string, end int) (*SequenceMatch, error) {
	if len(seq) != 2 {
		return nil, nil
	}

	val, err := strconv.ParseUint(seq, 16, 8)
	if err != nil {
		return nil, nil
	}

	return &SequenceMatch{
		Bytes:  []byte{byte(val)},
		Length: end + 1,
	}, nil
}

// tryControlSequence handles control code sequences
func (p *Parser) tryControlSequence(seq string, end int) (*SequenceMatch, error) {
	parts := strings.Fields(seq)
	if len(parts) == 0 {
		return nil, nil
	}

	cmd := strings.ToUpper(parts[0])
	var params []int

	// Parse parameters
	if len(parts) > 1 {
		// Clear any existing parameters
		p.currentParams = p.currentParams[:0]

		for _, param := range parts[1:] {
			n, err := strconv.Atoi(param)
			if err != nil {
				return nil, nil // Not a valid control sequence
			}
			// Validate individual parameter
			if err := p.validateControlParam(cmd, n, len(params)); err != nil {
				return nil, fmt.Errorf("line %d: %w", p.lineCount, err)
			}
			params = append(params, n)
			p.currentParams = append(p.currentParams, n)
		}
	}

	// Get control sequence bytes
	result, err := p.getControlBytes(cmd, params)
	if err != nil {
		return nil, fmt.Errorf("line %d: %w", p.lineCount, err)
	}
	if result == nil {
		return nil, nil
	}

	return &SequenceMatch{
		Bytes:  result,
		Length: end + 1,
	}, nil
}

// validateControlParam checks if a parameter is valid for a control code
func (p *Parser) validateControlParam(cmd string, param int, paramIndex int) error {
	switch cmd {
	case "AT":
		switch paramIndex {
		case 0: // Row
			if param < 0 || param > 23 {
				return fmt.Errorf("AT row must be between 0 and 23")
			}
		case 1: // Column
			if param < 0 || param > 31 {
				return fmt.Errorf("AT column must be between 0 and 31")
			}
		default:
			return fmt.Errorf("AT requires exactly 2 parameters")
		}
	case "TAB":
		if param < 0 || param > 31 {
			return fmt.Errorf("TAB position must be between 0 and 31")
		}
	case "INK", "PAPER":
		if param < 0 || param > 7 {
			return fmt.Errorf("%s color must be between 0 and 7", cmd)
		}
	case "FLASH", "BRIGHT", "INVERSE", "OVER":
		if param < 0 || param > 1 {
			return fmt.Errorf("%s must be 0 or 1", cmd)
		}
	default:
		return fmt.Errorf("unknown control sequence: %s", cmd)
	}
	return nil
}

// getControlBytes returns the bytes for a control sequence
func (p *Parser) getControlBytes(cmd string, params []int) ([]byte, error) {
	// Check for AT/TAB in PRINT context
	if (cmd == "AT" || cmd == "TAB") && !p.inPrint {
		return nil, fmt.Errorf("%s only allowed in PRINT statements", cmd)
	}

	switch cmd {
	case "AT":
		if len(params) != 2 {
			return nil, fmt.Errorf("AT requires row and column parameters")
		}
		return []byte{0x16, byte(params[0]), byte(params[1])}, nil
	case "TAB":
		if len(params) != 1 {
			return nil, fmt.Errorf("TAB requires one parameter")
		}
		return []byte{0x17, byte(params[0])}, nil
	case "INK":
		if len(params) != 1 {
			return nil, fmt.Errorf("INK requires one parameter")
		}
		return []byte{0x10, byte(params[0])}, nil
	case "PAPER":
		if len(params) != 1 {
			return nil, fmt.Errorf("PAPER requires one parameter")
		}
		return []byte{0x11, byte(params[0])}, nil
	case "FLASH":
		if len(params) != 1 {
			return nil, fmt.Errorf("FLASH requires one parameter")
		}
		return []byte{0x12, byte(params[0])}, nil
	case "BRIGHT":
		if len(params) != 1 {
			return nil, fmt.Errorf("BRIGHT requires one parameter")
		}
		return []byte{0x13, byte(params[0])}, nil
	case "INVERSE":
		if len(params) != 1 {
			return nil, fmt.Errorf("INVERSE requires one parameter")
		}
		return []byte{0x14, byte(params[0])}, nil
	case "OVER":
		if len(params) != 1 {
			return nil, fmt.Errorf("OVER requires one parameter")
		}
		return []byte{0x15, byte(params[0])}, nil
	}
	return nil, nil
}

// adjustSpaces handles the stripSpaces option by including trailing spaces in Length
func (p *Parser) adjustSpaces(match *SequenceMatch, text string, stripSpaces bool) (*SequenceMatch, error) {
	if stripSpaces {
		// Count trailing spaces
		i := match.Length
		for i < len(text) && text[i] == ' ' {
			i++
		}
		match.Length = i
	}
	return match, nil
}