package h264

import "github.com/tachode/bitstream-go/bits"

type HrdParameters struct {
	CpbCntMinus1                       uint64   `descriptor:"ue(v)" json:"cpb_cnt_minus1"`
	BitRateScale                       uint8    `descriptor:"u(4)" json:"bit_rate_scale"`
	CpbSizeScale                       uint8    `descriptor:"u(4)" json:"cpb_size_scale"`
	BitRateValueMinus1                 []uint64 `descriptor:"ue(v)" json:"bit_rate_value_minus1"`
	CpbSizeValueMinus1                 []uint64 `descriptor:"ue(v)" json:"cpb_size_value_minus1"`
	CbrFlag                            []bool   `descriptor:"u(1)" json:"cbr_flag"`
	InitialCpbRemovalDelayLengthMinus1 uint8    `descriptor:"u(5)" json:"initial_cpb_removal_delay_length_minus1"`
	CpbRemovalDelayLengthMinus1        uint8    `descriptor:"u(5)" json:"cpb_removal_delay_length_minus1"`
	DpbOutputDelayLengthMinus1         uint8    `descriptor:"u(5)" json:"dpb_output_delay_length_minus1"`
	TimeOffsetLength                   uint8    `descriptor:"u(5)" json:"time_offset_length"`
}

func (e *HrdParameters) Read(d bits.Decoder) error {
	d.Decode(e, "CpbCntMinus1")
	d.Decode(e, "BitRateScale")
	d.Decode(e, "CpbSizeScale")
	for SchedSelIdx := 0; SchedSelIdx <= int(e.CpbCntMinus1); SchedSelIdx++ {
		d.DecodeIndex(e, "BitRateValueMinus1", SchedSelIdx)
		d.DecodeIndex(e, "CpbSizeValueMinus1", SchedSelIdx)
		d.DecodeIndex(e, "CbrFlag", SchedSelIdx)
	}
	d.Decode(e, "InitialCpbRemovalDelayLengthMinus1")
	d.Decode(e, "CpbRemovalDelayLengthMinus1")
	d.Decode(e, "DpbOutputDelayLengthMinus1")
	d.Decode(e, "TimeOffsetLength")
	return d.Error()
}
