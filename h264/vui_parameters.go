package h264

import "github.com/tachode/bitstream-go/bits"

type VuiParameters struct {
	AspectRatioInfoPresentFlag         bool           `descriptor:"u(1)" json:"aspect_ratio_info_present_flag"`
	AspectRatioIdc                     uint8          `descriptor:"u(8)" json:"aspect_ratio_idc"`
	SarWidth                           uint16         `descriptor:"u(16)" json:"sar_width"`
	SarHeight                          uint16         `descriptor:"u(16)" json:"sar_height"`
	OverscanInfoPresentFlag            bool           `descriptor:"u(1)" json:"overscan_info_present_flag"`
	OverscanAppropriateFlag            bool           `descriptor:"u(1)" json:"overscan_appropriate_flag"`
	VideoSignalTypePresentFlag         bool           `descriptor:"u(1)" json:"video_signal_type_present_flag"`
	VideoFormat                        uint8          `descriptor:"u(3)" json:"video_format"`
	VideoFullRangeFlag                 bool           `descriptor:"u(1)" json:"video_full_range_flag"`
	ColourDescriptionPresentFlag       bool           `descriptor:"u(1)" json:"colour_description_present_flag"`
	ColourPrimaries                    uint8          `descriptor:"u(8)" json:"colour_primaries"`
	TransferCharacteristics            uint8          `descriptor:"u(8)" json:"transfer_characteristics"`
	MatrixCoefficients                 uint8          `descriptor:"u(8)" json:"matrix_coefficients"`
	ChromaLocInfoPresentFlag           bool           `descriptor:"u(1)" json:"chroma_loc_info_present_flag"`
	ChromaSampleLocTypeTopField        uint64         `descriptor:"ue(v)" json:"chroma_sample_loc_type_top_field"`
	ChromaSampleLocTypeBottomField     uint64         `descriptor:"ue(v)" json:"chroma_sample_loc_type_bottom_field"`
	TimingInfoPresentFlag              bool           `descriptor:"u(1)" json:"timing_info_present_flag"`
	NumUnitsInTick                     uint32         `descriptor:"u(32)" json:"num_units_in_tick"`
	TimeScale                          uint32         `descriptor:"u(32)" json:"time_scale"`
	FixedFrameRateFlag                 bool           `descriptor:"u(1)" json:"fixed_frame_rate_flag"`
	NalHrdParametersPresentFlag        bool           `descriptor:"u(1)" json:"nal_hrd_parameters_present_flag"`
	NalHrdParameters                   *HrdParameters `json:"nal_hrd_parameters,omitempty"`
	VclHrdParametersPresentFlag        bool           `descriptor:"u(1)" json:"vcl_hrd_parameters_present_flag"`
	VclHrdParameters                   *HrdParameters `json:"vcl_hrd_parameters,omitempty"`
	LowDelayHrdFlag                    bool           `descriptor:"u(1)" json:"low_delay_hrd_flag"`
	PicStructPresentFlag               bool           `descriptor:"u(1)" json:"pic_struct_present_flag"`
	BitstreamRestrictionFlag           bool           `descriptor:"u(1)" json:"bitstream_restriction_flag"`
	MotionVectorsOverPicBoundariesFlag bool           `descriptor:"u(1)" json:"motion_vectors_over_pic_boundaries_flag"`
	MaxBytesPerPicDenom                uint64         `descriptor:"ue(v)" json:"max_bytes_per_pic_denom"`
	MaxBitsPerMbDenom                  uint64         `descriptor:"ue(v)" json:"max_bits_per_mb_denom"`
	Log2MaxMvLengthHorizontal          uint64         `descriptor:"ue(v)" json:"log2_max_mv_length_horizontal"`
	Log2MaxMvLengthVertical            uint64         `descriptor:"ue(v)" json:"log2_max_mv_length_vertical"`
	MaxNumReorderFrames                uint64         `descriptor:"ue(v)" json:"max_num_reorder_frames"`
	MaxDecFrameBuffering               uint64         `descriptor:"ue(v)" json:"max_dec_frame_buffering"`
}

func (e *VuiParameters) Read(d bits.Decoder) error {
	const Extended_SAR = 255
	d.Decode(e, "AspectRatioInfoPresentFlag")
	if e.AspectRatioInfoPresentFlag {
		d.Decode(e, "AspectRatioIdc")
		if e.AspectRatioIdc == Extended_SAR {
			d.Decode(e, "SarWidth")
			d.Decode(e, "SarHeight")
		}
	}
	d.Decode(e, "OverscanInfoPresentFlag")
	if e.OverscanInfoPresentFlag {
		d.Decode(e, "OverscanAppropriateFlag")
	}
	d.Decode(e, "VideoSignalTypePresentFlag")
	if e.VideoSignalTypePresentFlag {
		d.Decode(e, "VideoFormat")
		d.Decode(e, "VideoFullRangeFlag")
		d.Decode(e, "ColourDescriptionPresentFlag")
		if e.ColourDescriptionPresentFlag {
			d.Decode(e, "ColourPrimaries")
			d.Decode(e, "TransferCharacteristics")
			d.Decode(e, "MatrixCoefficients")
		}
	}
	d.Decode(e, "ChromaLocInfoPresentFlag")
	if e.ChromaLocInfoPresentFlag {
		d.Decode(e, "ChromaSampleLocTypeTopField")
		d.Decode(e, "ChromaSampleLocTypeBottomField")
	}
	d.Decode(e, "TimingInfoPresentFlag")
	if e.TimingInfoPresentFlag {
		d.Decode(e, "NumUnitsInTick")
		d.Decode(e, "TimeScale")
		d.Decode(e, "FixedFrameRateFlag")
	}
	d.Decode(e, "NalHrdParametersPresentFlag")
	if e.NalHrdParametersPresentFlag {
		e.NalHrdParameters = &HrdParameters{}
		e.NalHrdParameters.Read(d)
	}
	d.Decode(e, "VclHrdParametersPresentFlag")
	if e.VclHrdParametersPresentFlag {
		e.VclHrdParameters = &HrdParameters{}
		e.VclHrdParameters.Read(d)
	}
	if e.NalHrdParametersPresentFlag || e.VclHrdParametersPresentFlag {
		d.Decode(e, "LowDelayHrdFlag")
	}
	d.Decode(e, "PicStructPresentFlag")
	d.Decode(e, "BitstreamRestrictionFlag")
	if e.BitstreamRestrictionFlag {
		d.Decode(e, "MotionVectorsOverPicBoundariesFlag")
		d.Decode(e, "MaxBytesPerPicDenom")
		d.Decode(e, "MaxBitsPerMbDenom")
		d.Decode(e, "Log2MaxMvLengthHorizontal")
		d.Decode(e, "Log2MaxMvLengthVertical")
		d.Decode(e, "MaxNumReorderFrames")
		d.Decode(e, "MaxDecFrameBuffering")
	}
	return d.Error()
}
