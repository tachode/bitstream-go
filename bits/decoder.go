package bits

type Decoder interface {
	Decode(v any, field string)
	DecodeRange(v any, start string, end string)
	DecodeIndex(v any, field string, i int)
	Error() error
}
