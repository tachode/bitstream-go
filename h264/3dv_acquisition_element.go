package h264

import (
	"fmt"

	"github.com/tachode/bitstream-go/bits"
)

type distance string

const (
	Near distance = "Near"
	Far  distance = "Far"
)

type Struct3dvAcquisitionElement struct {
	ElementEqualFlag  bool   `descriptor:"u(1)" json:"element_equal_flag"`
	MantissaLenMinus1 uint8  `descriptor:"u(5)" json:"mantissa_len_minus1"`
	Sign0             uint8  `descriptor:"u(1)" json:"sign0"`
	Exponent0         uint64 `descriptor:"u(v)" json:"exponent0"` // expLen bits
	Mantissa0         uint64 `descriptor:"u(v)" json:"mantissa0"` // manLen bits
	SkipFlag          bool   `descriptor:"u(1)" json:"skip_flag"`
	Sign1             uint8  `descriptor:"u(1)" json:"sign1"`
	ExponentSkipFlag  bool   `descriptor:"u(1)" json:"exponent_skip_flag"`
	Exponent1         uint64 `descriptor:"u(v)" json:"exponent1"` // expLen bits
	MantissaDiff      int64  `descriptor:"se(v)" json:"mantissa_diff"`
}

func (e *Struct3dvAcquisitionElement) Read(
	d bits.Decoder,
	numViews int,
	predDirection int,
	expLen int,
	index int,
	distance distance,
) (
	outSign [][]uint8,
	outExp [][]uint64,
	outMantissa [][]uint64,
	outManLen [][]uint8,
	err error,
) {
	ref_dps_id0, _ := d.Value("RefDpsId0").(uint64)
	ref_dps_id1, _ := d.Value("RefDpsId1").(uint64)
	predWeight0, _ := d.Value("PredWeight0").(uint8)

	// Values associated with this element
	outSign = make([][]uint8, index+1)
	outSign[index] = make([]uint8, numViews)
	outExp = make([][]uint64, index+1)
	outExp[index] = make([]uint64, numViews)
	outMantissa = make([][]uint64, index+1)
	outMantissa[index] = make([]uint64, numViews)
	outManLen = make([][]uint8, index+1)
	outManLen[index] = make([]uint8, numViews)

	// Global Values -- need to make sure there's enough space to hold
	// the new values
	Sign, _ := d.Value(fmt.Sprintf("Z%sSign", distance)).([][]uint8)
	Exp, _ := d.Value(fmt.Sprintf("Z%sExp", distance)).([][]uint64)
	Mantissa, _ := d.Value(fmt.Sprintf("Z%sMantissa", distance)).([][]uint64)
	ManLen, _ := d.Value(fmt.Sprintf("Z%sManLen", distance)).([][]uint8)

	// Make sure the global values are large enough for all the following operations
	for _, i := range []int{index, int(ref_dps_id0), int(ref_dps_id1)} {
		if Sign == nil {
			Sign = outSign
		}
		for len(Sign) < i {
			Sign = append(Sign, []uint8{})
		}
		for len(Sign[i]) < numViews {
			Sign[i] = append(Sign[i], 0)
		}
		if Exp == nil {
			Exp = outExp
		}
		for len(Exp) < i {
			Exp = append(Exp, []uint64{})
		}
		for len(Exp[i]) < numViews {
			Exp[i] = append(Exp[i], 0)
		}
		if Mantissa == nil {
			Mantissa = outMantissa
		}
		for len(Mantissa) < i {
			Mantissa = append(Mantissa, []uint64{})
		}
		for len(Mantissa[i]) < numViews {
			Mantissa[i] = append(Mantissa[i], 0)
		}
		if ManLen == nil {
			ManLen = outManLen
		}
		for len(ManLen) < i {
			ManLen = append(ManLen, []uint8{})
		}
		for len(ManLen[i]) < numViews {
			ManLen[i] = append(ManLen[i], 0)
		}
	}

	if numViews > 1 {
		d.Decode(e, "ElementEqualFlag")
	}
	var numValues int
	if !e.ElementEqualFlag {
		numValues = numViews
	} else {
		numValues = 1
	}
	for i := 0; i < numValues; i++ {
		var manLen int
		if predDirection == 2 && i == 0 {
			d.Decode(e, "MantissaLenMinus1")
			outManLen[index][i] = e.MantissaLenMinus1 + 1
			ManLen[index][i] = e.MantissaLenMinus1 + 1
			manLen = int(e.MantissaLenMinus1 + 1)
		}
		if predDirection == 2 {
			d.Decode(e, "Sign0")
			outSign[index][i] = e.Sign0
			Sign[index][i] = e.Sign0
			d.SetValueLength("Exponent0", expLen)
			d.Decode(e, "Exponent0")
			outExp[index][i] = e.Exponent0
			Exp[index][i] = e.Exponent0
			d.SetValueLength("Mantissa0", manLen)
			d.Decode(e, "Mantissa0")
			outMantissa[index][i] = e.Mantissa0
			Mantissa[index][i] = e.Mantissa0
		} else {
			d.Decode(e, "SkipFlag")
			if !e.SkipFlag {
				d.Decode(e, "Sign1")
				outSign[index][i] = e.Sign1
				Sign[index][i] = e.Sign1
				d.Decode(e, "ExponentSkipFlag")
				if !e.ExponentSkipFlag {
					d.SetValueLength("Exponent1", expLen)
					d.Decode(e, "Exponent1")
					outExp[index][i] = e.Exponent1
					Exp[index][i] = e.Exponent1
				} else {
					outExp[index][i] = Exp[ref_dps_id0][i]
					Exp[index][i] = Exp[ref_dps_id0][i]
				}
				d.Decode(e, "MantissaDiff")
				if predDirection == 0 {
					mantissaPred := ((Mantissa[ref_dps_id0][i]*uint64(predWeight0) + Mantissa[ref_dps_id1][i]*(64-uint64(predWeight0)) + 32) >> 6)
					outMantissa[index][i] = uint64(int64(mantissaPred) + e.MantissaDiff)
					Mantissa[index][i] = uint64(int64(mantissaPred) + e.MantissaDiff)
				} else {
					mantissaPred := Mantissa[ref_dps_id0][i]
					outMantissa[index][i] = uint64(int64(mantissaPred) + e.MantissaDiff)
					Mantissa[index][i] = uint64(int64(mantissaPred) + e.MantissaDiff)
				}
				outManLen[index][i] = ManLen[ref_dps_id0][i]
				ManLen[index][i] = ManLen[ref_dps_id0][i]
			} else {
				outSign[index][i] = Sign[ref_dps_id0][i]
				Sign[index][i] = Sign[ref_dps_id0][i]
				outExp[index][i] = Exp[ref_dps_id0][i]
				Exp[index][i] = Exp[ref_dps_id0][i]
				outMantissa[index][i] = Mantissa[ref_dps_id0][i]
				Mantissa[index][i] = Mantissa[ref_dps_id0][i]
				outManLen[index][i] = ManLen[ref_dps_id0][i]
				ManLen[index][i] = ManLen[ref_dps_id0][i]
			}
		}
	}
	if e.ElementEqualFlag {
		num_views_minus1, _ := d.Value("NumViewsMinus1").(uint64)
		// ????!??!?!? the string `deltaFlag` appears exactly once in the entire H.264 spec
		const deltaFlag = 0
		limit := int(num_views_minus1 + 1 - deltaFlag)
		for limit > numViews {
			outSign[index] = append(outSign[index], 0)
			outExp[index] = append(outExp[index], 0)
			outMantissa[index] = append(outMantissa[index], 0)
			outManLen[index] = append(outManLen[index], 0)
			numViews++
		}
		// Make sure the global values have enough space
		for len(Sign[index]) < numViews {
			Sign[index] = append(Sign[index], 0)
		}
		for len(Exp[index]) < numViews {
			Exp[index] = append(Exp[index], 0)
		}
		for len(Mantissa[index]) < numViews {
			Mantissa[index] = append(Mantissa[index], 0)
		}
		for len(ManLen[index]) < numViews {
			ManLen[index] = append(ManLen[index], 0)
		}
		for i := 1; i < limit; i++ {
			outSign[index][i] = outSign[index][0]
			Sign[index][i] = outSign[index][0]
			outExp[index][i] = outExp[index][0]
			Exp[index][i] = outExp[index][0]
			outMantissa[index][i] = outMantissa[index][0]
			Mantissa[index][i] = outMantissa[index][0]
			outManLen[index][i] = outManLen[index][0]
			ManLen[index][i] = outManLen[index][0]
		}
	}
	err = d.Error()

	// Store the modified global values back in the decoder
	d.SetValue(fmt.Sprintf("Z%sSign", distance), Sign)
	d.SetValue(fmt.Sprintf("Z%sExp", distance), Exp)
	d.SetValue(fmt.Sprintf("Z%sMantissa", distance), Mantissa)
	d.SetValue(fmt.Sprintf("Z%sManLen", distance), ManLen)

	return
}
