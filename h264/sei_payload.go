package h264

import (
	"fmt"
	"reflect"

	"github.com/tachode/bitstream-go/bits"
)

type SeiPayload struct {
	PayloadType    SeiType         `json:"payload_type"`
	Payload        SeiPayloadValue `json:"payload,omitempty"`
	BitEqualToOne  bool            `descriptor:"f(1)=1" json:"bit_equal_to_one"`
	BitEqualToZero bool            `descriptor:"f(1)=0" json:"bit_equal_to_zero"`
}

type SeiPayloadValue interface {
	Read(d bits.Decoder, payloadSize int) error
}

var seiPayloadTypes map[SeiType]SeiPayloadValue

func RegisterSeiPayloadType(payloadType SeiType, prototype SeiPayloadValue) {
	if seiPayloadTypes == nil {
		seiPayloadTypes = make(map[SeiType]SeiPayloadValue)
	}
	seiPayloadTypes[payloadType] = prototype
}

func (e *SeiPayload) Read(d bits.Decoder, payloadType int, payloadSize int) error {
	e.PayloadType = SeiType(payloadType)
	prototype, ok := seiPayloadTypes[SeiType(payloadType)]
	if ok {
		copy := reflect.New(reflect.Indirect(reflect.ValueOf(prototype)).Type()).Interface()
		e.Payload, ok = copy.(SeiPayloadValue)
		if !ok {
			return fmt.Errorf("invalid registered type %T does not implement SeiPayloadValue interface", prototype)
		}
	} else {
		e.Payload = &ReservedSeiMessage{}
	}

	e.Payload.Read(d, payloadSize)

	if !d.ByteAligned() {
		d.Decode(e, "BitEqualToOne")
		for !d.ByteAligned() {
			err := d.Decode(e, "BitEqualToZero")
			if err != nil {
				return err
			}
		}
	}
	return d.Error()
}
