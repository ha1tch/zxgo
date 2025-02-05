package basic

import "fmt"

// TokenMap contains all ZX Spectrum BASIC tokens and their properties
var TokenMap = [256]Token{
	// ASCII control characters (0-31)
	0:  {Text: "", Type: TokenTypeless},
	1:  {Text: "", Type: TokenTypeless},
	2:  {Text: "", Type: TokenTypeless},
	3:  {Text: "", Type: TokenTypeless},
	4:  {Text: "", Type: TokenTypeless},
	5:  {Text: "", Type: TokenTypeless},
	6:  {Text: "", Type: TokenTypeless, KeywordClass: KeywordClass{0}}, // Print '
	7:  {Text: "", Type: TokenTypeless},
	8:  {Text: "", Type: TokenTypeless},
	9:  {Text: "", Type: TokenTypeless},
	10: {Text: "", Type: TokenTypeless},
	11: {Text: "", Type: TokenTypeless},
	12: {Text: "", Type: TokenTypeless},
	13: {Text: "(eoln)", Type: TokenTypeless}, // CR
	14: {Text: "", Type: TokenTypeless},       // Number
	15: {Text: "", Type: TokenTypeless},
	16: {Text: "", Type: TokenTypeless}, // INK
	17: {Text: "", Type: TokenTypeless}, // PAPER
	18: {Text: "", Type: TokenTypeless}, // FLASH
	19: {Text: "", Type: TokenTypeless}, // BRIGHT
	20: {Text: "", Type: TokenTypeless}, // INVERSE
	21: {Text: "", Type: TokenTypeless}, // OVER
	22: {Text: "", Type: TokenTypeless}, // AT
	23: {Text: "", Type: TokenTypeless}, // TAB
	24: {Text: "", Type: TokenTypeless},
	25: {Text: "", Type: TokenTypeless},
	26: {Text: "", Type: TokenTypeless},
	27: {Text: "", Type: TokenTypeless},
	28: {Text: "", Type: TokenTypeless},
	29: {Text: "", Type: TokenTypeless},
	30: {Text: "", Type: TokenTypeless},
	31: {Text: "", Type: TokenTypeless},

	// ASCII printable characters (32-127) are initialized in init()

	// 128K BASIC tokens
	0xA3: {Text: "SPECTRUM", Type: TokenKeyword, KeywordClass: KeywordClass{ClassNone}},
	0xA4: {Text: "PLAY", Type: TokenKeyword, KeywordClass: KeywordClass{ClassStrList}},

	// Expression tokens
	0xA5: {Text: "RND", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassNone}},
	0xA6: {Text: "INKEY$", Type: TokenStrExpr, KeywordClass: KeywordClass{ClassNone}},
	0xA7: {Text: "PI", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassNone}},
	0xA8: {Text: "FN", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassLet, '(', ClassExprList, ')', 0}},
	0xA9: {Text: "POINT", Type: TokenNumExpr, KeywordClass: KeywordClass{'(', ClassTwoNum, ')', 0}},
	0xAA: {Text: "SCREEN$", Type: TokenStrExpr, KeywordClass: KeywordClass{'(', ClassTwoNum, ')', 0}},
	0xAB: {Text: "ATTR", Type: TokenNumExpr, KeywordClass: KeywordClass{'(', ClassTwoNum, ')', 0}},
	0xAC: {Text: "AT", Type: TokenPrint, KeywordClass: KeywordClass{ClassTwoNum}},
	0xAD: {Text: "TAB", Type: TokenPrint, KeywordClass: KeywordClass{ClassNumExpr}},
	0xAE: {Text: "VAL$", Type: TokenStrExpr, KeywordClass: KeywordClass{ClassStrExpr}},
	0xAF: {Text: "CODE", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassStrExpr}},
	0xB0: {Text: "VAL", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassStrExpr}},
	0xB1: {Text: "LEN", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassStrExpr}},
	0xB2: {Text: "SIN", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassNumExpr}},
	0xB3: {Text: "COS", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassNumExpr}},
	0xB4: {Text: "TAN", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassNumExpr}},
	0xB5: {Text: "ASN", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassNumExpr}},
	0xB6: {Text: "ACS", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassNumExpr}},
	0xB7: {Text: "ATN", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassNumExpr}},
	0xB8: {Text: "LN", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassNumExpr}},
	0xB9: {Text: "EXP", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassNumExpr}},
	0xBA: {Text: "INT", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassNumExpr}},
	0xBB: {Text: "SQR", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassNumExpr}},
	0xBC: {Text: "SGN", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassNumExpr}},
	0xBD: {Text: "ABS", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassNumExpr}},
	0xBE: {Text: "PEEK", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassNumExpr}},
	0xBF: {Text: "IN", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassNumExpr}},
	0xC0: {Text: "USR", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassNumExpr}},
	0xC1: {Text: "STR$", Type: TokenStrExpr, KeywordClass: KeywordClass{ClassNumExpr}},
	0xC2: {Text: "CHR$", Type: TokenStrExpr, KeywordClass: KeywordClass{ClassNumExpr}},
	0xC3: {Text: "NOT", Type: TokenNumExpr, KeywordClass: KeywordClass{ClassNumExpr}},
	0xC4: {Text: "BIN", Type: TokenTypeless},
	0xC5: {Text: "OR", Type: TokenTypeless, KeywordClass: KeywordClass{ClassItems}},
	0xC6: {Text: "AND", Type: TokenTypeless, KeywordClass: KeywordClass{ClassItems}},
	0xC7: {Text: "<=", Type: TokenTypeless, KeywordClass: KeywordClass{ClassItems}},
	0xC8: {Text: ">=", Type: TokenTypeless, KeywordClass: KeywordClass{ClassItems}},
	0xC9: {Text: "<>", Type: TokenTypeless, KeywordClass: KeywordClass{ClassItems}},
	0xCA: {Text: "LINE", Type: TokenTypeless},
	0xCB: {Text: "THEN", Type: TokenTypeless},
	0xCC: {Text: "TO", Type: TokenTypeless},
	0xCD: {Text: "STEP", Type: TokenTypeless},
	0xCE: {Text: "DEF FN", Type: TokenKeyword, KeywordClass: KeywordClass{ClassDEFFN}},
	0xCF: {Text: "CAT", Type: TokenKeyword, KeywordClass: KeywordClass{ClassTape}},
	0xD0: {Text: "FORMAT", Type: TokenKeyword, KeywordClass: KeywordClass{ClassTape}},
	0xD1: {Text: "MOVE", Type: TokenKeyword, KeywordClass: KeywordClass{ClassTape}},
	0xD2: {Text: "ERASE", Type: TokenKeyword, KeywordClass: KeywordClass{ClassTape}},
	0xD3: {Text: "OPEN #", Type: TokenKeyword, KeywordClass: KeywordClass{ClassTape}},
	0xD4: {Text: "CLOSE #", Type: TokenKeyword, KeywordClass: KeywordClass{ClassTape}},
	0xD5: {Text: "MERGE", Type: TokenKeyword, KeywordClass: KeywordClass{ClassTape}},
	0xD6: {Text: "VERIFY", Type: TokenKeyword, KeywordClass: KeywordClass{ClassTape}},
	0xD7: {Text: "BEEP", Type: TokenKeyword, KeywordClass: KeywordClass{ClassTwoNum}},
	0xD8: {Text: "CIRCLE", Type: TokenKeyword, KeywordClass: KeywordClass{ClassTwoNumCol, ',', ClassNumExpr}},
	0xD9: {Text: "INK", Type: TokenColour, KeywordClass: KeywordClass{ClassColour}},
	0xDA: {Text: "PAPER", Type: TokenColour, KeywordClass: KeywordClass{ClassColour}},
	0xDB: {Text: "FLASH", Type: TokenColour, KeywordClass: KeywordClass{ClassColour}},
	0xDC: {Text: "BRIGHT", Type: TokenColour, KeywordClass: KeywordClass{ClassColour}},
	0xDD: {Text: "INVERSE", Type: TokenColour, KeywordClass: KeywordClass{ClassColour}},
	0xDE: {Text: "OVER", Type: TokenColour, KeywordClass: KeywordClass{ClassColour}},
	0xDF: {Text: "OUT", Type: TokenKeyword, KeywordClass: KeywordClass{ClassTwoNum}},
	0xE0: {Text: "LPRINT", Type: TokenKeyword, KeywordClass: KeywordClass{ClassItems}},
	0xE1: {Text: "LLIST", Type: TokenKeyword, KeywordClass: KeywordClass{ClassExprOpt}},
	0xE2: {Text: "STOP", Type: TokenKeyword, KeywordClass: KeywordClass{ClassNone}},
	0xE3: {Text: "READ", Type: TokenKeyword, KeywordClass: KeywordClass{ClassVarList}},
	0xE4: {Text: "DATA", Type: TokenColour, KeywordClass: KeywordClass{ClassExprList}},
	0xE5: {Text: "RESTORE", Type: TokenKeyword, KeywordClass: KeywordClass{ClassExprOpt}},
	0xE6: {Text: "NEW", Type: TokenKeyword, KeywordClass: KeywordClass{ClassNone}},
	0xE7: {Text: "BORDER", Type: TokenKeyword, KeywordClass: KeywordClass{ClassNumExpr}},
	0xE8: {Text: "CONTINUE", Type: TokenKeyword, KeywordClass: KeywordClass{ClassNone}},
	0xE9: {Text: "DIM", Type: TokenKeyword, KeywordClass: KeywordClass{ClassLet, '(', ClassExprList, ')', 0}},
	0xEA: {Text: "REM", Type: TokenKeyword, KeywordClass: KeywordClass{ClassItems}},
	0xEB: {Text: "FOR", Type: TokenKeyword, KeywordClass: KeywordClass{ClassVarChar, '=', ClassNumExpr, 0xCC, ClassNumExpr, 0xCD, ClassNumExpr}},
	0xEC: {Text: "GO TO", Type: TokenKeyword, KeywordClass: KeywordClass{ClassNumExpr}},
	0xED: {Text: "GO SUB", Type: TokenKeyword, KeywordClass: KeywordClass{ClassNumExpr}},
	0xEE: {Text: "INPUT", Type: TokenKeyword, KeywordClass: KeywordClass{ClassItems}},
	0xEF: {Text: "LOAD", Type: TokenKeyword, KeywordClass: KeywordClass{ClassTape}},
	0xF0: {Text: "LIST", Type: TokenKeyword, KeywordClass: KeywordClass{ClassExprOpt}},
	0xF1: {Text: "LET", Type: TokenKeyword, KeywordClass: KeywordClass{ClassLet, '=', ClassLetExpr}},
	0xF2: {Text: "PAUSE", Type: TokenKeyword, KeywordClass: KeywordClass{ClassNumExpr}},
	0xF3: {Text: "NEXT", Type: TokenKeyword, KeywordClass: KeywordClass{ClassVarChar}},
	0xF4: {Text: "POKE", Type: TokenKeyword, KeywordClass: KeywordClass{ClassTwoNum}},
	0xF5: {Text: "PRINT", Type: TokenKeyword, KeywordClass: KeywordClass{ClassItems}},
	0xF6: {Text: "PLOT", Type: TokenKeyword, KeywordClass: KeywordClass{ClassTwoNumCol}},
	0xF7: {Text: "RUN", Type: TokenKeyword, KeywordClass: KeywordClass{ClassExprOpt}},
	0xF8: {Text: "SAVE", Type: TokenKeyword, KeywordClass: KeywordClass{ClassTape}},
	0xF9: {Text: "RANDOMIZE", Type: TokenKeyword, KeywordClass: KeywordClass{ClassExprOpt}},
	0xFA: {Text: "IF", Type: TokenKeyword, KeywordClass: KeywordClass{ClassNumExpr, 0xCB}},
	0xFB: {Text: "CLS", Type: TokenKeyword, KeywordClass: KeywordClass{ClassNone}},
	0xFC: {Text: "DRAW", Type: TokenKeyword, KeywordClass: KeywordClass{ClassTwoNumCol, ',', ClassNumExpr}},
	0xFD: {Text: "CLEAR", Type: TokenKeyword, KeywordClass: KeywordClass{ClassExprOpt}},
	0xFE: {Text: "RETURN", Type: TokenKeyword, KeywordClass: KeywordClass{ClassNone}},
	0xFF: {Text: "COPY", Type: TokenKeyword, KeywordClass: KeywordClass{ClassNone}},
}

func init() {
	// Initialize ASCII printable characters (32-127)
	for i := byte(32); i < 128; i++ {
		TokenMap[i] = Token{
			Text: string([]byte{i}),
			Type: TokenTypeless,
		}
	}

	// Special handling for colon
	TokenMap[':'] = Token{Text: ":", Type: TokenKeyword, KeywordClass: KeywordClass{0}}

	// Block graphics without shift (0x80-0x87)
	for i := byte(0x80); i <= 0x87; i++ {
		TokenMap[i] = Token{
			Text: fmt.Sprintf("{-%d}", i-0x7F),
			Type: TokenTypeless,
		}
	}

	// Block graphics with shift (0x88-0x8F)
	for i := byte(0x88); i <= 0x8F; i++ {
		TokenMap[i] = Token{
			Text: fmt.Sprintf("{+%d}", i-0x87),
			Type: TokenTypeless,
		}
	}

	// UDGs (0x90-0xA2)
	for i := byte(0x90); i <= 0xA2; i++ {
		TokenMap[i] = Token{
			Text: fmt.Sprintf("{%c}", 'A'+i-0x90),
			Type: TokenTypeless,
		}
	}
}
