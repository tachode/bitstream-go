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

	result, err := ituReader.UE()
	assert.NoError(t, err)
	assert.Equal(t, uint64(3), result)
}

func TestItuReader_SE_Positive(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected int64
	}{
		{"SE(0)", 0b1_1111111, 0},
		{"SE(1)", 0b010_11111, 1},
		{"SE(2)", 0b00100_111, 2},
		{"SE(3)", 0b00110_111, 3},
		{"SE(4)", 0b0001000_1, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &bits.ReadBuffer{Reader: bytes.NewBuffer([]byte{tt.input})}
			ituReader := &bits.ItuReader{Reader: reader}

			result, err := ituReader.SE()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestItuReader_SE_Negative(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected int64
	}{
		{"SE(-1)", 0b011_11111, -1},
		{"SE(-2)", 0b00101_111, -2},
		{"SE(-3)", 0b00111_111, -3},
		{"SE(-4)", 0b0001001_1, -4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &bits.ReadBuffer{Reader: bytes.NewBuffer([]byte{tt.input})}
			ituReader := &bits.ItuReader{Reader: reader}

			result, err := ituReader.SE()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestItuReader_EOF(t *testing.T) {
	reader := &bits.ReadBuffer{Reader: bytes.NewBuffer([]byte{})}
	ituReader := &bits.ItuReader{Reader: reader}

	_, err := ituReader.UE()
	assert.Error(t, err)
	assert.ErrorIs(t, io.EOF, err)
}
