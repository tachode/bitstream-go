package h264

import "github.com/tachode/bitstream-go/bits"

type MvcVuiParametersExtension struct {
	VuiMvcNumOpsMinus1                uint64         `descriptor:"ue(v)" json:"vui_mvc_num_ops_minus1"`
	VuiMvcTemporalId                  []uint8        `descriptor:"u(3)" json:"vui_mvc_temporal_id"`
	VuiMvcNumTargetOutputViewsMinus1  []uint64       `descriptor:"ue(v)" json:"vui_mvc_num_target_output_views_minus1"`
	VuiMvcViewId                      [][]uint64     `descriptor:"ue(v)" json:"vui_mvc_view_id"`
	VuiMvcTimingInfoPresentFlag       []bool         `descriptor:"u(1)" json:"vui_mvc_timing_info_present_flag"`
	VuiMvcNumUnitsInTick              []uint32       `descriptor:"u(32)" json:"vui_mvc_num_units_in_tick"`
	VuiMvcTimeScale                   []uint32       `descriptor:"u(32)" json:"vui_mvc_time_scale"`
	VuiMvcFixedFrameRateFlag          []bool         `descriptor:"u(1)" json:"vui_mvc_fixed_frame_rate_flag"`
	VuiMvcNalHrdParametersPresentFlag []bool         `descriptor:"u(1)" json:"vui_mvc_nal_hrd_parameters_present_flag"`
	NalHrdParameters                  *HrdParameters `json:"nal_hrd_parameters,omitempty"`
	VuiMvcVclHrdParametersPresentFlag []bool         `descriptor:"u(1)" json:"vui_mvc_vcl_hrd_parameters_present_flag"`
	VclHrdParameters                  *HrdParameters `json:"vcl_hrd_parameters,omitempty"`
	VuiMvcLowDelayHrdFlag             []bool         `descriptor:"u(1)" json:"vui_mvc_low_delay_hrd_flag"`
	VuiMvcPicStructPresentFlag        []bool         `descriptor:"u(1)" json:"vui_mvc_pic_struct_present_flag"`
}

func (e *MvcVuiParametersExtension) Read(d bits.Decoder) error {
	d.Decode(e, "VuiMvcNumOpsMinus1")
	for i := 0; i <= int(e.VuiMvcNumOpsMinus1); i++ {
		d.DecodeIndex(e, "VuiMvcTemporalId", i)
		d.DecodeIndex(e, "VuiMvcNumTargetOutputViewsMinus1", i)
		for j := 0; j <= int(e.VuiMvcNumTargetOutputViewsMinus1[i]); j++ {
			d.DecodeIndex(e, "VuiMvcViewId", i, j)
		}
		d.DecodeIndex(e, "VuiMvcTimingInfoPresentFlag", i)
		if e.VuiMvcTimingInfoPresentFlag[i] {
			d.DecodeIndex(e, "VuiMvcNumUnitsInTick", i)
			d.DecodeIndex(e, "VuiMvcTimeScale", i)
			d.DecodeIndex(e, "VuiMvcFixedFrameRateFlag", i)
		}
		d.DecodeIndex(e, "VuiMvcNalHrdParametersPresentFlag", i)
		if e.VuiMvcNalHrdParametersPresentFlag[i] {
			e.NalHrdParameters = &HrdParameters{}
			e.NalHrdParameters.Read(d)
		}
		d.DecodeIndex(e, "VuiMvcVclHrdParametersPresentFlag", i)
		if e.VuiMvcVclHrdParametersPresentFlag[i] {
			e.VclHrdParameters = &HrdParameters{}
			e.VclHrdParameters.Read(d)
		}
		if e.VuiMvcNalHrdParametersPresentFlag[i] || e.VuiMvcVclHrdParametersPresentFlag[i] {
			d.DecodeIndex(e, "VuiMvcLowDelayHrdFlag", i)
		}
		d.DecodeIndex(e, "VuiMvcPicStructPresentFlag", i)
	}
	return d.Error()
}
