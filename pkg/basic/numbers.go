package basic

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Number types in the ZX Spectrum's Sinclair BASIC:
// 1. Small integers (-65535 to +65535) - stored in 6 bytes
// 2. Floating point numbers - stored in 6 bytes using Spectrum's custom format
// 3. Binary numbers (BIN) - stored as small integers

const (
	numberMarker = 0x0E                  // Marks the start of a number in BASIC
	shift31Bits  = float64(2147483648.0) // 2^31, used for mantissa calculation

	// Number format limits
	maxInt   = 65535
	minInt   = -65535
	maxExp   = 126
	minExp   = -129
	mantMask = 0x7F // Mask for mantissa bits in first byte
	signMask = 0x80 // Mask for sign bit
	expBias  = 0x81 // Exponent bias for floating point format
)

// parseNumber tries to parse and encode a number from the text
// Returns the encoded bytes and the number of characters consumed
func (p *Parser) parseNumber(text string) ([]byte, int, error) {
	// Find the end of the number
	i := 0
	hasDecimal := false
	hasExponent := false

	// Look for integer/decimal part
	for i < len(text) && (isDigit(text[i]) || text[i] == '.') {
		if text[i] == '.' {
			if hasDecimal {
				return nil, 0, fmt.Errorf("%s: multiple decimal points in number", p.errorPrefix)
			}
			hasDecimal = true
		}
		i++
	}

	// Look for exponent
	if i < len(text) && (text[i] == 'e' || text[i] == 'E') {
		hasExponent = true
		i++
		if i < len(text) && (text[i] == '+' || text[i] == '-') {
			i++
		}
		for i < len(text) && isDigit(text[i]) {
			i++
		}
	}

	if i == 0 {
		return nil, 0, nil // Not a number
	}

	// Parse the number
	numStr := text[:i]
	val, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: invalid number format: %s", p.errorPrefix, numStr)
	}

	// Check for integer representation
	if !hasDecimal && !hasExponent {
		intVal := int(math.Floor(val))
		if intVal >= minInt && intVal <= maxInt {
			return encodeSmallInt(intVal), i, nil
		}
	}

	// Encode as floating point
	result, err := encodeFloat(val)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", p.errorPrefix, err)
	}
	return result, i, nil
}

// parseBinaryNumber handles BIN format numbers (e.g., BIN 01101)
func (p *Parser) parseBinaryNumber(text string) ([]byte, int, error) {
	if len(text) < 4 || !strings.HasPrefix(text, "BIN") { // Must start with "BIN"
		return nil, 0, nil
	}

	pos := 3
	// Skip spaces after BIN
	for pos < len(text) && text[pos] == ' ' {
		pos++
	}

	if pos >= len(text) {
		return nil, 0, fmt.Errorf("%s: expected binary digits after BIN", p.errorPrefix)
	}

	// Read binary digits
	value := 0
	start := pos
	for pos < len(text) && (text[pos] == '0' || text[pos] == '1') {
		value = value*2 + int(text[pos]-'0')
		if value > maxInt {
			return nil, 0, fmt.Errorf("%s: binary number too large (maximum is %d)", p.errorPrefix, maxInt)
		}
		pos++
	}

	if pos == start {
		return nil, 0, fmt.Errorf("%s: expected binary digits after BIN", p.errorPrefix)
	}

	// Return as small integer
	return encodeSmallInt(value), pos, nil
}

// encodeSmallInt encodes a small integer in Spectrum format
// Format:
// - byte 0: numberMarker (0x0E)
// - byte 1: 0x00 (small integer flag)
// - byte 2: sign (0x00 for positive, 0xFF for negative)
// - byte 3: low byte
// - byte 4: high byte
// - byte 5: 0x00
func encodeSmallInt(val int) []byte {
	// Defensive programming - ensure value is in range
	if val < minInt {
		val = minInt
	} else if val > maxInt {
		val = maxInt
	}

	result := make([]byte, 6)
	result[0] = numberMarker
	result[1] = 0x00

	if val < 0 {
		result[2] = 0xFF
		val = -val
	} else {
		result[2] = 0x00
	}

	result[3] = byte(val & 0xFF)
	result[4] = byte((val >> 8) & 0xFF)
	result[5] = 0x00

	return result
}

// encodeFloat encodes a floating point number in Spectrum format
// Format:
// - byte 0: numberMarker (0x0E)
// - byte 1: exponent (biased by 0x81)
// - byte 2: bits 31-25 of mantissa with sign in bit 7
// - byte 3: bits 24-17 of mantissa
// - byte 4: bits 16-9 of mantissa
// - byte 5: bits 8-1 of mantissa
func encodeFloat(val float64) ([]byte, error) {
	result := make([]byte, 6)
	result[0] = numberMarker

	// Handle zero specially
	if val == 0 {
		return result, nil // All bytes are already 0
	}

	// Handle negative numbers
	sign := byte(0)
	if val < 0 {
		sign = signMask
		val = -val
	}

	// Calculate exponent and mantissa
	exp := math.Floor(math.Log2(val))
	if exp < minExp || exp > maxExp {
		return nil, fmt.Errorf("number out of range (exponent %d not in range %d to %d)", int(exp), minExp, maxExp)
	}

	// Calculate mantissa between 1 and 2
	mantissa := val / math.Pow(2, exp)
	// Convert to integer between 2^31 and 2^32
	mantissa = (mantissa - 1.0) * shift31Bits
	// Round to nearest integer
	mantissaInt := uint32(math.Floor(mantissa + 0.5))

	// Pack the bytes
	result[1] = byte(exp) + expBias // Bias the exponent
	result[2] = byte((mantissaInt>>24)&mantMask) | sign
	result[3] = byte((mantissaInt >> 16) & 0xFF)
	result[4] = byte((mantissaInt >> 8) & 0xFF)
	result[5] = byte(mantissaInt & 0xFF)

	return result, nil
}
