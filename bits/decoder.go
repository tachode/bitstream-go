package bits

type Decoder interface {
	Decode(v any, field string) error
	DecodeRange(v any, start string, end string) error
	DecodeIndex(v any, field string, i int) error
	Error() error
}
