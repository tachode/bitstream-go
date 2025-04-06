package h264

import (
	"io"

	"github.com/google/uuid"
	"github.com/tachode/bitstream-go/bits"
)

func init() { RegisterSeiPayloadType(SeiTypeUserDataUnregistered, &UserDataUnregistered{}) }

type UserDataUnregistered struct {
	UuidIsoIec11578 uuid.UUID `json:"uuid_iso_iec_11578"`
	UserDataPayload []byte    `json:"user_data_payload_byte"`
}

func (e *UserDataUnregistered) Read(d bits.Decoder, payloadSize int) error {
	id := make([]byte, 16)
	n, err := d.Read(id)
	if err != nil {
		return err
	}
	if n != 16 {
		return io.ErrUnexpectedEOF
	}
	uid, err := uuid.FromBytes(id)
	if err != nil {
		return err
	}
	e.UuidIsoIec11578 = uid
	e.UserDataPayload = make([]byte, payloadSize-16)
	_, err = d.Read(e.UserDataPayload)
	return err
}
