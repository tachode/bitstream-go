package bits

import "io"

type Decoder interface {
	io.Reader
	Decode(v any, field string) error
	DecodeRange(v any, start string, end string) error
	DecodeIndex(v any, field string, i int) error
	Value(name string) any
	Error() error
}
