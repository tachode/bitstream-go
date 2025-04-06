package h264

import "github.com/tachode/bitstream-go/bits"

type ReservedSeiMessage struct {
	ReservedSeiMessagePayload []byte `json:"reserved_sei_message_payload"`
}

func (e *ReservedSeiMessage) Read(d bits.Decoder, payloadSize int) error {
	e.ReservedSeiMessagePayload = make([]byte, payloadSize)
	_, err := d.Read(e.ReservedSeiMessagePayload)
	return err
}
