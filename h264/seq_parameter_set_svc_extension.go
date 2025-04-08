package h264

import "github.com/tachode/bitstream-go/bits"

type SeqParameterSetSvcExtension struct {
	InterLayerDeblockingFilterControlPresentFlag bool  `descriptor:"u(1)" json:"inter_layer_deblocking_filter_control_present_flag"`
	ExtendedSpatialScalabilityIdc                uint8 `descriptor:"u(2)" json:"extended_spatial_scalability_idc"`
	ChromaPhaseXPlus1Flag                        bool  `descriptor:"u(1)" json:"chroma_phase_x_plus1_flag"`
	ChromaPhaseYPlus1                            uint8 `descriptor:"u(2)" json:"chroma_phase_y_plus1"`
	SeqRefLayerChromaPhaseXPlus1Flag             bool  `descriptor:"u(1)" json:"seq_ref_layer_chroma_phase_x_plus1_flag"`
	SeqRefLayerChromaPhaseYPlus1                 uint8 `descriptor:"u(2)" json:"seq_ref_layer_chroma_phase_y_plus1"`
	SeqScaledRefLayerLeftOffset                  int64 `descriptor:"se(v)" json:"seq_scaled_ref_layer_left_offset"`
	SeqScaledRefLayerTopOffset                   int64 `descriptor:"se(v)" json:"seq_scaled_ref_layer_top_offset"`
	SeqScaledRefLayerRightOffset                 int64 `descriptor:"se(v)" json:"seq_scaled_ref_layer_right_offset"`
	SeqScaledRefLayerBottomOffset                int64 `descriptor:"se(v)" json:"seq_scaled_ref_layer_bottom_offset"`
	SeqTcoeffLevelPredictionFlag                 bool  `descriptor:"u(1)" json:"seq_tcoeff_level_prediction_flag"`
	AdaptiveTcoeffLevelPredictionFlag            bool  `descriptor:"u(1)" json:"adaptive_tcoeff_level_prediction_flag"`
	SliceHeaderRestrictionFlag                   bool  `descriptor:"u(1)" json:"slice_header_restriction_flag"`
}

func (e *SeqParameterSetSvcExtension) Read(d bits.Decoder) error {
	ChromaArrayType, _ := d.Value("ChromaArrayType").(uint64)
	d.Decode(e, "InterLayerDeblockingFilterControlPresentFlag")
	d.Decode(e, "ExtendedSpatialScalabilityIdc")
	if ChromaArrayType == 1 || ChromaArrayType == 2 {
		d.Decode(e, "ChromaPhaseXPlus1Flag")
	}
	if ChromaArrayType == 1 {
		d.Decode(e, "ChromaPhaseYPlus1")
	}
	if e.ExtendedSpatialScalabilityIdc == 1 {
		if ChromaArrayType > 0 {
			d.Decode(e, "SeqRefLayerChromaPhaseXPlus1Flag")
			d.Decode(e, "SeqRefLayerChromaPhaseYPlus1")
		}
		d.Decode(e, "SeqScaledRefLayerLeftOffset")
		d.Decode(e, "SeqScaledRefLayerTopOffset")
		d.Decode(e, "SeqScaledRefLayerRightOffset")
		d.Decode(e, "SeqScaledRefLayerBottomOffset")
	}
	d.Decode(e, "SeqTcoeffLevelPredictionFlag")
	if e.SeqTcoeffLevelPredictionFlag {
		d.Decode(e, "AdaptiveTcoeffLevelPredictionFlag")
	}
	d.Decode(e, "SliceHeaderRestrictionFlag")
	return d.Error()
}
