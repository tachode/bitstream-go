package h264

import "github.com/tachode/bitstream-go/bits"

type MvcdVuiParametersExtension struct {
	VuiMvcdNumOpsMinus1                uint64         `descriptor:"ue(v)" json:"vui_mvcd_num_ops_minus1"`
	VuiMvcdTemporalId                  []uint8        `descriptor:"u(3)" json:"vui_mvcd_temporal_id"`
	VuiMvcdNumTargetOutputViewsMinus1  []uint64       `descriptor:"ue(v)" json:"vui_mvcd_num_target_output_views_minus1"`
	VuiMvcdViewId                      [][]uint64     `descriptor:"ue(v)" json:"vui_mvcd_view_id"`
	VuiMvcdDepthFlag                   [][]bool       `descriptor:"u(1)" json:"vui_mvcd_depth_flag"`
	VuiMvcdTextureFlag                 [][]bool       `descriptor:"u(1)" json:"vui_mvcd_texture_flag"`
	VuiMvcdTimingInfoPresentFlag       []bool         `descriptor:"u(1)" json:"vui_mvcd_timing_info_present_flag"`
	VuiMvcdNumUnitsInTick              []uint32       `descriptor:"u(32)" json:"vui_mvcd_num_units_in_tick"`
	VuiMvcdTimeScale                   []uint32       `descriptor:"u(32)" json:"vui_mvcd_time_scale"`
	VuiMvcdFixedFrameRateFlag          []bool         `descriptor:"u(1)" json:"vui_mvcd_fixed_frame_rate_flag"`
	VuiMvcdNalHrdParametersPresentFlag []bool         `descriptor:"u(1)" json:"vui_mvcd_nal_hrd_parameters_present_flag"`
	NalHrdParameters                   *HrdParameters `json:"nal_hrd_parameters,omitempty"`
	VuiMvcdVclHrdParametersPresentFlag []bool         `descriptor:"u(1)" json:"vui_mvcd_vcl_hrd_parameters_present_flag"`
	VclHrdParameters                   *HrdParameters `json:"vcl_hrd_parameters,omitempty"`
	VuiMvcdLowDelayHrdFlag             []bool         `descriptor:"u(1)" json:"vui_mvcd_low_delay_hrd_flag"`
	VuiMvcdPicStructPresentFlag        []bool         `descriptor:"u(1)" json:"vui_mvcd_pic_struct_present_flag"`
}

func (e *MvcdVuiParametersExtension) Read(d bits.Decoder) error {
	d.Decode(e, "VuiMvcdNumOpsMinus1")
	for i := 0; i <= int(e.VuiMvcdNumOpsMinus1); i++ {
		d.DecodeIndex(e, "VuiMvcdTemporalId", i)
		d.DecodeIndex(e, "VuiMvcdNumTargetOutputViewsMinus1", i)
		for j := 0; j <= int(e.VuiMvcdNumTargetOutputViewsMinus1[i]); j++ {
			d.DecodeIndex(e, "VuiMvcdViewId", i, j)
			d.DecodeIndex(e, "VuiMvcdDepthFlag", i, j)
			d.DecodeIndex(e, "VuiMvcdTextureFlag", i, j)
		}
		d.DecodeIndex(e, "VuiMvcdTimingInfoPresentFlag", i)
		if e.VuiMvcdTimingInfoPresentFlag[i] {
			d.DecodeIndex(e, "VuiMvcdNumUnitsInTick", i)
			d.DecodeIndex(e, "VuiMvcdTimeScale", i)
			d.DecodeIndex(e, "VuiMvcdFixedFrameRateFlag", i)
		}
		d.DecodeIndex(e, "VuiMvcdNalHrdParametersPresentFlag", i)
		if e.VuiMvcdNalHrdParametersPresentFlag[i] {
			e.NalHrdParameters = &HrdParameters{}
			e.NalHrdParameters.Read(d)
		}
		d.DecodeIndex(e, "VuiMvcdVclHrdParametersPresentFlag", i)
		if e.VuiMvcdVclHrdParametersPresentFlag[i] {
			e.VclHrdParameters = &HrdParameters{}
			e.VclHrdParameters.Read(d)
		}
		if e.VuiMvcdNalHrdParametersPresentFlag[i] || e.VuiMvcdVclHrdParametersPresentFlag[i] {
			d.DecodeIndex(e, "VuiMvcdLowDelayHrdFlag", i)
		}
		d.DecodeIndex(e, "VuiMvcdPicStructPresentFlag", i)
	}
	return d.Error()
}
