package h264

import (
	"fmt"

	"github.com/tachode/bitstream-go/bits"
)

type SeqParameterSet3davcExtension struct {
	Var3dvAcquisitionIdc           uint64       `descriptor:"ue(v)" json:"3dv_acquisition_idc"`
	ViewId3dv                      []uint64     `descriptor:"ue(v)" json:"view_id_3dv"`
	DepthRanges                    *DepthRanges `json:"depth_ranges,omitempty"`
	VspParam                       *VspParam    `json:"vsp_param,omitempty"`
	ReducedResolutionFlag          bool         `descriptor:"u(1)" json:"reduced_resolution_flag"`
	DepthPicWidthInMbsMinus1       uint64       `descriptor:"ue(v)" json:"depth_pic_width_in_mbs_minus1"`
	DepthPicHeightInMapUnitsMinus1 uint64       `descriptor:"ue(v)" json:"depth_pic_height_in_map_units_minus1"`
	DepthHorMultMinus1             uint64       `descriptor:"ue(v)" json:"depth_hor_mult_minus1"`
	DepthVerMultMinus1             uint64       `descriptor:"ue(v)" json:"depth_ver_mult_minus1"`
	DepthHorRsh                    uint64       `descriptor:"ue(v)" json:"depth_hor_rsh"`
	DepthVerRsh                    uint64       `descriptor:"ue(v)" json:"depth_ver_rsh"`
	DepthFrameCroppingFlag         bool         `descriptor:"u(1)" json:"depth_frame_cropping_flag"`
	DepthFrameCropLeftOffset       uint64       `descriptor:"ue(v)" json:"depth_frame_crop_left_offset"`
	DepthFrameCropRightOffset      uint64       `descriptor:"ue(v)" json:"depth_frame_crop_right_offset"`
	DepthFrameCropTopOffset        uint64       `descriptor:"ue(v)" json:"depth_frame_crop_top_offset"`
	DepthFrameCropBottomOffset     uint64       `descriptor:"ue(v)" json:"depth_frame_crop_bottom_offset"`
	GridPosNumViews                uint64       `descriptor:"ue(v)" json:"grid_pos_num_views"`
	GridPosViewId                  []uint64     `descriptor:"ue(v)" json:"grid_pos_view_id"`
	GridPosX                       []int64      `descriptor:"se(v)" json:"grid_pos_x"`
	GridPosY                       []int64      `descriptor:"se(v)" json:"grid_pos_y"`
	SliceHeaderPredictionFlag      bool         `descriptor:"u(1)" json:"slice_header_prediction_flag"`
	SeqViewSynthesisFlag           bool         `descriptor:"u(1)" json:"seq_view_synthesis_flag"`
	AlcSpsEnableFlag               bool         `descriptor:"u(1)" json:"alc_sps_enable_flag"`
	EnableRleSkipFlag              bool         `descriptor:"u(1)" json:"enable_rle_skip_flag"`
	NumAnchorRefsL0                []uint64     `descriptor:"ue(v)" json:"num_anchor_refs_l0"`
	AnchorRefL0                    [][]uint64   `descriptor:"ue(v)" json:"anchor_ref_l0"`
	NumAnchorRefsL1                []uint64     `descriptor:"ue(v)" json:"num_anchor_refs_l1"`
	AnchorRefL1                    [][]uint64   `descriptor:"ue(v)" json:"anchor_ref_l1"`
	NumNonAnchorRefsL0             []uint64     `descriptor:"ue(v)" json:"num_non_anchor_refs_l0"`
	NonAnchorRefL0                 [][]uint64   `descriptor:"ue(v)" json:"non_anchor_ref_l0"`
	NumNonAnchorRefsL1             []uint64     `descriptor:"ue(v)" json:"num_non_anchor_refs_l1"`
	NonAnchorRefL1                 [][]uint64   `descriptor:"ue(v)" json:"non_anchor_ref_l1"`
}

func (e *SeqParameterSet3davcExtension) Read(d bits.Decoder) error {
	NumDepthViews := int(d.Value("NumDepthViews").(uint64))
	if NumDepthViews > 0 {
		d.Decode(e, "Var3dvAcquisitionIdc")
		for i := 0; i < NumDepthViews; i++ {
			d.DecodeIndex(e, "ViewId3dv", i)
		}
		if e.Var3dvAcquisitionIdc != 0 {
			e.DepthRanges = &DepthRanges{}
			e.DepthRanges.Read(d, NumDepthViews, 2, 0)
			e.VspParam = &VspParam{}
			e.VspParam.Read(d, NumDepthViews, 2, 0)
		}
		d.Decode(e, "ReducedResolutionFlag")
		if e.ReducedResolutionFlag {
			d.Decode(e, "DepthPicWidthInMbsMinus1")
			d.Decode(e, "DepthPicHeightInMapUnitsMinus1")
			d.Decode(e, "DepthHorMultMinus1")
			d.Decode(e, "DepthVerMultMinus1")
			d.Decode(e, "DepthHorRsh")
			d.Decode(e, "DepthVerRsh")
		}
		d.Decode(e, "DepthFrameCroppingFlag")
		if e.DepthFrameCroppingFlag {
			d.Decode(e, "DepthFrameCropLeftOffset")
			d.Decode(e, "DepthFrameCropRightOffset")
			d.Decode(e, "DepthFrameCropTopOffset")
			d.Decode(e, "DepthFrameCropBottomOffset")
		}
		d.Decode(e, "GridPosNumViews")
		for i := 0; i < int(e.GridPosNumViews); i++ {
			d.DecodeIndex(e, "GridPosViewId", i)
			d.DecodeIndex(e, "GridPosX", int(e.GridPosViewId[i]))
			d.DecodeIndex(e, "GridPosY", int(e.GridPosViewId[i]))
		}
		d.Decode(e, "SliceHeaderPredictionFlag")
		d.Decode(e, "SeqViewSynthesisFlag")
	}
	d.Decode(e, "AlcSpsEnableFlag")
	d.Decode(e, "EnableRleSkipFlag")
	num_views_minus1, _ := d.Value("NumViewsMinus1").(uint64)

	// H.264 figure J-9
	AllViewsPairedFlag := true
	for i := 1; i <= int(num_views_minus1); i++ {
		texture_view_present_flag, _ := d.Value(fmt.Sprintf("TextureViewPresentFlag[%d]", i)).(bool)
		depth_view_present_flag, _ := d.Value(fmt.Sprintf("DepthViewPresentFlag[%d]", i)).(bool)
		AllViewsPairedFlag = (AllViewsPairedFlag && depth_view_present_flag && texture_view_present_flag)
	}

	if !AllViewsPairedFlag {
		for i := 1; i <= int(num_views_minus1); i++ {
			texture_view_present_flag, _ := d.Value(fmt.Sprintf("TextureViewPresentFlag[%d]", i)).(bool)
			if texture_view_present_flag {
				d.DecodeIndex(e, "NumAnchorRefsL0", i)
				for j := 0; j < int(e.NumAnchorRefsL0[i]); j++ {
					d.DecodeIndex(e, "AnchorRefL0", i, j)
				}
				d.DecodeIndex(e, "NumAnchorRefsL1", i)
				for j := 0; j < int(e.NumAnchorRefsL1[i]); j++ {
					d.DecodeIndex(e, "AnchorRefL1", i, j)
				}
			}
		}
		for i := 1; i <= int(num_views_minus1); i++ {
			texture_view_present_flag, _ := d.Value(fmt.Sprintf("TextureViewPresentFlag[%d]", i)).(bool)
			if texture_view_present_flag {
				d.DecodeIndex(e, "NumNonAnchorRefsL0", i)
				for j := 0; j < int(e.NumNonAnchorRefsL0[i]); j++ {
					d.DecodeIndex(e, "NonAnchorRefL0", i, j)
				}
				d.DecodeIndex(e, "NumNonAnchorRefsL1", i)
				for j := 0; j < int(e.NumNonAnchorRefsL1[i]); j++ {
					d.DecodeIndex(e, "NonAnchorRefL1", i, j)
				}
			}
		}
	}
	return d.Error()
}
