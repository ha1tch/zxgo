package basic

import (
	"fmt"
	"strings"
)

// TokenMatch represents a matched token with its byte value
type TokenMatch struct {
	Value  byte   // The byte value for this token in ZX BASIC
	Length int    // Length of the matched text
}

// matchToken looks for a token at the start of the text
// Returns nil if no token matched
// wantKeyword indicates whether we're expecting a keyword at this position
func (p *Parser) matchToken(text string, wantKeyword bool) (*TokenMatch, error) {
	var longestMatch int
	var matchedToken byte
	found := false

	// Check all tokens from 0xA3 (first keyword) up
	for token := byte(0xA3); token <= 0xFF; token++ {
		tokenDef := TokenMap[token]
		if tokenDef.Text == "" {
			continue
		}

		var matches bool
		if p.caseIndependent {
			matches = strings.HasPrefix(strings.ToUpper(text), strings.ToUpper(tokenDef.Text))
		} else {
			matches = strings.HasPrefix(text, tokenDef.Text)
		}

		if matches {
			length := len(tokenDef.Text)
			if length > longestMatch {
				// Make sure we don't match part of a longer word
				// e.g., "INT" shouldn't match in "PRINT"
				if length < len(text) {
					nextChar := text[length]
					if isAlpha(nextChar) && isAlpha(text[length-1]) {
						continue
					}
				}
				
				longestMatch = length
				matchedToken = token
				found = true
			}
		}
	}

	if !found {
		return nil, nil
	}

	// Get token type and class
	tokenType := TokenMap[matchedToken].Type
	tokenClass := TokenMap[matchedToken].KeywordClass

	// Check if token type matches what we want
	if wantKeyword && tokenType != TokenKeyword && tokenType != TokenColour {
		return nil, nil
	}
	if !wantKeyword && tokenType == TokenKeyword {
		return nil, nil
	}

	// Validate token in current context
	if err := p.validateTokenContext(matchedToken, tokenType, tokenClass); err != nil {
		return nil, err
	}

	// Handle program type detection
	if err := p.handleProgramType(matchedToken); err != nil {
		return nil, err
	}

	return &TokenMatch{
		Value:  matchedToken,
		Length: longestMatch,
	}, nil
}

// validateTokenContext checks if the token is valid in the current context
func (p *Parser) validateTokenContext(token byte, tokenType TokenType, tokenClass KeywordClass) error {
	// Always allow statement separators
	if token == ':' {
		return nil
	}

	// Check for AT/TAB in PRINT context
	if (token == 0xAC || token == 0xAD) && !p.inPrint { // AT or TAB
		return fmt.Errorf("%s: AT/TAB only allowed in PRINT statements", p.errorPrefix)
	}

	// Check keyword class sequence
	if tokenType == TokenKeyword {
		if err := p.validateKeywordClass(token, tokenClass); err != nil {
			return err
		}
	}

	// Validate expression tokens in correct context
	if tokenType == TokenNumExpr || tokenType == TokenStrExpr {
		if p.handlingDEFFN && !p.insideDEFFN {
			return fmt.Errorf("%s: expression not allowed here in DEF FN", p.errorPrefix)
		}
	}

	return nil
}

// validateKeywordClass checks if the keyword's class is valid in current context
func (p *Parser) validateKeywordClass(token byte, class KeywordClass) error {
	if len(class) == 0 {
		return nil // No class restrictions
	}

	// Special cases
	switch token {
	case 0xEA: // REM
		return nil // REM can appear anywhere
	case 0xFA: // IF
		if p.bracketCount > 0 {
			return fmt.Errorf("%s: IF not allowed within brackets", p.errorPrefix)
		}
	case 0xF1: // LET
		if p.bracketCount > 0 {
			return fmt.Errorf("%s: LET not allowed within brackets", p.errorPrefix)
		}
	}

	// Check first class requirement
	firstClass := class[0]
	switch firstClass {
	case ClassNone:
		return nil
	case ClassLet: // Variable required
		if p.tokenBracket {
			return fmt.Errorf("%s: variable required here", p.errorPrefix)
		}
	case ClassLetExpr: // Expression required
		if !p.tokenBracket && p.bracketCount == 0 {
			return fmt.Errorf("%s: expression required here", p.errorPrefix)
		}
	case ClassVarChar: // Single character variable required
		if p.tokenBracket {
			return fmt.Errorf("%s: single character variable required here", p.errorPrefix)
		}
	}

	return nil
}

// handleProgramType updates program type flags based on tokens found
func (p *Parser) handleProgramType(token byte) error {
	// Check for 128K-specific tokens
	if token == 0xA3 || token == 0xA4 { // SPECTRUM or PLAY
		switch p.is48K {
		case -1: // unknown
			p.is48K = 0 // mark as 128K
		case 1: // already marked as 48K
			return fmt.Errorf("%s: program contains 128K keywords but was already marked as 48K", p.errorPrefix)
		}
	}

	// Check for 48K-specific features in UDGs
	if token >= 0x90 && token <= 0xA2 { // UDG range
		udgT := byte(0x90 + ('T' - 'A'))
		udgU := byte(0x90 + ('U' - 'A'))
		if token == udgT || token == udgU {
			switch p.is48K {
			case -1: // unknown
				p.is48K = 1 // mark as 48K
			case 0: // already marked as 128K
				return fmt.Errorf("%s: program contains 48K UDGs but was already marked as 128K", p.errorPrefix)
			}
		}
	}

	return nil
}

// isAlpha returns true if the byte is an ASCII letter
func isAlpha(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

// isSpace returns true if the byte is a space character
func isSpace(c byte) bool {
	return c == ' ' || c == '\t'
}

// skipSpaces advances past any whitespace
func skipSpaces(text string) string {
	for i := 0; i < len(text); i++ {
		if !isSpace(text[i]) {
			return text[i:]
		}
	}
	return ""
}