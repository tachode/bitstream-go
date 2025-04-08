package h264

import "github.com/tachode/bitstream-go/bits"

type SvcVuiParametersExtension struct {
	VuiExtNumEntriesMinus1            uint64         `descriptor:"ue(v)" json:"vui_ext_num_entries_minus1"`
	VuiExtDependencyId                []uint8        `descriptor:"u(3)" json:"vui_ext_dependency_id"`
	VuiExtQualityId                   []uint8        `descriptor:"u(4)" json:"vui_ext_quality_id"`
	VuiExtTemporalId                  []uint8        `descriptor:"u(3)" json:"vui_ext_temporal_id"`
	VuiExtTimingInfoPresentFlag       []bool         `descriptor:"u(1)" json:"vui_ext_timing_info_present_flag"`
	VuiExtNumUnitsInTick              []uint32       `descriptor:"u(32)" json:"vui_ext_num_units_in_tick"`
	VuiExtTimeScale                   []uint32       `descriptor:"u(32)" json:"vui_ext_time_scale"`
	VuiExtFixedFrameRateFlag          []bool         `descriptor:"u(1)" json:"vui_ext_fixed_frame_rate_flag"`
	VuiExtNalHrdParametersPresentFlag []bool         `descriptor:"u(1)" json:"vui_ext_nal_hrd_parameters_present_flag"`
	NalHrdParameters                  *HrdParameters `json:"nal_hrd_parameters,omitempty"`
	VuiExtVclHrdParametersPresentFlag []bool         `descriptor:"u(1)" json:"vui_ext_vcl_hrd_parameters_present_flag"`
	VclHrdParameters                  *HrdParameters `json:"vcl_hrd_parameters,omitempty"`
	VuiExtLowDelayHrdFlag             []bool         `descriptor:"u(1)" json:"vui_ext_low_delay_hrd_flag"`
	VuiExtPicStructPresentFlag        []bool         `descriptor:"u(1)" json:"vui_ext_pic_struct_present_flag"`
}

func (e *SvcVuiParametersExtension) Read(d bits.Decoder) error {
	d.Decode(e, "VuiExtNumEntriesMinus1")
	for i := 0; i <= int(e.VuiExtNumEntriesMinus1); i++ {
		d.DecodeIndex(e, "VuiExtDependencyId", i)
		d.DecodeIndex(e, "VuiExtQualityId", i)
		d.DecodeIndex(e, "VuiExtTemporalId", i)
		d.DecodeIndex(e, "VuiExtTimingInfoPresentFlag", i)
		if e.VuiExtTimingInfoPresentFlag[i] {
			d.DecodeIndex(e, "VuiExtNumUnitsInTick", i)
			d.DecodeIndex(e, "VuiExtTimeScale", i)
			d.DecodeIndex(e, "VuiExtFixedFrameRateFlag", i)
		}
		d.DecodeIndex(e, "VuiExtNalHrdParametersPresentFlag", i)
		if e.VuiExtNalHrdParametersPresentFlag[i] {
			e.NalHrdParameters = &HrdParameters{}
			e.NalHrdParameters.Read(d)
		}
		d.DecodeIndex(e, "VuiExtVclHrdParametersPresentFlag", i)
		if e.VuiExtVclHrdParametersPresentFlag[i] {
			e.VclHrdParameters = &HrdParameters{}
			e.VclHrdParameters.Read(d)
		}
		if e.VuiExtNalHrdParametersPresentFlag[i] || e.VuiExtVclHrdParametersPresentFlag[i] {
			d.DecodeIndex(e, "VuiExtLowDelayHrdFlag", i)
		}
		d.DecodeIndex(e, "VuiExtPicStructPresentFlag", i)
	}
	return d.Error()
}
