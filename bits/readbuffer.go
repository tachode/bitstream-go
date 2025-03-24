package bits

import "io"

type ReadBuffer struct {
	Reader io.Reader
	buffer uint64
	bits   int
}

// Read reads the specified number of bits from the buffer and returns the value as a uint64.
func (b *ReadBuffer) Read(bits int) (uint64, error) {
	if bits <= 0 || bits > 64 {
		return 0, io.ErrShortBuffer
	}

	// Ensure the buffer has enough bits
	for b.bits < bits {
		var byteBuf [1]byte
		_, err := b.Reader.Read(byteBuf[:])
		if err != nil {
			return 0, err
		}
		b.buffer = (b.buffer << 8) | uint64(byteBuf[0])
		b.bits += 8
	}

	// Extract the requested bits
	shift := b.bits - bits
	value := (b.buffer >> shift) & ((1 << bits) - 1)
	b.buffer &= (1 << shift) - 1
	b.bits -= bits

	return value, nil
}

// Align discards bits until the bitstream is aligned on a byte boundary.
func (b *ReadBuffer) Align() error {
	remainingBits := b.bits % 8
	if remainingBits > 0 {
		_, err := b.Read(remainingBits)
		if err != nil {
			return err
		}
	}
	return nil
}
