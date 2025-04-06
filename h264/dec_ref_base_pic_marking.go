package h264

import "github.com/tachode/bitstream-go/bits"

type DecRefBasePicMarking struct {
	AdaptiveRefBasePicMarkingModeFlag    bool   `descriptor:"u(1)" json:"adaptive_ref_base_pic_marking_mode_flag"`
	MemoryManagementBaseControlOperation uint64 `descriptor:"ue(v)" json:"memory_management_base_control_operation"`
	DifferenceOfBasePicNumsMinus1        uint64 `descriptor:"ue(v)" json:"difference_of_base_pic_nums_minus1"`
	LongTermBasePicNum                   uint64 `descriptor:"ue(v)" json:"long_term_base_pic_num"`
}

func (e *DecRefBasePicMarking) Read(d bits.Decoder) error {
	d.Decode(e, "AdaptiveRefBasePicMarkingModeFlag")
	if e.AdaptiveRefBasePicMarkingModeFlag {
		for {
			d.Decode(e, "MemoryManagementBaseControlOperation")
			if e.MemoryManagementBaseControlOperation == 1 {
				d.Decode(e, "DifferenceOfBasePicNumsMinus1")
			}
			if e.MemoryManagementBaseControlOperation == 2 {
				d.Decode(e, "LongTermBasePicNum")
			}
			if e.MemoryManagementBaseControlOperation == 0 {
				break
			}
		}
	}
	return d.Error()
}
