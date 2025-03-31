package h264

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/tachode/bitstream-go/bits"
)

type NalPayload interface {
	Read(d bits.Decoder) error
}

var nalPayloadRegistry map[NalUnitType]NalPayload

func RegisterNalPayload(typ NalUnitType, v NalPayload) {
	if nalPayloadRegistry == nil {
		nalPayloadRegistry = make(map[NalUnitType]NalPayload)
	}
	nalPayloadRegistry[typ] = v
}

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

	prototype, ok := nalPayloadRegistry[nal.NalUnitType]
	if !ok {
		// We don't understand this NAL type, but its payload is in the RBSP bytes,
		// so the application may be able to make use of it anyway.
		return nal, decoder.Error()
	}
	copy := reflect.New(reflect.Indirect(reflect.ValueOf(prototype)).Type()).Interface()
	payload, ok := copy.(NalPayload)
	if !ok {
		return nil, fmt.Errorf("invalid registered type %T does not implement NalPayload interface", prototype)
	}
	payload.Read(decoder)
	nal.Payload = payload
	return nal, decoder.Error()
}
