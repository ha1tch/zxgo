package basic

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	MaxLineLength = 1024
	MaxStatements = 127 // Maximum number of statements per line
)

// Parser holds the state for parsing BASIC programs
type Parser struct {
	// Configuration
	caseIndependent bool

	// Program type detection
	is48K int // -1=unknown, 0=128K, 1=48K

	// Line tracking
	lineCount      int
	statementCount int
	previousLine   int

	// Statement context
	bracketCount  int
	handlingDEFFN bool
	insideDEFFN   bool
	tokenBracket  bool
	inPrint       bool
	currentParams []int

	// Error reporting
	errorPrefix string
}

// Option defines a parser configuration option
type Option func(*Parser)

// WithCaseIndependent sets case-independent token matching
func WithCaseIndependent(v bool) Option {
	return func(p *Parser) {
		p.caseIndependent = v
	}
}

// NewParser creates a new BASIC parser with the given options
func NewParser(options ...Option) *Parser {
	p := &Parser{
		is48K:         -1,
		previousLine:  -1,
		currentParams: make([]int, 0, 8),
	}
	for _, opt := range options {
		opt(p)
	}
	return p
}

// Public accessor methods

// Is128K returns true if the program is for 128K machines
func (p *Parser) Is128K() bool {
	return p.is48K == 0
}

// LineCount returns the current source line number
func (p *Parser) LineCount() int {
	return p.lineCount
}

// StatementCount returns the current statement number within the line
func (p *Parser) StatementCount() int {
	return p.statementCount
}

// Parse processes BASIC text into binary format suitable for TAP
func (p *Parser) Parse(r io.Reader) ([]byte, error) {
	scanner := bufio.NewScanner(r)
	var output bytes.Buffer

	p.lineCount = 0
	p.previousLine = -1

	for scanner.Scan() {
		p.lineCount++
		line := scanner.Text()

		// Check line length
		if len(line) > MaxLineLength {
			return nil, fmt.Errorf("line %d: exceeds maximum length of %d characters", p.lineCount, MaxLineLength)
		}

		// Skip empty lines and comments
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		p.errorPrefix = fmt.Sprintf("line %d", p.lineCount)

		// Process the line
		basicLine, lineNum, err := p.parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", p.errorPrefix, err)
		}

		// Check line number sequence
		if p.previousLine >= 0 {
			if lineNum < p.previousLine {
				return nil, fmt.Errorf("%s: number %d is smaller than previous line number %d",
					p.errorPrefix, lineNum, p.previousLine)
			}
			if lineNum == p.previousLine {
				fmt.Printf("Warning: Duplicate use of line number %d\n", lineNum)
			}
		}
		p.previousLine = lineNum

		// Write to output buffer
		output.Write(basicLine)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}

	return output.Bytes(), nil
}

// parseLine converts a single line of text into BASIC binary format
func (p *Parser) parseLine(text string) ([]byte, int, error) {
	// Extract line number
	lineNum, rest, err := p.extractLineNumber(text)
	if err != nil {
		return nil, 0, err
	}

	if lineNum < 0 || lineNum > 9999 {
		return nil, 0, fmt.Errorf("line number must be between 0 and 9999")
	}

	// Check for empty line (just a line number)
	if rest == "" {
		return nil, 0, fmt.Errorf("line contains no statements")
	}

	var lineBuf bytes.Buffer

	// Write line number (little-endian)
	lineBuf.WriteByte(byte(lineNum & 0xFF))
	lineBuf.WriteByte(byte(lineNum >> 8))

	// Reserve space for line length (will be filled in later)
	lengthPos := lineBuf.Len()
	lineBuf.WriteByte(0)
	lineBuf.WriteByte(0)

	// Reset state for this line
	p.resetState()

	// Convert the line contents
	if err := p.convertLine(rest, &lineBuf); err != nil {
		return nil, 0, err
	}

	// Check statement count
	if p.statementCount > MaxStatements {
		return nil, 0, fmt.Errorf("too many statements (maximum is %d)", MaxStatements)
	}

	// Check final bracket count
	if p.bracketCount != 0 {
		return nil, 0, fmt.Errorf("mismatched brackets (count: %d)", p.bracketCount)
	}

	// Add end of line marker
	lineBuf.WriteByte(0x0D)

	// Calculate and write line length
	lineBytes := lineBuf.Bytes()
	length := uint16(len(lineBytes) - lengthPos - 2)
	lineBytes[lengthPos] = byte(length & 0xFF)
	lineBytes[lengthPos+1] = byte(length >> 8)

	return lineBytes, lineNum, nil
}

// convertLine processes the content of a BASIC line after the line number
func (p *Parser) convertLine(text string, out *bytes.Buffer) error {
	var inString bool
	var inRem bool
	pos := 0

	expectKeyword := true

	for pos < len(text) {
		// Skip whitespace unless in string or REM
		if !inString && !inRem {
			for pos < len(text) && text[pos] == ' ' {
				pos++
			}
			if pos >= len(text) {
				break
			}
		}

		if inRem {
			// After REM, copy everything as-is, expanding sequences
			for pos < len(text) {
				if match, err := p.expandSequence(text[pos:], false); err != nil {
					return fmt.Errorf("in REM: %w", err)
				} else if match != nil {
					out.Write(match.Bytes)
					pos += match.Length
				} else {
					out.WriteByte(text[pos])
					pos++
				}
			}
			continue
		}

		if inString {
			if text[pos] == '"' {
				inString = false
				out.WriteByte('"')
				pos++
				continue
			}
			// Look for special sequences in strings
			if match, err := p.expandSequence(text[pos:], false); err != nil {
				return fmt.Errorf("in string: %w", err)
			} else if match != nil {
				out.Write(match.Bytes)
				pos += match.Length
			} else {
				out.WriteByte(text[pos])
				pos++
			}
			continue
		}

		// Start of string?
		if text[pos] == '"' {
			inString = true
			out.WriteByte('"')
			pos++
			continue
		}

		// Handle brackets
		if text[pos] == '(' {
			p.bracketCount++
			if p.handlingDEFFN && !p.insideDEFFN {
				p.insideDEFFN = true
			}
			p.tokenBracket = true
			out.WriteByte('(')
			pos++
			continue
		}
		if text[pos] == ')' {
			p.bracketCount--
			if p.bracketCount < 0 {
				return fmt.Errorf("too many closing brackets")
			}
			if p.handlingDEFFN && p.insideDEFFN {
				// Insert room for evaluator (call by value)
				out.WriteByte(0x0E)
				out.WriteByte(0x00)
				out.WriteByte(0x00)
				out.WriteByte(0x00)
				out.WriteByte(0x00)
				out.WriteByte(0x00)
				p.insideDEFFN = false
				p.handlingDEFFN = false
			}
			p.tokenBracket = false
			out.WriteByte(')')
			pos++
			continue
		}

		// Try to match a token
		if match, err := p.matchToken(text[pos:], expectKeyword); err != nil {
			return err
		} else if match != nil {
			// Handle special tokens
			switch match.Value {
			case 0xEA: // REM
				inRem = true
			case 0xCE: // DEF FN
				p.handlingDEFFN = true
			case 0xF5, 0xE0: // PRINT or LPRINT
				p.inPrint = true
			case ':':
				p.statementCount++
				expectKeyword = true
				p.inPrint = false
				// Reset parameter state
				p.currentParams = p.currentParams[:0]
				pos++
				continue
			}

			out.WriteByte(match.Value)
			pos += match.Length

			// After a token, the next token can be any type
			expectKeyword = false
			continue
		}

		// Try to parse a number
		if bytes, consumed, err := p.parseNumber(text[pos:]); err != nil {
			return fmt.Errorf("parsing number: %w", err)
		} else if consumed > 0 {
			out.Write(bytes)
			pos += consumed
			expectKeyword = false
			continue
		}

		// Try to parse a binary number
		if bytes, consumed, err := p.parseBinaryNumber(text[pos:]); err != nil {
			return fmt.Errorf("parsing binary: %w", err)
		} else if consumed > 0 {
			out.Write(bytes)
			pos += consumed
			expectKeyword = false
			continue
		}

		// Look for special sequences
		if match, err := p.expandSequence(text[pos:], true); err != nil {
			return fmt.Errorf("expanding sequence: %w", err)
		} else if match != nil {
			out.Write(match.Bytes)
			pos += match.Length
			continue
		}

		// Just copy any other character
		out.WriteByte(text[pos])
		pos++
		
		// Reset expectKeyword if this wasn't whitespace
		if text[pos-1] != ' ' {
			expectKeyword = false
		}
	}

	return nil
}

// extractLineNumber gets the BASIC line number from the start of the line
func (p *Parser) extractLineNumber(line string) (int, string, error) {
	line = strings.TrimSpace(line)

	// Find first non-digit
	i := 0
	for i < len(line) && isDigit(line[i]) {
		i++
	}

	if i == 0 {
		return 0, "", fmt.Errorf("line must start with a number")
	}

	// Parse line number
	num, err := strconv.Atoi(line[:i])
	if err != nil {
		return 0, "", fmt.Errorf("invalid line number: %s", line[:i])
	}

	// Return line number and rest of line
	rest := strings.TrimSpace(line[i:])
	return num, rest, nil
}

// resetState resets the parser state for a new line
func (p *Parser) resetState() {
	p.statementCount = 0
	p.bracketCount = 0
	p.handlingDEFFN = false
	p.insideDEFFN = false
	p.tokenBracket = false
	p.inPrint = false
	p.currentParams = p.currentParams[:0]
	p.errorPrefix = fmt.Sprintf("line %d", p.lineCount)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}