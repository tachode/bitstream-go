package h264

import (
	"bufio"
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

	// Create a slice with the trailing zeros removed
	nalUnitBuffer := nal.RbspByte
	for len(nalUnitBuffer) > 0 && nalUnitBuffer[len(nalUnitBuffer)-1] == 0 {
		nalUnitBuffer = nalUnitBuffer[:len(nalUnitBuffer)-1]
	}

	// In order to implement several H.264 functions, parsing of RBSP payload
	// need to be buffered. The functions that require this are more_rbsp_data()
	// more_rbsp_trailing_data(), and next_bits().
	reader = &bits.ReadBuffer{Reader: bufio.NewReader(bytes.NewBuffer(nalUnitBuffer))}
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
