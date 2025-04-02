package h264

import "github.com/tachode/bitstream-go/bits"

func init() { RegisterNalPayload(NalUnitTypeEndOfStream, &EndOfStream{}) }

type EndOfStream struct {
}

func (e *EndOfStream) Read(d bits.Decoder) error {
	return d.Error()
}
