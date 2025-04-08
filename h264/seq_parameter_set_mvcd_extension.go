package h264

import (
	"fmt"

	"github.com/tachode/bitstream-go/bits"
)

type SeqParameterSetMvcdExtension struct {
	NumViewsMinus1                    uint64                      `descriptor:"ue(v)" json:"num_views_minus1"`
	ViewId                            []uint64                    `descriptor:"ue(v)" json:"view_id"`
	DepthViewPresentFlag              []bool                      `descriptor:"u(1)" json:"depth_view_present_flag"`
	TextureViewPresentFlag            []bool                      `descriptor:"u(1)" json:"texture_view_present_flag"`
	NumAnchorRefsL0                   []uint64                    `descriptor:"ue(v)" json:"num_anchor_refs_l0"`
	AnchorRefL0                       [][]uint64                  `descriptor:"ue(v)" json:"anchor_ref_l0"`
	NumAnchorRefsL1                   []uint64                    `descriptor:"ue(v)" json:"num_anchor_refs_l1"`
	AnchorRefL1                       [][]uint64                  `descriptor:"ue(v)" json:"anchor_ref_l1"`
	NumNonAnchorRefsL0                []uint64                    `descriptor:"ue(v)" json:"num_non_anchor_refs_l0"`
	NonAnchorRefL0                    [][]uint64                  `descriptor:"ue(v)" json:"non_anchor_ref_l0"`
	NumNonAnchorRefsL1                []uint64                    `descriptor:"ue(v)" json:"num_non_anchor_refs_l1"`
	NonAnchorRefL1                    [][]uint64                  `descriptor:"ue(v)" json:"non_anchor_ref_l1"`
	NumLevelValuesSignalledMinus1     uint64                      `descriptor:"ue(v)" json:"num_level_values_signalled_minus1"`
	LevelIdc                          []uint8                     `descriptor:"u(8)" json:"level_idc"`
	NumApplicableOpsMinus1            []uint64                    `descriptor:"ue(v)" json:"num_applicable_ops_minus1"`
	ApplicableOpTemporalId            [][]uint8                   `descriptor:"u(3)" json:"applicable_op_temporal_id"`
	ApplicableOpNumTargetViewsMinus1  [][]uint64                  `descriptor:"ue(v)" json:"applicable_op_num_target_views_minus1"`
	ApplicableOpTargetViewId          [][][]uint64                `descriptor:"ue(v)" json:"applicable_op_target_view_id"`
	ApplicableOpDepthFlag             [][][]bool                  `descriptor:"u(1)" json:"applicable_op_depth_flag"`
	ApplicableOpTextureFlag           [][][]bool                  `descriptor:"u(1)" json:"applicable_op_texture_flag"`
	ApplicableOpNumTextureViewsMinus1 [][]uint64                  `descriptor:"ue(v)" json:"applicable_op_num_texture_views_minus1"`
	ApplicableOpNumDepthViews         [][]uint64                  `descriptor:"ue(v)" json:"applicable_op_num_depth_views"`
	MvcdVuiParametersPresentFlag      bool                        `descriptor:"u(1)" json:"mvcd_vui_parameters_present_flag"`
	MvcdVuiParametersExtension        *MvcdVuiParametersExtension `json:"mvcd_vui_parameters_extension,omitempty"`
	TextureVuiParametersPresentFlag   bool                        `descriptor:"u(1)" json:"texture_vui_parameters_present_flag"`
	MvcVuiParametersExtension         *MvcVuiParametersExtension  `json:"mvc_vui_parameters_extension,omitempty"`
}

func (e *SeqParameterSetMvcdExtension) Read(d bits.Decoder) error {
	d.Decode(e, "NumViewsMinus1")
	NumDepthViews := 0
	for i := 0; i <= int(e.NumViewsMinus1); i++ {
		d.DecodeIndex(e, "ViewId", i)
		d.DecodeIndex(e, "DepthViewPresentFlag", i)
		d.SetValue(fmt.Sprintf("DepthViewId[%d]", NumDepthViews), e.ViewId[i])
		if e.DepthViewPresentFlag[i] {
			NumDepthViews++
		}
		d.DecodeIndex(e, "TextureViewPresentFlag", i)
	}
	d.SetValue("NumDepthViews", uint64(NumDepthViews))
	for i := 1; i <= int(e.NumViewsMinus1); i++ {
		if e.DepthViewPresentFlag[i] {
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
	for i := 1; i <= int(e.NumViewsMinus1); i++ {
		if e.DepthViewPresentFlag[i] {
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
	d.Decode(e, "NumLevelValuesSignalledMinus1")
	for i := 0; i <= int(e.NumLevelValuesSignalledMinus1); i++ {
		d.DecodeIndex(e, "LevelIdc", i)
		d.DecodeIndex(e, "NumApplicableOpsMinus1", i)
		for j := 0; j <= int(e.NumApplicableOpsMinus1[i]); j++ {
			d.DecodeIndex(e, "ApplicableOpTemporalId", i, j)
			d.DecodeIndex(e, "ApplicableOpNumTargetViewsMinus1", i, j)
			for k := 0; k <= int(e.ApplicableOpNumTargetViewsMinus1[i][j]); k++ {
				d.DecodeIndex(e, "ApplicableOpTargetViewId", i, j, k)
				d.DecodeIndex(e, "ApplicableOpDepthFlag", i, j, k)
				d.DecodeIndex(e, "ApplicableOpTextureFlag", i, j, k)
			}
			d.DecodeIndex(e, "ApplicableOpNumTextureViewsMinus1", i, j)
			d.DecodeIndex(e, "ApplicableOpNumDepthViews", i, j)
		}
	}
	d.Decode(e, "MvcdVuiParametersPresentFlag")
	if e.MvcdVuiParametersPresentFlag {
		e.MvcdVuiParametersExtension = &MvcdVuiParametersExtension{}
		e.MvcdVuiParametersExtension.Read(d)
	}
	d.Decode(e, "TextureVuiParametersPresentFlag")
	if e.TextureVuiParametersPresentFlag {
		e.MvcVuiParametersExtension = &MvcVuiParametersExtension{}
		e.MvcVuiParametersExtension.Read(d)
	}
	return d.Error()
}
