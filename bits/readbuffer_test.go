package bits_test

import (
	"bufio"
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
		result, err := readBuffer.ReadBits(test.bitsToRead)
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
	_, err := readBuffer.ReadBits(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Align to the next byte boundary
	err = readBuffer.Align()
	if err != nil {
		t.Fatalf("unexpected error during align: %v", err)
	}

	// Read the next byte
	result, err := readBuffer.ReadBits(8)
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

	_, err := readBuffer.ReadBits(-1)
	if err != io.ErrShortBuffer {
		t.Errorf("expected io.ErrShortBuffer, got %v", err)
	}

	_, err = readBuffer.ReadBits(65)
	if err != io.ErrShortBuffer {
		t.Errorf("expected io.ErrShortBuffer, got %v", err)
	}
}

func TestReadBuffer_MoreRbspData_1(t *testing.T) {
	data := []byte{0b10101010, 0b00000001, 0b10000000} // RBSP data with stop bit
	reader := bufio.NewReader(bytes.NewReader(data))
	readBuffer := &bits.ReadBuffer{Reader: reader}

	// Read some bits to simulate partial consumption
	_, err := readBuffer.ReadBits(8)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check if there's more RBSP data
	if !readBuffer.MoreRbspData() {
		t.Errorf("expected more RBSP data, but got false")
	}

	// Consume all remaining bits except one
	_, err = readBuffer.ReadBits(7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check if there's more RBSP data
	if !readBuffer.MoreRbspData() {
		t.Errorf("expected more RBSP data, but got false")
	}

	// Consume final bit
	_, err = readBuffer.ReadBits(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check again after consuming all data
	if readBuffer.MoreRbspData() {
		t.Errorf("expected no more RBSP data, but got true")
	}
}

func TestReadBuffer_MoreRbspData_2(t *testing.T) {
	data := []byte{0b10101010, 0b00000000, 0b10000000} // RBSP data with stop bit
	reader := bufio.NewReader(bytes.NewReader(data))
	readBuffer := &bits.ReadBuffer{Reader: reader}

	// Read some bits to simulate partial consumption
	_, err := readBuffer.ReadBits(8)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check if there's more RBSP data
	if !readBuffer.MoreRbspData() {
		t.Errorf("expected more RBSP data, but got false")
	}

	// Consume all remaining bits except one
	_, err = readBuffer.ReadBits(7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check if there's more RBSP data
	if !readBuffer.MoreRbspData() {
		t.Errorf("expected more RBSP data, but got false")
	}

	// Consume final bit
	_, err = readBuffer.ReadBits(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check again after consuming all data
	if readBuffer.MoreRbspData() {
		t.Errorf("expected no more RBSP data, but got true")
	}
}
