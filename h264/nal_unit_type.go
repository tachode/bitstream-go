package h264

import "encoding/json"

type NalUnitType uint8

//go:generate stringer -type=NalUnitType -trimprefix=NalUnitType
const (
	NalUnitTypeCodedSliceNonIdr         NalUnitType = 1  // Coded slice of a non-IDR picture
	NalUnitTypeCodedSliceDataPartitionA NalUnitType = 2  // Coded slice data partition A
	NalUnitTypeCodedSliceDataPartitionB NalUnitType = 3  // Coded slice data partition B
	NalUnitTypeCodedSliceDataPartitionC NalUnitType = 4  // Coded slice data partition C
	NalUnitTypeCodedSliceIdr            NalUnitType = 5  // Coded slice of an IDR picture
	NalUnitTypeSEI                      NalUnitType = 6  // Supplemental enhancement information (SEI)
	NalUnitTypeSPS                      NalUnitType = 7  // Sequence parameter set
	NalUnitTypePPS                      NalUnitType = 8  // Picture parameter set
	NalUnitTypeAUD                      NalUnitType = 9  // Access unit delimiter
	NalUnitTypeEndOfSequence            NalUnitType = 10 // End of sequence
	NalUnitTypeEndOfStream              NalUnitType = 11 // End of stream
	NalUnitTypeFiller                   NalUnitType = 12 // Filler data
	NalUnitTypeSpsExt                   NalUnitType = 13 // Sequence parameter set extension
	NalUnitTypePrefixNalUnit            NalUnitType = 14 // Prefix NAL unit
	NalUnitTypeSubsetSeqParameterSet    NalUnitType = 15 // Subset sequence parameter set
	NalUnitTypeDepthParameterSet        NalUnitType = 16 // Depth parameter set
	NalUnitTypeCodedSliceAux            NalUnitType = 19 // Coded slice of an auxiliary coded picture without partitioning
	NalUnitTypeCodedSliceExtension      NalUnitType = 20 // Coded slice extension
	NalUnitTypeCodedSliceExtension3D    NalUnitType = 21 // Coded slice extension for a depth view component or a 3D-AVC texture view component
)

func (t NalUnitType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *NalUnitType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	for i := NalUnitTypeCodedSliceNonIdr; i <= NalUnitTypeCodedSliceExtension3D; i++ {
		if i.String() == s {
			*t = i
			return nil
		}
	}

	return json.Unmarshal(data, (*uint8)(t))
}
