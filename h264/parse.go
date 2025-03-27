package h264

import (
	"bytes"

	"github.com/tachode/bitstream-go/bits"
)

func Parse(buffer []byte) (*NalUnit, error) {
	reader := &bits.ReadBuffer{Reader: bytes.NewBuffer(buffer)}
	ituReader := &bits.ItuReader{Reader: reader}
	decoder := bits.NewItuDecoder(ituReader)
	nal := &NalUnit{}
	err := nal.Read(decoder, len(buffer))
	if err != nil {
		return nil, err
	}

	reader = &bits.ReadBuffer{Reader: bytes.NewBuffer(nal.RbspByte)}
	ituReader = &bits.ItuReader{Reader: reader}
	decoder = bits.NewItuDecoder(ituReader)

	switch nal.NalUnitType {
	case NalUnitTypeCodedSliceNonIdr:
		// payload := &SliceLayerWithoutPartioning{}
		// payload.Read(decoder)
		// nal.Payload = payload
	case NalUnitTypeCodedSliceDataPartitionA:
		// payload := &SliceDataPartitionALayer{}
		// payload.Read(decoder)
		// nal.Payload = payload
	case NalUnitTypeCodedSliceDataPartitionB:
		// payload := &SliceDataPartitionBLayer{}
		// payload.Read(decoder)
		// nal.Payload = payload
	case NalUnitTypeCodedSliceDataPartitionC:
		// payload := &SliceDataPartitionCLayer{}
		// payload.Read(decoder)
		// nal.Payload = payload
	case NalUnitTypeCodedSliceIdr:
		// payload := &SliceLayerWithoutPartioning{}
		// payload.Read(decoder)
		// nal.Payload = payload
	case NalUnitTypeSEI:
		// payload := &Sei{}
		// payload.Read(decoder)
		// nal.Payload = payload
	case NalUnitTypeSPS:
		payload := &SeqParameterSet{}
		payload.Read(decoder)
		nal.Payload = payload
	case NalUnitTypePPS:
		// payload := &PicParameterSet{}
		// payload.Read(decoder)
		// nal.Payload = payload
	case NalUnitTypeAUD:
		// payload := &AccessUnitDelimiter{}
		// payload.Read(decoder)
		// nal.Payload = payload
	case NalUnitTypeEndOfSequence:
		// payload := &EndOfSeq{}
		// payload.Read(decoder)
		// nal.Payload = payload
	case NalUnitTypeEndOfStream:
		// payload := &EndOfStream{}
		// payload.Read(decoder)
		// nal.Payload = payload
	case NalUnitTypeFiller:
		// payload := &FillerData{}
		// payload.Read(decoder)
		// nal.Payload = payload
	case NalUnitTypeSpsExt:
		// payload := &SeqParameterSetExtension{}
		// payload.Read(decoder)
		// nal.Payload = payload
	case NalUnitTypePrefixNalUnit:
		// payload := &PrefixNalUnit{}
		// payload.Read(decoder)
		// nal.Payload = payload
	case NalUnitTypeSubsetSeqParameterSet:
		// payload := &SubsetSeqParameterSet{}
		// payload.Read(decoder)
		// nal.Payload = payload
	case NalUnitTypeDepthParameterSet:
		// payload := &DepthParameterSet{}
		// payload.Read(decoder)
		// nal.Payload = payload
	case NalUnitTypeCodedSliceAux:
		// payload := &SliceLayerWithoutPartioning{}
		// payload.Read(decoder)
		// nal.Payload = payload
	case NalUnitTypeCodedSliceExtension:
		// payload := &SliceLayerExtension{}
		// payload.Read(decoder)
		// nal.Payload = payload
	case NalUnitTypeCodedSliceExtension3D:
		// payload := &SliceLayerExtension{}
		// payload.Read(decoder)
		// nal.Payload = payload
	}

	return nal, decoder.Error()
}
