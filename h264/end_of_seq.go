package h264

import "github.com/tachode/bitstream-go/bits"

func init() { RegisterNalPayloadType(NalUnitTypeEndOfSequence, &EndOfSeq{}) }

type EndOfSeq struct {
}

func (e *EndOfSeq) Read(d bits.Decoder) error {
	return d.Error()
}
