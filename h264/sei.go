package h264

import (
	"github.com/tachode/bitstream-go/bits"
)

func init() { RegisterNalPayloadType(NalUnitTypeSEI, &Sei{}) }

type Sei struct {
	SeiMessage []*SeiMessage `json:"sei_message,omitempty"`
}

func (e *Sei) Read(d bits.Decoder) error {
	for d.MoreRbspData() {
		seiMessage := &SeiMessage{}
		seiMessage.Read(d)
		e.SeiMessage = append(e.SeiMessage, seiMessage)
	}
	return d.Error()
}
