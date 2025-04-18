package h264_test

import (
	"testing"

	"github.com/tachode/bitstream-go/bits"
	"github.com/tachode/bitstream-go/h264"
)

func TestNalUnitHeader3davcExtension_Read(t *testing.T) {
	decoder := bits.NewItuDecoder([]byte{0b10101010, 0b1_1_011_01_1})

	var header h264.NalUnitHeader3davcExtension
	err := header.Read(decoder)
	if err != nil {
		t.Fatalf("Failed to read NalUnitHeader3davcExtension: %v", err)
	}

	expected := h264.NalUnitHeader3davcExtension{
		ViewIdx:       170, // 0b10101010
		DepthFlag:     true,
		NonIdrFlag:    true,
		TemporalId:    3, // 0b011
		AnchorPicFlag: false,
		InterViewFlag: true,
	}

	if header != expected {
		t.Errorf("Decoded header does not match expected. Got %+v, want %+v", header, expected)
	}
}
