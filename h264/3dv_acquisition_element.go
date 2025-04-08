package h264

import "github.com/tachode/bitstream-go/bits"

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
) (
	outSign [][]uint8,
	outExp [][]uint64,
	outMantissa [][]uint64,
	outManLen [][]uint8,
	err error,
) {
	/*
		ref_dps_id0, _ := d.Value("RefDpsId0").(uint64)
		ref_dps_id1, _ := d.Value("RefDpsId1").(uint64)
		predWeight0, _ := d.Value("PredWeight0").(uint8)
	*/

	outSign = make([][]uint8, index+1)
	outSign[index] = make([]uint8, numViews)
	outExp = make([][]uint64, index+1)
	outExp[index] = make([]uint64, numViews)
	outMantissa = make([][]uint64, index+1)
	outMantissa[index] = make([]uint64, numViews)
	outManLen = make([][]uint8, index+1)
	outManLen[index] = make([]uint8, numViews)

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
			manLen = int(e.MantissaLenMinus1 + 1)
		}
		if predDirection == 2 {
			d.Decode(e, "Sign0")
			outSign[index][i] = e.Sign0
			d.SetValueLength("Exponent0", expLen)
			d.Decode(e, "Exponent0")
			outExp[index][i] = e.Exponent0
			d.SetValueLength("Mantissa0", manLen)
			d.Decode(e, "Mantissa0")
			outMantissa[index][i] = e.Mantissa0
		} else {
			d.Decode(e, "SkipFlag")
			if !e.SkipFlag {
				d.Decode(e, "Sign1")
				outSign[index][i] = e.Sign1
				d.Decode(e, "ExponentSkipFlag")
				if !e.ExponentSkipFlag {
					d.SetValueLength("Exponent1", expLen)
					d.Decode(e, "Exponent1")
					outExp[index][i] = e.Exponent1
				} else {
					// TODO -- read this value from the decoder
					// outExp[index][i] = outExp[ref_dps_id0][i]
				}
				d.Decode(e, "MantissaDiff")
				if predDirection == 0 {
					// TODO -- read these values from the decoder
					/*
						mantissaPred := ((outMantissa[ref_dps_id0][i]*uint64(predWeight0) + outMantissa[ref_dps_id1][i]*(64-uint64(predWeight0)) + 32) >> 6)
						outMantissa[index][i] = uint64(int64(mantissaPred) + e.MantissaDiff)
					*/
				} else {
					// TODO -- read these values from the decoder
					/*
						mantissaPred := outMantissa[ref_dps_id0][i]
						outMantissa[index][i] = uint64(int64(mantissaPred) + e.MantissaDiff)
					*/
				}
				// TODO -- read this value from the decoder
				// outManLen[index][i] = outManLen[ref_dps_id0][i]
			} else {
				// TODO -- read these values from the decoder
				/*
					outSign[index][i] = outSign[ref_dps_id0][i]
					outExp[index][i] = outExp[ref_dps_id0][i]
					outMantissa[index][i] = outMantissa[ref_dps_id0][i]
					outManLen[index][i] = outManLen[ref_dps_id0][i]
				*/
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
		for i := 1; i < limit; i++ {
			outSign[index][i] = outSign[index][0]
			outExp[index][i] = outExp[index][0]
			outMantissa[index][i] = outMantissa[index][0]
			outManLen[index][i] = outManLen[index][0]
		}
	}
	err = d.Error()
	return
}
