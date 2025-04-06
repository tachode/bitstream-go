package h264

import (
	"github.com/tachode/bitstream-go/bits"
)

func init() { RegisterSeiPayloadType(SeiTypeBufferingPeriod, &BufferingPeriod{}) }

type BufferingPeriod struct {
	SeqParameterSetId               uint64   `descriptor:"ue(v)" json:"seq_parameter_set_id"`
	NalInitialCpbRemovalDelay       []uint64 `descriptor:"u(v)" json:"nal_initial_cpb_removal_delay"`
	NalInitialCpbRemovalDelayOffset []uint64 `descriptor:"u(v)" json:"nal_initial_cpb_removal_delay_offset"`
	VclInitialCpbRemovalDelay       []uint64 `descriptor:"u(v)" json:"vcl_initial_cpb_removal_delay"`
	VclInitialCpbRemovalDelayOffset []uint64 `descriptor:"u(v)" json:"vcl_initial_cpb_removal_delay_offset"`
}

func (e *BufferingPeriod) Read(d bits.Decoder, payloadSize int) error {
	// From VUI
	NalHrdBpPresentFlag, _ := d.Value("NalHrdParametersPresentFlag").(bool)
	VclHrdBpPresentFlag, _ := d.Value("VclHrdParametersPresentFlag").(bool)

	// From HRD parameters, inside VUI, which is inside SPS
	cpb_cnt_minus1, _ := d.Value("CpbCntMinus1").(uint64)
	initial_cpb_removal_delay_length_minus1, _ := d.Value("InitialCpbRemovalDelayLengthMinus1").(uint64)
	d.SetValueLength("NalInitialCpbRemovalDelay", int(initial_cpb_removal_delay_length_minus1+1))
	d.SetValueLength("NalInitialCpbRemovalDelayOffset", int(initial_cpb_removal_delay_length_minus1+1))
	d.SetValueLength("VclInitialCpbRemovalDelay", int(initial_cpb_removal_delay_length_minus1+1))
	d.SetValueLength("VclInitialCpbRemovalDelayOffset", int(initial_cpb_removal_delay_length_minus1+1))

	d.Decode(e, "SeqParameterSetId")
	if NalHrdBpPresentFlag {
		for SchedSelIdx := 0; SchedSelIdx <= int(cpb_cnt_minus1); SchedSelIdx++ {
			d.DecodeIndex(e, "NalInitialCpbRemovalDelay", SchedSelIdx)
			d.DecodeIndex(e, "NalInitialCpbRemovalDelayOffset", SchedSelIdx)
		}
	}
	if VclHrdBpPresentFlag {
		for SchedSelIdx := 0; SchedSelIdx <= int(cpb_cnt_minus1); SchedSelIdx++ {
			d.DecodeIndex(e, "VclInitialCpbRemovalDelay", SchedSelIdx)
			d.DecodeIndex(e, "VclInitialCpbRemovalDelayOffset", SchedSelIdx)
		}
	}
	return d.Error()
}
