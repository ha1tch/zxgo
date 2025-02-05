package basic

import (
	"bytes"
	"testing"
)

func TestExpandSequence(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		stripSpaces bool
		want        []byte
		wantLength  int
		wantErr     bool
	}{
		{
			name:       "Copyright symbol",
			input:      "{(C)}",
			want:       []byte{0x7F},
			wantLength: 5,
		},
		{
			name:       "UDG",
			input:      "{A}",
			want:       []byte{0x90},
			wantLength: 3,
		},
		{
			name:       "UDG T (48K)",
			input:      "{T}",
			want:       []byte{0x90 + ('T' - 'A')},
			wantLength: 3,
		},
		{
			name:       "Block graphics +",
			input:      "{+1}",
			want:       []byte{0x88},
			wantLength: 4,
		},
		{
			name:       "Block graphics -",
			input:      "{-1}",
			want:       []byte{0x80},
			wantLength: 4,
		},
		{
			name:       "Hex value",
			input:      "{7F}",
			want:       []byte{0x7F},
			wantLength: 4,
		},
		{
			name:        "AT with spaces",
			input:       "{AT 10 20}  ",
			stripSpaces: true,
			want:        []byte{0x16, 10, 20},
			wantLength:  12,
		},
		{
			name:    "Invalid sequence",
			input:   "{INVALID}",
			wantErr: false, // Should return nil, not error
		},
		{
			name:    "AT without PRINT",
			input:   "{AT 0 0}",
			wantErr: true,
		},
		{
			name:    "Unclosed sequence",
			input:   "{UNCLOSED",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			if tt.input == "{AT 0 0}" {
				p.inPrint = true
			}

			got, err := p.expandSequence(tt.input, tt.stripSpaces)

			if (err != nil) != tt.wantErr {
				t.Errorf("expandSequence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if tt.want == nil {
				if got != nil {
					t.Errorf("expandSequence() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatal("expandSequence() = nil, want match")
			}

			if !bytes.Equal(got.Bytes, tt.want) {
				t.Errorf("expandSequence() bytes = %v, want %v", got.Bytes, tt.want)
			}

			if got.Length != tt.wantLength {
				t.Errorf("expandSequence() length = %d, want %d", got.Length, tt.wantLength)
			}
		})
	}
}

func TestControlSequences(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []byte
		wantErr bool
	}{
		{
			name:  "INK 2",
			input: "{INK 2}",
			want:  []byte{0x10, 2},
		},
		{
			name:    "INK out of range",
			input:   "{INK 8}",
			wantErr: true,
		},
		{
			name:  "FLASH 1",
			input: "{FLASH 1}",
			want:  []byte{0x12, 1},
		},
		{
			name:    "FLASH invalid",
			input:   "{FLASH 2}",
			wantErr: true,
		},
		{
			name:  "BRIGHT 0",
			input: "{BRIGHT 0}",
			want:  []byte{0x13, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, err := p.expandSequence(tt.input, false)

			if (err != nil) != tt.wantErr {
				t.Errorf("expandSequence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if !bytes.Equal(got.Bytes, tt.want) {
				t.Errorf("expandSequence() = %v, want %v", got.Bytes, tt.want)
			}
		})
	}
}