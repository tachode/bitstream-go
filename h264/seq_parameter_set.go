package h264

import "github.com/tachode/bitstream-go/bits"

type SeqParameterSet struct {
	ProfileIdc                      uint8          `descriptor:"u(8)" json:"profile_idc"`
	ConstraintSet0Flag              bool           `descriptor:"u(1)" json:"constraint_set0_flag"`
	ConstraintSet1Flag              bool           `descriptor:"u(1)" json:"constraint_set1_flag"`
	ConstraintSet2Flag              bool           `descriptor:"u(1)" json:"constraint_set2_flag"`
	ConstraintSet3Flag              bool           `descriptor:"u(1)" json:"constraint_set3_flag"`
	ConstraintSet4Flag              bool           `descriptor:"u(1)" json:"constraint_set4_flag"`
	ConstraintSet5Flag              bool           `descriptor:"u(1)" json:"constraint_set5_flag"`
	ReservedZero2bits               uint8          `descriptor:"u(2)" json:"reserved_zero_2bits"`
	LevelIdc                        uint8          `descriptor:"u(8)" json:"level_idc"`
	SeqParameterSetId               uint64         `descriptor:"ue(v)" json:"seq_parameter_set_id"`
	ChromaFormatIdc                 uint64         `descriptor:"ue(v)" json:"chroma_format_idc"`
	SeparateColourPlaneFlag         bool           `descriptor:"u(1)" json:"separate_colour_plane_flag"`
	BitDepthLumaMinus8              uint64         `descriptor:"ue(v)" json:"bit_depth_luma_minus8"`
	BitDepthChromaMinus8            uint64         `descriptor:"ue(v)" json:"bit_depth_chroma_minus8"`
	QpprimeYZeroTransformBypassFlag bool           `descriptor:"u(1)" json:"qpprime_y_zero_transform_bypass_flag"`
	SeqScalingMatrixPresentFlag     bool           `descriptor:"u(1)" json:"seq_scaling_matrix_present_flag"`
	SeqScalingListPresentFlag       []bool         `descriptor:"u(1)" json:"seq_scaling_list_present_flag"`
	ScalingList                     *ScalingList   `json:"scaling_list,omitempty"`
	Log2MaxFrameNumMinus4           uint64         `descriptor:"ue(v)" json:"log2_max_frame_num_minus4"`
	PicOrderCntType                 uint64         `descriptor:"ue(v)" json:"pic_order_cnt_type"`
	Log2MaxPicOrderCntLsbMinus4     uint64         `descriptor:"ue(v)" json:"log2_max_pic_order_cnt_lsb_minus4"`
	DeltaPicOrderAlwaysZeroFlag     bool           `descriptor:"u(1)" json:"delta_pic_order_always_zero_flag"`
	OffsetForNonRefPic              int64          `descriptor:"se(v)" json:"offset_for_non_ref_pic"`
	OffsetForTopToBottomField       int64          `descriptor:"se(v)" json:"offset_for_top_to_bottom_field"`
	NumRefFramesInPicOrderCntCycle  uint64         `descriptor:"ue(v)" json:"num_ref_frames_in_pic_order_cnt_cycle"`
	OffsetForRefFrame               []int64        `descriptor:"se(v)" json:"offset_for_ref_frame"`
	MaxNumRefFrames                 uint64         `descriptor:"ue(v)" json:"max_num_ref_frames"`
	GapsInFrameNumValueAllowedFlag  bool           `descriptor:"u(1)" json:"gaps_in_frame_num_value_allowed_flag"`
	PicWidthInMbsMinus1             uint64         `descriptor:"ue(v)" json:"pic_width_in_mbs_minus1"`
	PicHeightInMapUnitsMinus1       uint64         `descriptor:"ue(v)" json:"pic_height_in_map_units_minus1"`
	FrameMbsOnlyFlag                bool           `descriptor:"u(1)" json:"frame_mbs_only_flag"`
	MbAdaptiveFrameFieldFlag        bool           `descriptor:"u(1)" json:"mb_adaptive_frame_field_flag"`
	Direct8x8InferenceFlag          bool           `descriptor:"u(1)" json:"direct_8x8_inference_flag"`
	FrameCroppingFlag               bool           `descriptor:"u(1)" json:"frame_cropping_flag"`
	FrameCropLeftOffset             uint64         `descriptor:"ue(v)" json:"frame_crop_left_offset"`
	FrameCropRightOffset            uint64         `descriptor:"ue(v)" json:"frame_crop_right_offset"`
	FrameCropTopOffset              uint64         `descriptor:"ue(v)" json:"frame_crop_top_offset"`
	FrameCropBottomOffset           uint64         `descriptor:"ue(v)" json:"frame_crop_bottom_offset"`
	VuiParametersPresentFlag        bool           `descriptor:"u(1)" json:"vui_parameters_present_flag"`
	VuiParameters                   *VuiParameters `json:"vui_parameters,omitempty"`
}

func (e *SeqParameterSet) Read(d bits.Decoder) error {
	d.Decode(e, "ProfileIdc")
	d.Decode(e, "ConstraintSet0Flag")
	d.Decode(e, "ConstraintSet1Flag")
	d.Decode(e, "ConstraintSet2Flag")
	d.Decode(e, "ConstraintSet3Flag")
	d.Decode(e, "ConstraintSet4Flag")
	d.Decode(e, "ConstraintSet5Flag")
	d.Decode(e, "ReservedZero2bits")
	d.Decode(e, "LevelIdc")
	d.Decode(e, "SeqParameterSetId")
	if e.ProfileIdc == 100 || e.ProfileIdc == 110 || e.ProfileIdc == 122 || e.ProfileIdc == 244 || e.ProfileIdc == 44 || e.ProfileIdc == 83 || e.ProfileIdc == 86 || e.ProfileIdc == 118 || e.ProfileIdc == 128 || e.ProfileIdc == 138 || e.ProfileIdc == 139 || e.ProfileIdc == 134 || e.ProfileIdc == 135 {
		d.Decode(e, "ChromaFormatIdc")
		if e.ChromaFormatIdc == 3 {
			d.Decode(e, "SeparateColourPlaneFlag")
		}
		d.Decode(e, "BitDepthLumaMinus8")
		d.Decode(e, "BitDepthChromaMinus8")
		d.Decode(e, "QpprimeYZeroTransformBypassFlag")
		d.Decode(e, "SeqScalingMatrixPresentFlag")
		if e.SeqScalingMatrixPresentFlag {
			len := 0
			if e.ChromaFormatIdc != 3 {
				len = 8
			} else {
				len = 12
			}
			for i := 0; i < len; i++ {
				d.DecodeIndex(e, "SeqScalingListPresentFlag", i)
				if e.SeqScalingListPresentFlag[i] {
					// TODO -- add these fields, adapt ScalingList to provide output parameters, don't store in e
					/*
					   if i < 6 {
					       e.ScalingList = &ScalingList{}
					       e.ScalingList.Read(d, ScalingList4x4[i], 16, UseDefaultScalingMatrix4x4Flag[i])
					   } else {
					       e.ScalingList = &ScalingList{}
					       e.ScalingList.Read(d, ScalingList8x8[i−6], 64, UseDefaultScalingMatrix8x8Flag[i−6])
					   }
					*/
				}
			}
		}
	}
	d.Decode(e, "Log2MaxFrameNumMinus4")
	d.Decode(e, "PicOrderCntType")
	if e.PicOrderCntType == 0 {
		d.Decode(e, "Log2MaxPicOrderCntLsbMinus4")
	} else if e.PicOrderCntType == 1 {
		d.Decode(e, "DeltaPicOrderAlwaysZeroFlag")
		d.Decode(e, "OffsetForNonRefPic")
		d.Decode(e, "OffsetForTopToBottomField")
		d.Decode(e, "NumRefFramesInPicOrderCntCycle")
		for i := 0; i < int(e.NumRefFramesInPicOrderCntCycle); i++ {
			d.DecodeIndex(e, "OffsetForRefFrame", i)
		}
	}
	d.Decode(e, "MaxNumRefFrames")
	d.Decode(e, "GapsInFrameNumValueAllowedFlag")
	d.Decode(e, "PicWidthInMbsMinus1")
	d.Decode(e, "PicHeightInMapUnitsMinus1")
	d.Decode(e, "FrameMbsOnlyFlag")
	if !e.FrameMbsOnlyFlag {
		d.Decode(e, "MbAdaptiveFrameFieldFlag")
	}
	d.Decode(e, "Direct8x8InferenceFlag")
	d.Decode(e, "FrameCroppingFlag")
	if e.FrameCroppingFlag {
		d.Decode(e, "FrameCropLeftOffset")
		d.Decode(e, "FrameCropRightOffset")
		d.Decode(e, "FrameCropTopOffset")
		d.Decode(e, "FrameCropBottomOffset")
	}
	d.Decode(e, "VuiParametersPresentFlag")
	if e.VuiParametersPresentFlag {
		e.VuiParameters = &VuiParameters{}
		e.VuiParameters.Read(d)
	}
	return d.Error()
}
