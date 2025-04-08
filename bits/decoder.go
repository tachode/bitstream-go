package bits

import "io"

type Decoder interface {
	io.Reader
	Decode(v any, field string) error
	DecodeRange(v any, start string, end string) error
	DecodeIndex(v any, field string, i int, j ...int) error
	Value(name string) any
	SetValue(name string, value any)
	SetValueLength(name string, length int)
	Reset(in any) error
	Error() error

	// Functions from H.264, §7.2
	ByteAligned() bool
	MoreRbspData() bool
	NextBits(bits int) uint64
}
