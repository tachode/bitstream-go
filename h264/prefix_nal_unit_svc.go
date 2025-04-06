package h264

import "github.com/tachode/bitstream-go/bits"

type PrefixNalUnitSvc struct {
	StoreRefBasePicFlag                      bool                  `descriptor:"u(1)" json:"store_ref_base_pic_flag"`
	DecRefBasePicMarking                     *DecRefBasePicMarking `json:"dec_ref_base_pic_marking,omitempty"`
	AdditionalPrefixNalUnitExtensionFlag     bool                  `descriptor:"u(1)" json:"additional_prefix_nal_unit_extension_flag"`
	AdditionalPrefixNalUnitExtensionDataFlag bool                  `descriptor:"u(1)" json:"additional_prefix_nal_unit_extension_data_flag"`
}

func (e *PrefixNalUnitSvc) Read(d bits.Decoder) error {
	nal_ref_idc, _ := d.Value("NalRefIdc").(uint64)
	idr_flag, _ := d.Value("IdrFlag").(bool)
	use_ref_base_pic_flag := d.Value("UseRefBasePicFlag").(bool)

	if nal_ref_idc != 0 {
		d.Decode(e, "StoreRefBasePicFlag")
		if (use_ref_base_pic_flag || e.StoreRefBasePicFlag) && !idr_flag {
			e.DecRefBasePicMarking = &DecRefBasePicMarking{}
			e.DecRefBasePicMarking.Read(d)
		}
		d.Decode(e, "AdditionalPrefixNalUnitExtensionFlag")
		if e.AdditionalPrefixNalUnitExtensionFlag {
			for d.MoreRbspData() {
				d.Decode(e, "AdditionalPrefixNalUnitExtensionDataFlag")
			}
		}
	} else if d.MoreRbspData() {
		for d.MoreRbspData() {
			d.Decode(e, "AdditionalPrefixNalUnitExtensionDataFlag")
		}
	}
	return d.Error()
}
