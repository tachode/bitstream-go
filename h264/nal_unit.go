package h264

import "github.com/tachode/bitstream-go/bits"

type NalUnit struct {
	ForbiddenZeroBit             bool                         `descriptor:"f(1)=0" json:"forbidden_zero_bit"`
	NalRefIdc                    uint8                        `descriptor:"u(2)" json:"nal_ref_idc"`
	NalUnitType                  NalUnitType                  `descriptor:"u(5)" json:"nal_unit_type"`
	SvcExtensionFlag             bool                         `descriptor:"u(1)" json:"svc_extension_flag"`
	Avc3dExtensionFlag           bool                         `descriptor:"u(1)" json:"avc_3d_extension_flag"`
	NalUnitHeaderSvcExtension    *NalUnitHeaderSvcExtension   `json:"nal_unit_header_svc_extension,omitempty"`
	NalUnitHeader3davcExtension  *NalUnitHeader3davcExtension `json:"nal_unit_header_3davc_extension,omitempty"`
	NalUnitHeaderMvcExtension    *NalUnitHeaderMvcExtension   `json:"nal_unit_header_mvc_extension,omitempty"`
	EmulationPreventionThreeByte uint8                        `descriptor:"f(8)=3" json:"emulation_prevention_three_byte"`
	RbspByte                     []byte                       `descriptor:"b(8)" json:"rbsp_byte"`
	Payload                      any                          `json:"payload"`
}

func (e *NalUnit) Read(d bits.Decoder, NumBytesInNALunit int) error {
	d.Decode(e, "ForbiddenZeroBit")
	d.Decode(e, "NalRefIdc")
	d.Decode(e, "NalUnitType")
	// NumBytesInRBSP := 0
	nalUnitHeaderBytes := 1
	if e.NalUnitType == 14 || e.NalUnitType == 20 || e.NalUnitType == 21 {
		if e.NalUnitType != 21 {
			d.Decode(e, "SvcExtensionFlag")
		} else {
			d.Decode(e, "Avc3dExtensionFlag")
		}
		if e.SvcExtensionFlag {
			e.NalUnitHeaderSvcExtension = &NalUnitHeaderSvcExtension{}
			e.NalUnitHeaderSvcExtension.Read(d)
			nalUnitHeaderBytes += 3
		} else if e.Avc3dExtensionFlag {
			e.NalUnitHeader3davcExtension = &NalUnitHeader3davcExtension{}
			e.NalUnitHeader3davcExtension.Read(d)
			nalUnitHeaderBytes += 2
		} else {
			e.NalUnitHeaderMvcExtension = &NalUnitHeaderMvcExtension{}
			e.NalUnitHeaderMvcExtension.Read(d)
			nalUnitHeaderBytes += 3
		}
	}
	// TODO
	/*
		for i := nalUnitHeaderBytes; i < NumBytesInNALunit; i++ {
			if i+2 < NumBytesInNALunit && next_bits(24) == 0x000003 {
				d.DecodeIndex(e, "RbspByte", NumBytesInRBSP)
				NumBytesInRBSP++
				d.DecodeIndex(e, "RbspByte", NumBytesInRBSP)
				NumBytesInRBSP++
				i += 2
				d.Decode(e, "EmulationPreventionThreeByte")
			} else {
				d.DecodeIndex(e, "RbspByte", NumBytesInRBSP)
				NumBytesInRBSP++
			}
		}
	*/
	return d.Error()
}
