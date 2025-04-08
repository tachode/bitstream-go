package h264

import (
	"github.com/tachode/bitstream-go/bits"
)

type SeqParameterSetMvcExtension struct {
	NumViewsMinus1                   uint64       `descriptor:"ue(v)" json:"num_views_minus1"`
	ViewId                           []uint64     `descriptor:"ue(v)" json:"view_id"`
	NumAnchorRefsL0                  []uint64     `descriptor:"ue(v)" json:"num_anchor_refs_l0"`
	AnchorRefL0                      [][]uint64   `descriptor:"ue(v)" json:"anchor_ref_l0"`
	NumAnchorRefsL1                  []uint64     `descriptor:"ue(v)" json:"num_anchor_refs_l1"`
	AnchorRefL1                      [][]uint64   `descriptor:"ue(v)" json:"anchor_ref_l1"`
	NumNonAnchorRefsL0               []uint64     `descriptor:"ue(v)" json:"num_non_anchor_refs_l0"`
	NonAnchorRefL0                   [][]uint64   `descriptor:"ue(v)" json:"non_anchor_ref_l0"`
	NumNonAnchorRefsL1               []uint64     `descriptor:"ue(v)" json:"num_non_anchor_refs_l1"`
	NonAnchorRefL1                   [][]uint64   `descriptor:"ue(v)" json:"non_anchor_ref_l1"`
	NumLevelValuesSignalledMinus1    uint64       `descriptor:"ue(v)" json:"num_level_values_signalled_minus1"`
	LevelIdc                         []uint8      `descriptor:"u(8)" json:"level_idc"`
	NumApplicableOpsMinus1           []uint64     `descriptor:"ue(v)" json:"num_applicable_ops_minus1"`
	ApplicableOpTemporalId           [][]uint8    `descriptor:"u(3)" json:"applicable_op_temporal_id"`
	ApplicableOpNumTargetViewsMinus1 [][]uint64   `descriptor:"ue(v)" json:"applicable_op_num_target_views_minus1"`
	ApplicableOpTargetViewId         [][][]uint64 `descriptor:"ue(v)" json:"applicable_op_target_view_id"`
	ApplicableOpNumViewsMinus1       [][]uint64   `descriptor:"ue(v)" json:"applicable_op_num_views_minus1"`
	MfcFormatIdc                     uint8        `descriptor:"u(6)" json:"mfc_format_idc"`
	DefaultGridPositionFlag          bool         `descriptor:"u(1)" json:"default_grid_position_flag"`
	View0GridPositionX               uint8        `descriptor:"u(4)" json:"view0_grid_position_x"`
	View0GridPositionY               uint8        `descriptor:"u(4)" json:"view0_grid_position_y"`
	View1GridPositionX               uint8        `descriptor:"u(4)" json:"view1_grid_position_x"`
	View1GridPositionY               uint8        `descriptor:"u(4)" json:"view1_grid_position_y"`
	RpuFilterEnabledFlag             bool         `descriptor:"u(1)" json:"rpu_filter_enabled_flag"`
	RpuFieldProcessingFlag           bool         `descriptor:"u(1)" json:"rpu_field_processing_flag"`
}

func (e *SeqParameterSetMvcExtension) Read(d bits.Decoder) error {
	d.Decode(e, "NumViewsMinus1")
	for i := 0; i <= int(e.NumViewsMinus1); i++ {
		d.DecodeIndex(e, "ViewId", i)
	}
	for i := 1; i <= int(e.NumViewsMinus1); i++ {
		d.DecodeIndex(e, "NumAnchorRefsL0", i)
		for j := 0; j < int(e.NumAnchorRefsL0[i]); j++ {
			d.DecodeIndex(e, "AnchorRefL0", i, j)
		}
		d.DecodeIndex(e, "NumAnchorRefsL1", i)
		for j := 0; j < int(e.NumAnchorRefsL1[i]); j++ {
			d.DecodeIndex(e, "AnchorRefL1", i, j)
		}
	}
	for i := 1; i <= int(e.NumViewsMinus1); i++ {
		d.DecodeIndex(e, "NumNonAnchorRefsL0", i)
		for j := 0; j < int(e.NumNonAnchorRefsL0[i]); j++ {
			d.DecodeIndex(e, "NonAnchorRefL0", i, j)
		}
		d.DecodeIndex(e, "NumNonAnchorRefsL1", i)
		for j := 0; j < int(e.NumNonAnchorRefsL1[i]); j++ {
			d.DecodeIndex(e, "NonAnchorRefL1", i, j)
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
			}
			d.DecodeIndex(e, "ApplicableOpNumViewsMinus1", i, j)
		}
	}
	profile_idc := d.Value("ProfileIdc").(uint64)
	frame_mbs_only_flag := d.Value("FrameMbsOnlyFlag").(bool)
	if profile_idc == 134 {
		d.Decode(e, "MfcFormatIdc")
		if e.MfcFormatIdc == 0 || e.MfcFormatIdc == 1 {
			d.Decode(e, "DefaultGridPositionFlag")
			if !e.DefaultGridPositionFlag {
				d.Decode(e, "View0GridPositionX")
				d.Decode(e, "View0GridPositionY")
				d.Decode(e, "View1GridPositionX")
				d.Decode(e, "View1GridPositionY")
			}
		}
		d.Decode(e, "RpuFilterEnabledFlag")
		if !frame_mbs_only_flag {
			d.Decode(e, "RpuFieldProcessingFlag")
		}
	}
	return d.Error()
}
