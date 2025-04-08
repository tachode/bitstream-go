package h264

import "github.com/tachode/bitstream-go/bits"

type DepthRanges struct {
	ZNearFlag                 bool                         `descriptor:"u(1)" json:"z_near_flag"`
	ZFarFlag                  bool                         `descriptor:"u(1)" json:"z_far_flag"`
	Near3dvAcquisitionElement *Struct3dvAcquisitionElement `json:"near_3dv_acquisition_element,omitempty"`
	Far3dvAcquisitionElement  *Struct3dvAcquisitionElement `json:"far_3dv_acquisition_element,omitempty"`
	ZNearSign                 [][]uint8                    `json:"z_near_sign"`
	ZNearExp                  [][]uint64                   `json:"z_near_exp"`
	ZNearMantissa             [][]uint64                   `json:"z_near_mantissa"`
	ZNearManLen               [][]uint8                    `json:"z_near_man_len"`
	ZFarSign                  [][]uint8                    `json:"z_far_sign"`
	ZFarExp                   [][]uint64                   `json:"z_far_exp"`
	ZFarMantissa              [][]uint64                   `json:"z_far_mantissa"`
	ZFarManLen                [][]uint8                    `json:"z_far_man_len"`
}

func (e *DepthRanges) Read(d bits.Decoder, numViews int, predDirection int, index int) error {
	d.Decode(e, "ZNearFlag")
	d.Decode(e, "ZFarFlag")
	if e.ZNearFlag {
		e.Near3dvAcquisitionElement = &Struct3dvAcquisitionElement{}
		e.ZNearSign, e.ZNearExp, e.ZNearMantissa, e.ZNearManLen, _ =
			e.Near3dvAcquisitionElement.Read(d, numViews, predDirection, 7, index)
	}
	if e.ZFarFlag {
		e.Far3dvAcquisitionElement = &Struct3dvAcquisitionElement{}
		e.ZFarSign, e.ZFarExp, e.ZFarMantissa, e.ZFarManLen, _ =
			e.Far3dvAcquisitionElement.Read(d, numViews, predDirection, 7, index)
	}
	return d.Error()
}
