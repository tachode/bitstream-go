package bits_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/bitstream-go/bits"
)

///////////////////////////////////////////////////////////////////////////

func TestItuReader_UE(t *testing.T) {
	reader := &bits.ReadBuffer{Reader: bytes.NewBuffer([]byte{0b00100_111})}
	ituReader := &bits.ItuReader{Reader: reader}

	result, n, err := ituReader.UE()
	assert.NoError(t, err)
	assert.Equal(t, uint64(3), result)
	assert.Equal(t, 5, n)
}

func TestItuReader_SE_Positive(t *testing.T) {
	tests := []struct {
		name         string
		input        byte
		expected     int64
		expectedBits int
	}{
		{"SE(0)", 0b1_1111111, 0, 1},
		{"SE(1)", 0b010_11111, 1, 3},
		{"SE(2)", 0b00100_111, 2, 5},
		{"SE(3)", 0b00110_111, 3, 5},
		{"SE(4)", 0b0001000_1, 4, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &bits.ReadBuffer{Reader: bytes.NewBuffer([]byte{tt.input})}
			ituReader := &bits.ItuReader{Reader: reader}

			result, n, err := ituReader.SE()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectedBits, n)
		})
	}
}

func TestItuReader_SE_Negative(t *testing.T) {
	tests := []struct {
		name         string
		input        byte
		expected     int64
		expectedBits int
	}{
		{"SE(-1)", 0b011_11111, -1, 3},
		{"SE(-2)", 0b00101_111, -2, 5},
		{"SE(-3)", 0b00111_111, -3, 5},
		{"SE(-4)", 0b0001001_1, -4, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &bits.ReadBuffer{Reader: bytes.NewBuffer([]byte{tt.input})}
			ituReader := &bits.ItuReader{Reader: reader}

			result, n, err := ituReader.SE()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectedBits, n)
		})
	}
}

func TestItuReader_EOF(t *testing.T) {
	reader := &bits.ReadBuffer{Reader: bytes.NewBuffer([]byte{})}
	ituReader := &bits.ItuReader{Reader: reader}

	_, _, err := ituReader.UE()
	assert.Error(t, err)
	assert.ErrorIs(t, io.EOF, err)
}
func TestItuReader_I(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		bits     int
		expected int64
	}{
		{"I(0)  - 1 bit ", 0b0_1111111, 1, 0},
		{"I(1)  - 2 bits", 0b01_111111, 2, 1},
		{"I(-1) - 1 bit ", 0b1_1111111, 1, -1},
		{"I(-1) - 2 bits", 0b11_111111, 2, -1},
		{"I(2)  - 3 bits", 0b010_11111, 3, 2},
		{"I(-2) - 3 bits", 0b110_11111, 3, -2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &bits.ReadBuffer{Reader: bytes.NewBuffer([]byte{tt.input})}
			ituReader := &bits.ItuReader{Reader: reader}

			result, n, err := ituReader.I(tt.bits)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.bits, n)
		})
	}
}
