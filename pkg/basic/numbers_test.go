package basic

import (
	"bytes"
	"testing"
)

func TestParseNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []byte
		wantLen int
		wantErr bool
	}{
		{
			name:    "Small integer",
			input:   "123",
			want:    []byte{0x0E, 0x00, 0x00, 123, 0x00, 0x00},
			wantLen: 3,
		},
		{
			name:    "Negative integer",
			input:   "-42",
			want:    []byte{0x0E, 0x00, 0xFF, 42, 0x00, 0x00},
			wantLen: 3,
		},
		{
			name:    "Zero",
			input:   "0",
			want:    []byte{0x0E, 0x00, 0x00, 0x00, 0x00, 0x00},
			wantLen: 1,
		},
		{
			name:    "Double decimal point",
			input:   "1.2.3",
			wantErr: true,
		},
		{
			name:    "Not a number",
			input:   "abc",
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, gotLen, err := p.parseNumber(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if gotLen != tt.wantLen {
				t.Errorf("parseNumber() length = %d, want %d", gotLen, tt.wantLen)
			}

			if tt.want != nil && !bytes.Equal(got, tt.want) {
				t.Errorf("parseNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseBinaryNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []byte
		wantLen int
		wantErr bool
	}{
		{
			name:    "Valid binary",
			input:   "BIN 1010",
			want:    []byte{0x0E, 0x00, 0x00, 10, 0x00, 0x00},
			wantLen: 8,
		},
		{
			name:    "No digits",
			input:   "BIN ",
			wantErr: true,
		},
		{
			name:    "Invalid digit",
			input:   "BIN 102",
			wantLen: 5,
		},
		{
			name:    "Too large",
			input:   "BIN 1111111111111111111",
			wantErr: true,
		},
		{
			name:    "Not BIN",
			input:   "BINARY",
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, gotLen, err := p.parseBinaryNumber(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseBinaryNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if gotLen != tt.wantLen {
				t.Errorf("parseBinaryNumber() length = %d, want %d", gotLen, tt.wantLen)
			}

			if tt.want != nil && !bytes.Equal(got, tt.want) {
				t.Errorf("parseBinaryNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}