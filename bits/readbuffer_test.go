package bits_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/tachode/bitstream-go/bits"
)

func TestReadBuffer_Read(t *testing.T) {
	data := []byte{0b10101010, 0b11001100, 0b11110000}
	reader := bytes.NewReader(data)
	readBuffer := &bits.ReadBuffer{Reader: reader}

	tests := []struct {
		bitsToRead int
		expected   uint64
		expectErr  bool
	}{
		{4, 0b1010, false},
		{4, 0b1010, false},
		{8, 0b11001100, false},
		{4, 0b1111, false},
		{4, 0b0000, false},
		{1, 0, true}, // No more bits to read
	}

	for _, test := range tests {
		result, err := readBuffer.Read(test.bitsToRead)
		if test.expectErr {
			if err == nil {
				t.Errorf("expected error but got none")
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != test.expected {
				t.Errorf("expected %b, got %b", test.expected, result)
			}
		}
	}
}

func TestReadBuffer_Align(t *testing.T) {
	data := []byte{0b10101010, 0b11001100}
	reader := bytes.NewReader(data)
	readBuffer := &bits.ReadBuffer{Reader: reader}

	// Read 5 bits
	_, err := readBuffer.Read(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Align to the next byte boundary
	err = readBuffer.Align()
	if err != nil {
		t.Fatalf("unexpected error during align: %v", err)
	}

	// Read the next byte
	result, err := readBuffer.Read(8)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 0b11001100 {
		t.Errorf("expected %b, got %b", 0b11001100, result)
	}
}

func TestReadBuffer_ReadInvalidBits(t *testing.T) {
	reader := bytes.NewReader([]byte{0b10101010})
	readBuffer := &bits.ReadBuffer{Reader: reader}

	_, err := readBuffer.Read(0)
	if err != io.ErrShortBuffer {
		t.Errorf("expected io.ErrShortBuffer, got %v", err)
	}

	_, err = readBuffer.Read(65)
	if err != io.ErrShortBuffer {
		t.Errorf("expected io.ErrShortBuffer, got %v", err)
	}
}
