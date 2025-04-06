package h264

import "github.com/tachode/bitstream-go/bits"

func init() { RegisterNalPayloadType(NalUnitTypeAUD, &AccessUnitDelimiter{}) }

type AccessUnitDelimiter struct {
	PrimaryPicType uint8 `descriptor:"u(3)" json:"primary_pic_type"`
}

func (e *AccessUnitDelimiter) Read(d bits.Decoder) error {
	d.Decode(e, "PrimaryPicType")
	return d.Error()
}
