package bits

type ItuDecoder struct {
	reader *ItuReader
	err    error
}

func NewItuDecoder(r *ItuReader) *ItuDecoder {
	return &ItuDecoder{reader: r}
}

func (d *ItuDecoder) Decode(v any, field string) {
	d.DecodeRange(v, field, field)
}

func (d *ItuDecoder) DecodeRange(v any, start string, end string) {
	// TODO
}

func (d *ItuDecoder) DecodeIndex(v any, field string, i int) {
	// TODO
}

func (d *ItuDecoder) Error() error {
	return d.err
}
