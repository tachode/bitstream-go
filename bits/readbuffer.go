package bits

import (
	"bufio"
	"errors"
	"io"
)

type Reader interface {
	io.Reader
	ReadBits(bits int) (uint64, error)
	NextBits(bits int) (uint64, error)
	ByteAligned() bool
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

func (b ReadBuffer) ByteAligned() bool {
	return b.bits == 0
}

// This function requires that the ReadBuffer was constructed with a
// bufio.Reader, and also imposes some byte alignment constraints that
// are true for all uses of `next_bits()` in H.264
func (b *ReadBuffer) NextBits(bits int) (uint64, error) {
	buf, ok := b.Reader.(*bufio.Reader)
	if !ok {
		return 0, errors.New("cannot peek: ReadBuffer's Reader is not buffered")
	}
	if bits%8 != 0 {
		return 0, errors.New("can only peek whole bytes")
	}
	if !b.ByteAligned() {
		return 0, errors.New("can only peek on byte boundaries")
	}
	v, err := buf.Peek(bits / 8)
	if err != nil {
		return 0, err
	}
	val := uint64(0)
	for len(v) > 0 {
		val <<= 8
		val += uint64(v[0])
		v = v[1:]
	}
	return val, nil
}

// This function assumes that trailing RBSP zero bytes have been
// trimmed and that the underlying Reader is an iobuf.Reader.
func (b *ReadBuffer) MoreRbspData() bool {
	buf, ok := b.Reader.(*bufio.Reader)
	if !ok {
		panic("cannot call more_rbsp_data() on a non-RBSP payload")
	}

	// If there is still a zero bit to read, then there is more data
	// before the stop bit
	if b.bits > 0 {
		nextBit := b.buffer & (1 << (b.bits - 1))
		if nextBit == 0 {
			return true
		}
	}

	// b.buffer only has values in the lowest 8 bits between reads.
	// If there's more than one bit set, then there is more RBSP
	// data: the final bit is the rbsp_stop_one_bit, and the other
	// must be RBSP data.
	oneBitFound := false
	switch b.buffer {
	case 0:
		oneBitFound = false
	// These are the only values that have exactly one bit set
	case 1, 2, 4, 8, 16, 32, 64, 128:
		oneBitFound = true
	default:
		return true
	}

	// If we can peek two bytes, then there *must* be more data left.
	// The rsbp_stop_one_bit will be in the final byte (since we trimmed
	// all zero bytes), so the byte before it must contain RBSP data.
	_, err := buf.Peek(2)
	if err == nil {
		return true
	}

	// If we can peek one byte but not two, then it is the final byte in
	// the RBSP
	finalByte, err := buf.Peek(1)
	if err != nil {
		return false
	}

	// If there's a one-bit in b.buffer, then there's more data,
	// since we know the final byte can't be zero
	if oneBitFound {
		return true
	}

	// Finally, check if the final byte is simply a stop bit with no
	// additional bits
	if finalByte[0] == 0b1000_0000 {
		return false
	}
	return true
}
