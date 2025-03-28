package h264_test

import (
	"bytes"
	"testing"

	"github.com/tachode/bitstream-go/h264"
)

func TestUnescape(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "No escape sequences",
			input:    []byte{0x00, 0x01, 0x02, 0x03},
			expected: []byte{0x00, 0x01, 0x02, 0x03},
		},
		{
			name:     "Single escape sequence",
			input:    []byte{0x00, 0x00, 0x03, 0x04},
			expected: []byte{0x00, 0x00, 0x04},
		},
		{
			name:     "Multiple escape sequences",
			input:    []byte{0x00, 0x00, 0x03, 0x04, 0x00, 0x00, 0x03, 0x05},
			expected: []byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x05},
		},
		{
			name:     "Escape sequence at the end",
			input:    []byte{0x00, 0x00, 0x03},
			expected: []byte{0x00, 0x00},
		},
		{
			name:     "00 00 03 in the output",
			input:    []byte{0x00, 0x00, 0x03, 0x03, 0x00, 0x00, 0x03, 0x05},
			expected: []byte{0x00, 0x00, 0x03, 0x00, 0x00, 0x05},
		},
		{
			name:     "Empty input",
			input:    []byte{},
			expected: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := h264.Unescape(tt.input)
			if !bytes.Equal(output, tt.expected) {
				t.Errorf("Unescape(%v) = %v, want %v", tt.input, output, tt.expected)
			}
		})
	}
}
