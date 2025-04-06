package h264

import (
	"fmt"
	"reflect"

	"github.com/tachode/bitstream-go/bits"
)

type NalPayload interface {
	Read(d bits.Decoder) error
}

var nalPayloadRegistry map[NalUnitType]NalPayload

func RegisterNalPayloadType(typ NalUnitType, v NalPayload) {
	if nalPayloadRegistry == nil {
		nalPayloadRegistry = make(map[NalUnitType]NalPayload)
	}
	nalPayloadRegistry[typ] = v
}

type Parser struct {
	decoder bits.Decoder
}

func NewParser() *Parser {
	return &Parser{decoder: bits.NewItuDecoder([]byte{})}
}

func (p *Parser) Parse(buffer []byte) (*NalUnit, error) {
	p.decoder.Reset(buffer)
	nal := &NalUnit{}
	err := nal.Read(p.decoder, len(buffer))
	if err != nil {
		return nil, err
	}

	// Create a slice with the trailing zeros removed
	nalUnitBuffer := nal.Rbsp
	for len(nalUnitBuffer) > 0 && nalUnitBuffer[len(nalUnitBuffer)-1] == 0 {
		nalUnitBuffer = nalUnitBuffer[:len(nalUnitBuffer)-1]
	}

	p.decoder.Reset(nalUnitBuffer)

	prototype, ok := nalPayloadRegistry[nal.NalUnitType]
	if !ok {
		// We don't understand this NAL type, but its payload is in the RBSP bytes,
		// so the application may be able to make use of it anyway.
		return nal, p.decoder.Error()
	}
	copy := reflect.New(reflect.Indirect(reflect.ValueOf(prototype)).Type()).Interface()
	payload, ok := copy.(NalPayload)
	if !ok {
		return nil, fmt.Errorf("invalid registered type %T does not implement NalPayload interface", prototype)
	}
	payload.Read(p.decoder)
	nal.Payload = payload
	return nal, p.decoder.Error()
}
