package bits

import (
	"fmt"
	"io"
	"reflect"
)

type Writer interface {
	Write(v any, bits int) error
	Flush() error
}

type WriteBuffer struct {
	Writer io.Writer
	buffer uint64
	bits   int
}

func (b *WriteBuffer) Write(v any, bits int) error {
	val := reflect.ValueOf(v)
	var value uint64
	switch {
	case val.CanInt():
		value = uint64(val.Int())
	case val.CanUint():
		value = val.Uint()
	default:
		return fmt.Errorf("cannot convert %T into a uint64", v)
	}

	b.buffer = (b.buffer << bits) | (value & ((1 << bits) - 1))
	b.bits += bits

	for b.bits >= 8 {
		byteToWrite := byte(b.buffer >> (b.bits - 8))
		if _, err := b.Writer.Write([]byte{byteToWrite}); err != nil {
			return err
		}
		b.bits -= 8
		b.buffer &= (1 << b.bits) - 1
	}

	return nil
}

func (b *WriteBuffer) Flush() error {
	if b.bits > 0 {
		byteToWrite := byte(b.buffer << (8 - b.bits))
		if _, err := b.Writer.Write([]byte{byteToWrite}); err != nil {
			return err
		}
		b.buffer = 0
		b.bits = 0
	}
	return nil
}
