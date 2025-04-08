package h264

import "github.com/tachode/bitstream-go/bits"

type VspParam struct {
	DisparityDiffWji [][]uint64 `descriptor:"ue(v)" json:"disparity_diff_wji"`
	DisparityDiffOji [][]uint64 `descriptor:"ue(v)" json:"disparity_diff_oji"`
	DisparityDiffWij [][]uint64 `descriptor:"ue(v)" json:"disparity_diff_wij"`
	DisparityDiffOij [][]uint64 `descriptor:"ue(v)" json:"disparity_diff_oij"`
}

func (e *VspParam) Read(d bits.Decoder, numViews int, predDirection int, index int) error {
	for i := 0; i < numViews; i++ {
		for j := 0; j < i; j++ {
			d.DecodeIndex(e, "DisparityDiffWji", j, i)
			d.DecodeIndex(e, "DisparityDiffOji", j, i)
			d.DecodeIndex(e, "DisparityDiffWij", i, j)
			d.DecodeIndex(e, "DisparityDiffOij", i, j)
		}
	}
	return d.Error()
}
