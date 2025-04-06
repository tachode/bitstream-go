package h264

import "github.com/tachode/bitstream-go/bits"

func init() { RegisterNalPayloadType(NalUnitTypePPS, &PicParameterSet{}) }

type PicParameterSet struct {
	PicParameterSetId                     uint64       `descriptor:"ue(v)" json:"pic_parameter_set_id"`
	SeqParameterSetId                     uint64       `descriptor:"ue(v)" json:"seq_parameter_set_id"`
	EntropyCodingModeFlag                 bool         `descriptor:"u(1)" json:"entropy_coding_mode_flag"`
	BottomFieldPicOrderInFramePresentFlag bool         `descriptor:"u(1)" json:"bottom_field_pic_order_in_frame_present_flag"`
	NumSliceGroupsMinus1                  uint64       `descriptor:"ue(v)" json:"num_slice_groups_minus1"`
	SliceGroupMapType                     uint64       `descriptor:"ue(v)" json:"slice_group_map_type"`
	RunLengthMinus1                       []uint64     `descriptor:"ue(v)" json:"run_length_minus1"`
	TopLeft                               []uint64     `descriptor:"ue(v)" json:"top_left"`
	BottomRight                           []uint64     `descriptor:"ue(v)" json:"bottom_right"`
	SliceGroupChangeDirectionFlag         bool         `descriptor:"u(1)" json:"slice_group_change_direction_flag"`
	SliceGroupChangeRateMinus1            uint64       `descriptor:"ue(v)" json:"slice_group_change_rate_minus1"`
	PicSizeInMapUnitsMinus1               uint64       `descriptor:"ue(v)" json:"pic_size_in_map_units_minus1"`
	SliceGroupId                          []uint64     `descriptor:"u(v)" json:"slice_group_id"`
	NumRefIdxL0DefaultActiveMinus1        uint64       `descriptor:"ue(v)" json:"num_ref_idx_l0_default_active_minus1"`
	NumRefIdxL1DefaultActiveMinus1        uint64       `descriptor:"ue(v)" json:"num_ref_idx_l1_default_active_minus1"`
	WeightedPredFlag                      bool         `descriptor:"u(1)" json:"weighted_pred_flag"`
	WeightedBipredIdc                     uint8        `descriptor:"u(2)" json:"weighted_bipred_idc"`
	PicInitQpMinus26                      int64        `descriptor:"se(v)" json:"pic_init_qp_minus26"`
	PicInitQsMinus26                      int64        `descriptor:"se(v)" json:"pic_init_qs_minus26"`
	ChromaQpIndexOffset                   int64        `descriptor:"se(v)" json:"chroma_qp_index_offset"`
	DeblockingFilterControlPresentFlag    bool         `descriptor:"u(1)" json:"deblocking_filter_control_present_flag"`
	ConstrainedIntraPredFlag              bool         `descriptor:"u(1)" json:"constrained_intra_pred_flag"`
	RedundantPicCntPresentFlag            bool         `descriptor:"u(1)" json:"redundant_pic_cnt_present_flag"`
	Transform8x8ModeFlag                  bool         `descriptor:"u(1)" json:"transform_8x8_mode_flag"`
	PicScalingMatrixPresentFlag           bool         `descriptor:"u(1)" json:"pic_scaling_matrix_present_flag"`
	PicScalingListPresentFlag             []bool       `descriptor:"u(1)" json:"pic_scaling_list_present_flag"`
	ScalingList                           *ScalingList `json:"scaling_list,omitempty"`
	SecondChromaQpIndexOffset             int64        `descriptor:"se(v)" json:"second_chroma_qp_index_offset"`
}

func (e *PicParameterSet) Read(d bits.Decoder) error {
	d.Decode(e, "PicParameterSetId")
	d.Decode(e, "SeqParameterSetId")
	d.Decode(e, "EntropyCodingModeFlag")
	d.Decode(e, "BottomFieldPicOrderInFramePresentFlag")
	d.Decode(e, "NumSliceGroupsMinus1")
	if e.NumSliceGroupsMinus1 > 0 {
		d.Decode(e, "SliceGroupMapType")
		if e.SliceGroupMapType == 0 {
			for iGroup := 0; iGroup <= int(e.NumSliceGroupsMinus1); iGroup++ {
				d.DecodeIndex(e, "RunLengthMinus1", iGroup)
			}
		} else if e.SliceGroupMapType == 2 {
			for iGroup := 0; iGroup < int(e.NumSliceGroupsMinus1); iGroup++ {
				d.DecodeIndex(e, "TopLeft", iGroup)
				d.DecodeIndex(e, "BottomRight", iGroup)
			}
		} else if e.SliceGroupMapType == 3 || e.SliceGroupMapType == 4 || e.SliceGroupMapType == 5 {
			d.Decode(e, "SliceGroupChangeDirectionFlag")
			d.Decode(e, "SliceGroupChangeRateMinus1")
		} else if e.SliceGroupMapType == 6 {
			d.Decode(e, "PicSizeInMapUnitsMinus1")
			for i := 0; i <= int(e.PicSizeInMapUnitsMinus1); i++ {
				d.DecodeIndex(e, "SliceGroupId", i)
			}
		}
	}
	d.Decode(e, "NumRefIdxL0DefaultActiveMinus1")
	d.Decode(e, "NumRefIdxL1DefaultActiveMinus1")
	d.Decode(e, "WeightedPredFlag")
	d.Decode(e, "WeightedBipredIdc")
	d.Decode(e, "PicInitQpMinus26")
	d.Decode(e, "PicInitQsMinus26")
	d.Decode(e, "ChromaQpIndexOffset")
	d.Decode(e, "DeblockingFilterControlPresentFlag")
	d.Decode(e, "ConstrainedIntraPredFlag")
	d.Decode(e, "RedundantPicCntPresentFlag")
	chromaFormatIdc, _ := d.Value("ChromaFormatIdc").(uint64)
	if d.MoreRbspData() {
		d.Decode(e, "Transform8x8ModeFlag")
		d.Decode(e, "PicScalingMatrixPresentFlag")
		if e.PicScalingMatrixPresentFlag {
			e.ScalingList = &ScalingList{}
			for i := 0; i < 6+(If(chromaFormatIdc != 3, 2, 6))*If(e.Transform8x8ModeFlag, 1, 0); i++ {
				d.DecodeIndex(e, "PicScalingListPresentFlag", i)
				if e.PicScalingListPresentFlag[i] {
					if i < 6 {
						e.ScalingList.Read(d, i, 16)
					} else {
						e.ScalingList.Read(d, i-6, 64)
					}
				}
			}
		}
		d.Decode(e, "SecondChromaQpIndexOffset")
	}
	return d.Error()
}

func If[T any](condition bool, trueVal T, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}
