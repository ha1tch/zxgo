package basic

import (
	"testing"
)

func TestMatchToken(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantKeyword bool
		wantToken   byte
		wantLength  int
		wantErr     bool
	}{
		{
			name:        "PRINT keyword",
			input:       "PRINT \"Hello\"",
			wantKeyword: true,
			wantToken:   0xF5,
			wantLength:  5,
		},
		{
			name:        "PRINT in string",
			input:       "PRINT",
			wantKeyword: false,
			wantToken:   0,
			wantLength:  0,
		},
		{
			name:        "REM keyword",
			input:       "REM anything goes",
			wantKeyword: true,
			wantToken:   0xEA,
			wantLength:  3,
		},
		{
			name:        "INT inside PRINT",
			input:       "PRINT",
			wantKeyword: false,
			wantToken:   0,
			wantLength:  0,
		},
		{
			name:        "Case sensitivity",
			input:       "print",
			wantKeyword: true,
			wantToken:   0,
			wantLength:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			match, err := p.matchToken(tt.input, tt.wantKeyword)

			if (err != nil) != tt.wantErr {
				t.Errorf("matchToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantToken == 0 {
				if match != nil {
					t.Errorf("matchToken() = %v, want nil", match)
				}
				return
			}

			if match == nil {
				t.Fatal("matchToken() = nil, want match")
			}

			if match.Value != tt.wantToken {
				t.Errorf("matchToken() token = %#x, want %#x", match.Value, tt.wantToken)
			}

			if match.Length != tt.wantLength {
				t.Errorf("matchToken() length = %d, want %d", match.Length, tt.wantLength)
			}
		})
	}
}