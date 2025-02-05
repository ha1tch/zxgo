package basic

// TokenType identifies what kind of token this is
type TokenType byte

const (
	TokenNormal TokenType = iota  // No special meaning
	TokenKeyword                  // Always keyword
	TokenColour                   // Can be both keyword and non-keyword (colour parameters)
	TokenNumExpr                  // Numeric expression token
	TokenStrExpr                  // String expression token
	TokenPrint                    // May only appear in (L)PRINT statements (AT and TAB)
	TokenTypeless                 // Type-less (normal ASCII or expression token)
)

// KeywordClass defines what follows each token in Spectrum BASIC
// These match the classes in the Spectrum ROM
type KeywordClass []byte

const (
	// Standard classes from ROM
	ClassNone      byte = 0   // No further operands
	ClassLet       byte = 1   // Used in LET. A variable is required
	ClassLetExpr   byte = 2   // Used in LET. An expression must follow
	ClassExprOpt   byte = 3   // Optional numeric expression
	ClassVarChar   byte = 4   // Single character variable must follow
	ClassItems     byte = 5   // Set of items may be given
	ClassNumExpr   byte = 6   // Numeric expression must follow
	ClassColour    byte = 7   // Handles colour items
	ClassTwoNum    byte = 8   // Two numeric expressions with comma
	ClassTwoNumCol byte = 9   // Like Class8 but colour items may precede
	ClassStrExpr   byte = 10  // String expression must follow
	ClassTape      byte = 11  // Handles cassette routines

	// Additional classes needed
	ClassStrList   byte = 12  // One or more string expressions with commas
	ClassExprList  byte = 13  // One or more expressions with commas
	ClassVarList   byte = 14  // One or more variables with commas (READ)
	ClassDEFFN     byte = 15  // DEF FN form
)

// Token represents a ZX Spectrum BASIC token
type Token struct {
	Text         string        // The token text or character
	Type         TokenType     // Type of token
	KeywordClass KeywordClass  // Classes that define what can follow
}