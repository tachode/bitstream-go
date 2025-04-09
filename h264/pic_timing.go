package h264

import (
	"github.com/tachode/bitstream-go/bits"
)

func init() { RegisterSeiPayloadType(SeiTypePicTiming, &PicTiming{}) }

type PicTiming struct {
	CpbRemovalDelay    uint64 `descriptor:"u(v)" json:"cpb_removal_delay"`
	DpbOutputDelay     uint64 `descriptor:"u(v)" json:"dpb_output_delay"`
	PicStruct          uint8  `descriptor:"u(4)" json:"pic_struct"`
	ClockTimestampFlag []bool `descriptor:"u(1)" json:"clock_timestamp_flag"`
	CtType             uint8  `descriptor:"u(2)" json:"ct_type"`
	NuitFieldBasedFlag bool   `descriptor:"u(1)" json:"nuit_field_based_flag"`
	CountingType       uint8  `descriptor:"u(5)" json:"counting_type"`
	FullTimestampFlag  bool   `descriptor:"u(1)" json:"full_timestamp_flag"`
	DiscontinuityFlag  bool   `descriptor:"u(1)" json:"discontinuity_flag"`
	CntDroppedFlag     bool   `descriptor:"u(1)" json:"cnt_dropped_flag"`
	NFrames            uint8  `descriptor:"u(8)" json:"n_frames"`
	SecondsValue       uint8  `descriptor:"u(6)" json:"seconds_value"`
	MinutesValue       uint8  `descriptor:"u(6)" json:"minutes_value"`
	HoursValue         uint8  `descriptor:"u(5)" json:"hours_value"`
	SecondsFlag        bool   `descriptor:"u(1)" json:"seconds_flag"`
	MinutesFlag        bool   `descriptor:"u(1)" json:"minutes_flag"`
	HoursFlag          bool   `descriptor:"u(1)" json:"hours_flag"`
	TimeOffset         int64  `descriptor:"i(v)" json:"time_offset"`
}

func (e *PicTiming) Read(d bits.Decoder, payloadSize int) error {
	nal_hrd_parameters_present_flag, _ := d.Value("NalHrdParametersPresentFlag").(bool)
	vcl_hrd_parameters_present_flag, _ := d.Value("VclHrdParametersPresentFlag").(bool)
	CpbDpbDelaysPresentFlag := nal_hrd_parameters_present_flag || vcl_hrd_parameters_present_flag
	if CpbDpbDelaysPresentFlag {
		initial_cpb_removal_delay_length_minus1, _ := d.Value("InitialCpbRemovalDelayLengthMinus1").(uint64)
		dpb_output_delay_length_minus1, _ := d.Value("DpbOutputDelayLengthMinus1").(uint64)
		d.SetValueLength("CpbRemovalDelay", int(initial_cpb_removal_delay_length_minus1+1))
		d.SetValueLength("DpbOutputDelay", int(dpb_output_delay_length_minus1+1))
		d.Decode(e, "CpbRemovalDelay")
		d.Decode(e, "DpbOutputDelay")
	}
	pic_struct_present_flag, _ := d.Value("PicStructPresentFlag").(bool)
	if pic_struct_present_flag {
		d.Decode(e, "PicStruct")

		// H.264 table D-1
		var NumClockTS int
		switch e.PicStruct {
		case 0, 1, 2:
			NumClockTS = 1
		case 3, 4, 7:
			NumClockTS = 2
		case 5, 6, 8:
			NumClockTS = 3
		}

		d.SetValue("NumClockTS", NumClockTS)
		if d.Error() != nil {
			return d.Error()
		}
		for i := 0; i < NumClockTS; i++ {
			d.DecodeIndex(e, "ClockTimestampFlag", i)
			if e.ClockTimestampFlag[i] {
				d.Decode(e, "CtType")
				d.Decode(e, "NuitFieldBasedFlag")
				d.Decode(e, "CountingType")
				d.Decode(e, "FullTimestampFlag")
				d.Decode(e, "DiscontinuityFlag")
				d.Decode(e, "CntDroppedFlag")
				d.Decode(e, "NFrames")
				if e.FullTimestampFlag {
					d.Decode(e, "SecondsValue")
					d.Decode(e, "MinutesValue")
					d.Decode(e, "HoursValue")
				} else {
					d.Decode(e, "SecondsFlag")
					if e.SecondsFlag {
						d.Decode(e, "SecondsValue")
						d.Decode(e, "MinutesFlag")

						if e.MinutesFlag {
							d.Decode(e, "MinutesValue")
							d.Decode(e, "HoursFlag")
							if e.HoursFlag {
								d.Decode(e, "HoursValue")
							}
						}
					}
				}
				time_offset_length, _ := d.Value("TimeOffsetLength").(uint64)
				if time_offset_length > 0 {
					d.SetValueLength("TimeOffset", int(time_offset_length))
					d.Decode(e, "TimeOffset")
				}
			}
		}
	}
	return d.Error()
}
