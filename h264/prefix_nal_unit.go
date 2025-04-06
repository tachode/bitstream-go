package h264

import (
	"github.com/tachode/bitstream-go/bits"
)

func init() { RegisterNalPayloadType(NalUnitTypePrefixNalUnit, &PrefixNalUnit{}) }

type PrefixNalUnit struct {
	PrefixNalUnitSvc *PrefixNalUnitSvc `json:"prefix_nal_unit_svc,omitempty"`
}

func (e *PrefixNalUnit) Read(d bits.Decoder) error {
	svc_extension_flag, _ := d.Value("SvcExtensionFlag").(bool)
	if svc_extension_flag {
		e.PrefixNalUnitSvc = &PrefixNalUnitSvc{}
		e.PrefixNalUnitSvc.Read(d)
	}
	return d.Error()
}
