package h264

import "github.com/tachode/bitstream-go/bits"

func init() { RegisterNalPayload(NalUnitTypeFiller, &FillerDataRbsp{}) }

type FillerDataRbsp struct {
	FfByte uint8 `descriptor:"f(8)=255" json:"ff_byte"`
}

func (e *FillerDataRbsp) Read(d bits.Decoder) error {
	for d.NextBits(8) == 0xFF {
		d.Decode(e, "FfByte")
	}
	return d.Error()
}
