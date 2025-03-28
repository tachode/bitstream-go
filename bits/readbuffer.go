package bits

import "io"

type Reader interface {
	io.Reader
	ReadBits(bits int) (uint64, error)
	Align() error
}

type ReadBuffer struct {
	Reader io.Reader
	buffer uint64
	bits   int
}

func (b *ReadBuffer) Read(buf []byte) (n int, err error) {
	b.Align()
	return b.Reader.Read(buf)
}

// ReadBits reads the specified number of bits from the buffer and returns the value as a uint64.
func (b *ReadBuffer) ReadBits(bits int) (uint64, error) {
	if bits < 0 || bits > 64 {
		return 0, io.ErrShortBuffer
	}
	if bits == 0 {
		return 0, nil
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

func (b *ReadBuffer) Align() error {
	b.bits = 0
	b.buffer = 0
	return nil
}
