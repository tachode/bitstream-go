package bits_test

import (
	"bytes"
	"testing"

	"github.com/tachode/bitstream-go/bits"
)

func TestWriteBuffer_Write(t *testing.T) {
	var buf bytes.Buffer
	writer := &bits.WriteBuffer{Writer: &buf}

	tests := []struct {
		value    any
		bits     int
		expected []byte
	}{
		{value: uint8(0b101), bits: 3, expected: []byte{0b10100000}},
		{value: uint8(0b1101), bits: 4, expected: []byte{0b10111010}},
		{value: uint8(0b1111_1111), bits: 8, expected: []byte{0b10111011, 0b11111110}},
		{value: uint8(0b1), bits: 1, expected: []byte{0b10111011, 0b11111111}},
		{
			value: uint64(0xFACEFACEFACEFACE), bits: 64,
			expected: []byte{0b10111011, 0b11111111, 0xFA, 0xCE, 0xFA, 0xCE, 0xFA, 0xCE, 0xFA, 0xCE},
		},
	}

	for i, tt := range tests {
		buf.Reset()
		for j := 0; j <= i; j++ {
			err := writer.WriteBits(tests[j].value, tests[j].bits)
			if err != nil {
				t.Errorf("Write(0b%b, %d) returned error: %v", tests[j].value, tests[j].bits, err)
			}
		}
		err := writer.Flush()
		if err != nil {
			t.Errorf("Flush() returned error: %v", err)
		}
		if !bytes.Equal(buf.Bytes(), tt.expected) {
			t.Errorf("Write(0b%b, %d) = %08b, want %08b", tt.value, tt.bits, buf.Bytes(), tt.expected)
		}
	}
}

func TestWriteBuffer_Flush(t *testing.T) {
	var buf bytes.Buffer
	writer := &bits.WriteBuffer{Writer: &buf}

	err := writer.WriteBits(uint8(0b101), 3)
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	err = writer.Flush()
	if err != nil {
		t.Fatalf("Flush returned error: %v", err)
	}

	expected := []byte{0b10100000}
	if !bytes.Equal(buf.Bytes(), expected) {
		t.Errorf("Flush() = %08b, want %08b", buf.Bytes(), expected)
	}
}

func TestWriteBuffer_WriteInvalidType(t *testing.T) {
	var buf bytes.Buffer
	writer := &bits.WriteBuffer{Writer: &buf}

	err := writer.WriteBits("invalid", 3)
	if err == nil {
		t.Error("Write with invalid type did not return an error")
	}
}
