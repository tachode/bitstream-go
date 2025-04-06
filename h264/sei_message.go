package h264

import (
	"github.com/tachode/bitstream-go/bits"
)

type SeiMessage struct {
	FfByte              uint8       `descriptor:"f(8)=255" json:"ff_byte"`
	LastPayloadTypeByte uint8       `descriptor:"u(8)" json:"last_payload_type_byte"`
	LastPayloadSizeByte uint8       `descriptor:"u(8)" json:"last_payload_size_byte"`
	SeiPayload          *SeiPayload `json:"sei_payload,omitempty"`
}

func (e *SeiMessage) Read(d bits.Decoder) (err error) {
	payloadType := 0
	for d.NextBits(8) == 0xFF {
		err = d.Decode(e, "FfByte")
		if err != nil {
			return
		}
		payloadType += 255
	}
	d.Decode(e, "LastPayloadTypeByte")
	payloadType += int(e.LastPayloadTypeByte)
	payloadSize := 0
	for d.NextBits(8) == 0xFF {
		err = d.Decode(e, "FfByte")
		if err != nil {
			return
		}
		payloadSize += 255
	}
	d.Decode(e, "LastPayloadSizeByte")
	payloadSize += int(e.LastPayloadSizeByte)
	e.SeiPayload = &SeiPayload{}
	e.SeiPayload.Read(d, payloadType, payloadSize)
	return d.Error()
}
