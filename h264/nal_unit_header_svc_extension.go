package h264

import "github.com/tachode/bitstream-go/bits"

type NalUnitHeaderSvcExtension struct {
	IdrFlag              bool  `descriptor:"u(1)" json:"idr_flag"`
	PriorityId           uint8 `descriptor:"u(6)" json:"priority_id"`
	NoInterLayerPredFlag bool  `descriptor:"u(1)" json:"no_inter_layer_pred_flag"`
	DependencyId         uint8 `descriptor:"u(3)" json:"dependency_id"`
	QualityId            uint8 `descriptor:"u(4)" json:"quality_id"`
	TemporalId           uint8 `descriptor:"u(3)" json:"temporal_id"`
	UseRefBasePicFlag    bool  `descriptor:"u(1)" json:"use_ref_base_pic_flag"`
	DiscardableFlag      bool  `descriptor:"u(1)" json:"discardable_flag"`
	OutputFlag           bool  `descriptor:"u(1)" json:"output_flag"`
	ReservedThree2bits   uint8 `descriptor:"u(2)" json:"reserved_three_2bits"`
}

func (e *NalUnitHeaderSvcExtension) Read(d bits.Decoder) error {
	d.DecodeRange(e, "IdrFlag", "ReservedThree2bits")
	return d.Error()
}
