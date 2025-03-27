package bits_test

import (
	"bytes"
	"testing"

	"github.com/tachode/bitstream-go/bits"
)

type mockWriteBuffer struct {
	buffer bytes.Buffer
}

func (m *mockWriteBuffer) Write(value any, bits int) error {
	v := value.(uint64)
	for i := bits - 1; i >= 0; i-- {
		bit := (v >> i) & 1
		if bit == 1 {
			m.buffer.WriteByte('1')
		} else {
			m.buffer.WriteByte('0')
		}
	}
	return nil
}

func (m *mockWriteBuffer) Flush() error { return nil }

///////////////////////////////////////////////////////////////////////////

func TestItuWriter_UE(t *testing.T) {
	mockBuffer := &mockWriteBuffer{}
	writer := &bits.ItuWriter{Writer: mockBuffer}

	tests := []struct {
		input    uint64
		expected string
	}{
		{0, "1"},
		{1, "010"},
		{2, "011"},
		{3, "00100"},
		{4, "00101"},
	}

	for _, test := range tests {
		mockBuffer.buffer.Reset()
		err := writer.UE(test.input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if mockBuffer.buffer.String() != test.expected {
			t.Errorf("UE(%d): expected %s, got %s", test.input, test.expected, mockBuffer.buffer.String())
		}
	}
}

func TestItuWriter_SE(t *testing.T) {
	mockBuffer := &mockWriteBuffer{}
	writer := &bits.ItuWriter{Writer: mockBuffer}

	tests := []struct {
		input    int64
		expected string
	}{
		{0, "1"},
		{1, "010"},
		{-1, "011"},
		{2, "00100"},
		{-2, "00101"},
		{3, "00110"},
		{-3, "00111"},
		{4, "0001000"},
		{-4, "0001001"},
	}

	for _, test := range tests {
		mockBuffer.buffer.Reset()
		err := writer.SE(test.input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if mockBuffer.buffer.String() != test.expected {
			t.Errorf("SE(%d): expected %s, got %s", test.input, test.expected, mockBuffer.buffer.String())
		}
	}
}
