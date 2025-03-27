package h264

import "github.com/tachode/bitstream-go/bits"

type ScalingList struct {
	DeltaScale int64 `descriptor:"se(v)" json:"delta_scale"`
}

func (e *ScalingList) Read(d bits.Decoder) (scalingList []int, sizeOfScalingList int, useDefaultScalingMatrixFlag bool, err error) {
	scalingList = make([]int, sizeOfScalingList)
	lastScale := 8
	nextScale := 8
	for j := 0; j < sizeOfScalingList; j++ {
		if nextScale != 0 {
			d.Decode(e, "DeltaScale")
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
	err = d.Error()
	return
}
