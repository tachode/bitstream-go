package h264

import "github.com/tachode/bitstream-go/bits"

func init() { RegisterNalPayload(NalUnitTypeSpsExt, &SeqParameterSetExtension{}) }

type SeqParameterSetExtension struct {
	SeqParameterSetId       uint64 `descriptor:"ue(v)" json:"seq_parameter_set_id"`
	AuxFormatIdc            uint64 `descriptor:"ue(v)" json:"aux_format_idc"`
	BitDepthAuxMinus8       uint64 `descriptor:"ue(v)" json:"bit_depth_aux_minus8"`
	AlphaIncrFlag           bool   `descriptor:"u(1)" json:"alpha_incr_flag"`
	AlphaOpaqueValue        uint64 `descriptor:"u(v)" json:"alpha_opaque_value"`
	AlphaTransparentValue   uint64 `descriptor:"u(v)" json:"alpha_transparent_value"`
	AdditionalExtensionFlag bool   `descriptor:"u(1)" json:"additional_extension_flag"`
}

func (e *SeqParameterSetExtension) Read(d bits.Decoder) error {
	d.Decode(e, "SeqParameterSetId")
	d.Decode(e, "AuxFormatIdc")
	if e.AuxFormatIdc != 0 {
		d.Decode(e, "BitDepthAuxMinus8")
		d.Decode(e, "AlphaIncrFlag")
		d.Decode(e, "AlphaOpaqueValue")
		d.Decode(e, "AlphaTransparentValue")
	}
	d.Decode(e, "AdditionalExtensionFlag")
	return d.Error()
}
