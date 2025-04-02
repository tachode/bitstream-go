package h264

import "github.com/tachode/bitstream-go/bits"

type ScalingList struct {
	DeltaScale                     int64   `descriptor:"se(v)" json:"delta_scale"`
	ScalingList4x4                 [][]int `json:"scaling_list_4x4"`
	UseDefaultScalingMatrix4x4Flag []bool  `json:"use_default_scaling_matrix_4x4_flag"`
	ScalingList8x8                 [][]int `json:"scaling_list_8x8"`
	UseDefaultScalingMatrix8x8Flag []bool  `json:"use_default_scaling_matrix_8x8_flag"`
}

func (e *ScalingList) Read(d bits.Decoder, i int, sizeOfScalingList int) error {
	scalingList := make([]int, sizeOfScalingList)
	var useDefaultScalingMatrixFlag bool
	lastScale := 8
	nextScale := 8
	for j := 0; j < sizeOfScalingList; j++ {
		if nextScale != 0 {
			err := d.Decode(e, "DeltaScale")
			if err != nil {
				return err
			}
			nextScale := (lastScale + int(e.DeltaScale) + 256) % 256
			useDefaultScalingMatrixFlag = (j == 0 && nextScale == 0)
		}
		scale := nextScale
		if scale == 0 {
			scale = lastScale
		}
		scalingList[j] = scale
		lastScale = scalingList[j]
	}
	if sizeOfScalingList == 16 {
		if e.ScalingList4x4 == nil {
			e.ScalingList4x4 = make([][]int, 0, 6)
			e.UseDefaultScalingMatrix4x4Flag = make([]bool, 0, 6)
		}
		e.ScalingList4x4 = e.ScalingList4x4[:i+1]
		e.UseDefaultScalingMatrix4x4Flag = e.UseDefaultScalingMatrix4x4Flag[:i+1]
		e.ScalingList4x4[i] = scalingList
		e.UseDefaultScalingMatrix4x4Flag[i] = useDefaultScalingMatrixFlag
	} else {
		if e.ScalingList8x8 == nil {
			e.ScalingList8x8 = make([][]int, 0, 12)
			e.UseDefaultScalingMatrix8x8Flag = make([]bool, 0, 12)
		}
		e.ScalingList8x8 = e.ScalingList8x8[:i+1]
		e.UseDefaultScalingMatrix8x8Flag = e.UseDefaultScalingMatrix8x8Flag[:i+1]
		e.ScalingList8x8[i] = scalingList
		e.UseDefaultScalingMatrix8x8Flag[i] = useDefaultScalingMatrixFlag
	}
	return nil
}
