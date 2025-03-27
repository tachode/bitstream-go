package h264

import "github.com/tachode/bitstream-go/bits"

type NalUnitHeader3davcExtension struct {
	ViewIdx       uint8 `descriptor:"u(8)" json:"view_idx"`
	DepthFlag     bool  `descriptor:"u(1)" json:"depth_flag"`
	NonIdrFlag    bool  `descriptor:"u(1)" json:"non_idr_flag"`
	TemporalId    uint8 `descriptor:"u(3)" json:"temporal_id"`
	AnchorPicFlag bool  `descriptor:"u(1)" json:"anchor_pic_flag"`
	InterViewFlag bool  `descriptor:"u(1)" json:"inter_view_flag"`
}

func (e *NalUnitHeader3davcExtension) Read(d bits.Decoder) error {
	d.DecodeRange(e, "ViewIdx", "InterViewFlag")
	return d.Error()
}
