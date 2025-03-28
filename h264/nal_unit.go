package h264

import (
	"bytes"
	"fmt"

	"github.com/tachode/bitstream-go/bits"
)

type NalUnit struct {
	ForbiddenZeroBit            bool                         `descriptor:"f(1)=0" json:"forbidden_zero_bit"`
	NalRefIdc                   uint8                        `descriptor:"u(2)" json:"nal_ref_idc"`
	NalUnitType                 NalUnitType                  `descriptor:"u(5)" json:"nal_unit_type"`
	SvcExtensionFlag            bool                         `descriptor:"u(1)" json:"svc_extension_flag"`
	Avc3dExtensionFlag          bool                         `descriptor:"u(1)" json:"avc_3d_extension_flag"`
	NalUnitHeaderSvcExtension   *NalUnitHeaderSvcExtension   `json:"nal_unit_header_svc_extension,omitempty"`
	NalUnitHeader3davcExtension *NalUnitHeader3davcExtension `json:"nal_unit_header_3davc_extension,omitempty"`
	NalUnitHeaderMvcExtension   *NalUnitHeaderMvcExtension   `json:"nal_unit_header_mvc_extension,omitempty"`
	RbspByte                    []byte                       `descriptor:"b(8)" json:"rbsp_byte"`
	Payload                     any                          `json:"payload"`
}

func (e *NalUnit) Read(d bits.Decoder, NumBytesInNALunit int) error {
	d.Decode(e, "ForbiddenZeroBit")
	d.Decode(e, "NalRefIdc")
	d.Decode(e, "NalUnitType")
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

	if d.Error() != nil {
		return d.Error()
	}

	// The bitstream will be byte-aligned at this point, so we can use
	// normal byte operations to handle RBSP
	bytesRemaining := NumBytesInNALunit - nalUnitHeaderBytes
	if bytesRemaining <= 0 {
		return fmt.Errorf("cannot parse nal: payload computes as %d bytes long", bytesRemaining)
	}
	payloadBytes := make([]byte, bytesRemaining)
	_, err := d.Read(payloadBytes)
	if err != nil {
		return err
	}
	e.RbspByte = Unescape(payloadBytes)

	return d.Error()
}

func Unescape(src []byte) []byte {
	dst := make([]byte, 0, len(src))
	var emulationPreventionSequence = []byte{0, 0, 3}
	for i := bytes.Index(src, emulationPreventionSequence); i != -1; i = bytes.Index(src, emulationPreventionSequence) {
		dst = append(dst, src[:i+2]...)
		src = src[i+3:]
	}
	dst = append(dst, src...)
	return dst
}
