package h264

import (
	"github.com/tachode/bitstream-go/bits"
)

func init() { RegisterNalPayloadType(NalUnitTypeSubsetSeqParameterSet, &SubsetSeqParameterSet{}) }

type SubsetSeqParameterSet struct {
	SeqParameterSetData           *SeqParameterSet               `json:"seq_parameter_set_data,omitempty"`
	SeqParameterSetSvcExtension   *SeqParameterSetSvcExtension   `json:"seq_parameter_set_svc_extension,omitempty"`
	SvcVuiParametersPresentFlag   bool                           `descriptor:"u(1)" json:"svc_vui_parameters_present_flag"`
	SvcVuiParametersExtension     *SvcVuiParametersExtension     `json:"svc_vui_parameters_extension,omitempty"`
	BitEqualToOne                 bool                           `descriptor:"f(1)=1" json:"bit_equal_to_one"`
	SeqParameterSetMvcExtension   *SeqParameterSetMvcExtension   `json:"seq_parameter_set_mvc_extension,omitempty"`
	MvcVuiParametersPresentFlag   bool                           `descriptor:"u(1)" json:"mvc_vui_parameters_present_flag"`
	MvcVuiParametersExtension     *MvcVuiParametersExtension     `json:"mvc_vui_parameters_extension,omitempty"`
	SeqParameterSetMvcdExtension  *SeqParameterSetMvcdExtension  `json:"seq_parameter_set_mvcd_extension,omitempty"`
	SeqParameterSet3davcExtension *SeqParameterSet3davcExtension `json:"seq_parameter_set_3davc_extension,omitempty"`
	AdditionalExtension2Flag      bool                           `descriptor:"u(1)" json:"additional_extension2_flag"`
	AdditionalExtension2DataFlag  bool                           `descriptor:"u(1)" json:"additional_extension2_data_flag"`
}

func (e *SubsetSeqParameterSet) Read(d bits.Decoder) error {
	profile_idc, _ := d.Value("ProfileIdc").(uint64)
	e.SeqParameterSetData = &SeqParameterSet{}
	e.SeqParameterSetData.Read(d)
	if profile_idc == 83 || profile_idc == 86 {
		e.SeqParameterSetSvcExtension = &SeqParameterSetSvcExtension{}
		e.SeqParameterSetSvcExtension.Read(d)
		d.Decode(e, "SvcVuiParametersPresentFlag")
		if e.SvcVuiParametersPresentFlag {
			e.SvcVuiParametersExtension = &SvcVuiParametersExtension{}
			e.SvcVuiParametersExtension.Read(d)
		}
	} else if profile_idc == 118 || profile_idc == 128 || profile_idc == 134 {
		d.Decode(e, "BitEqualToOne")
		e.SeqParameterSetMvcExtension = &SeqParameterSetMvcExtension{}
		e.SeqParameterSetMvcExtension.Read(d)
		d.Decode(e, "MvcVuiParametersPresentFlag")
		if e.MvcVuiParametersPresentFlag {
			e.MvcVuiParametersExtension = &MvcVuiParametersExtension{}
			e.MvcVuiParametersExtension.Read(d)
		}
	} else if profile_idc == 138 || profile_idc == 135 {
		d.Decode(e, "BitEqualToOne")
		e.SeqParameterSetMvcdExtension = &SeqParameterSetMvcdExtension{}
		e.SeqParameterSetMvcdExtension.Read(d)
	} else if profile_idc == 139 {
		d.Decode(e, "BitEqualToOne")
		e.SeqParameterSetMvcdExtension = &SeqParameterSetMvcdExtension{}
		e.SeqParameterSetMvcdExtension.Read(d)
		e.SeqParameterSet3davcExtension = &SeqParameterSet3davcExtension{}
		e.SeqParameterSet3davcExtension.Read(d)
	}
	d.Decode(e, "AdditionalExtension2Flag")
	if e.AdditionalExtension2Flag {
		for d.MoreRbspData() {
			d.Decode(e, "AdditionalExtension2DataFlag")
		}
	}
	return d.Error()
}
