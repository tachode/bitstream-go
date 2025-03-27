package h264

import "github.com/tachode/bitstream-go/bits"

type NalUnitHeaderMvcExtension struct {
	NonIdrFlag     bool   `descriptor:"u(1)" json:"non_idr_flag"`
	PriorityId     uint8  `descriptor:"u(6)" json:"priority_id"`
	ViewId         uint16 `descriptor:"u(10)" json:"view_id"`
	TemporalId     uint8  `descriptor:"u(3)" json:"temporal_id"`
	AnchorPicFlag  bool   `descriptor:"u(1)" json:"anchor_pic_flag"`
	InterViewFlag  bool   `descriptor:"u(1)" json:"inter_view_flag"`
	ReservedOneBit bool   `descriptor:"u(1)" json:"reserved_one_bit"`
}

func (e *NalUnitHeaderMvcExtension) Read(d bits.Decoder) error {
	d.DecodeRange(e, "NonIdrFlag", "ReservedOneBit")
	return d.Error()
}
