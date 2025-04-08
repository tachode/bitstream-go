package h264

import "github.com/tachode/bitstream-go/bits"

func init() { RegisterNalPayloadType(NalUnitTypeDepthParameterSet, &DepthParameterSet{}) }

type DepthParameterSet struct {
	DepthParameterSetId                   uint64       `descriptor:"ue(v)" json:"depth_parameter_set_id"`
	PredDirection                         uint64       `descriptor:"ue(v)" json:"pred_direction"`
	RefDpsId0                             uint64       `descriptor:"ue(v)" json:"ref_dps_id0"`
	RefDpsId1                             uint64       `descriptor:"ue(v)" json:"ref_dps_id1"`
	PredWeight0                           uint8        `descriptor:"u(6)" json:"pred_weight0"`
	NumDepthViewsMinus1                   uint64       `descriptor:"ue(v)" json:"num_depth_views_minus1"`
	DepthRanges                           *DepthRanges `json:"depth_ranges,omitempty"`
	VspParamFlag                          bool         `descriptor:"u(1)" json:"vsp_param_flag"`
	VspParam                              *VspParam    `json:"vsp_param,omitempty"`
	DepthParamAdditionalExtensionFlag     bool         `descriptor:"u(1)" json:"depth_param_additional_extension_flag"`
	NonlinearDepthRepresentationNum       uint64       `descriptor:"ue(v)" json:"nonlinear_depth_representation_num"`
	NonlinearDepthRepresentationModel     []uint64     `descriptor:"ue(v)" json:"nonlinear_depth_representation_model"`
	DepthParamAdditionalExtensionDataFlag bool         `descriptor:"u(1)" json:"depth_param_additional_extension_data_flag"`
}

func (e *DepthParameterSet) Read(d bits.Decoder) error {
	d.Decode(e, "DepthParameterSetId")
	d.Decode(e, "PredDirection")
	if e.PredDirection == 0 || e.PredDirection == 1 {
		d.Decode(e, "RefDpsId0")
		predWeight0 := 64
		d.SetValue("PredWeight0", uint8(predWeight0))
	}
	if e.PredDirection == 0 {
		d.Decode(e, "RefDpsId1")
		d.Decode(e, "PredWeight0")
	}
	d.Decode(e, "NumDepthViewsMinus1")
	e.DepthRanges = &DepthRanges{}
	e.DepthRanges.Read(d, int(e.NumDepthViewsMinus1+1), int(e.PredDirection), int(e.DepthParameterSetId))
	d.Decode(e, "VspParamFlag")
	if e.VspParamFlag {
		e.VspParam = &VspParam{}
		e.VspParam.Read(d, int(e.NumDepthViewsMinus1+1), int(e.PredDirection), int(e.DepthParameterSetId))
	}
	d.Decode(e, "DepthParamAdditionalExtensionFlag")
	d.Decode(e, "NonlinearDepthRepresentationNum")
	for i := 1; i <= int(e.NonlinearDepthRepresentationNum); i++ {
		d.DecodeIndex(e, "NonlinearDepthRepresentationModel", i)
	}
	if e.DepthParamAdditionalExtensionFlag {
		for d.MoreRbspData() {
			d.Decode(e, "DepthParamAdditionalExtensionDataFlag")
		}
	}
	return d.Error()
}
