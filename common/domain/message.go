package domain

import (
	"encoding/json"
)

type Message struct {
	FHeader *FixedHeader `json:"f_header"`
	VHeader []byte       `json:"v_header"`
	Payload []byte       `json:"payload"`
}

type FixedHeader struct {
	V          byte   `json:"v"`
	Cmd        byte   `json:"cmd"`
	VarHLen    uint32 `json:"var_header_len"`
	PayloadLen uint32 `json:"payload_len"`
	Crc32sum   uint32 `json:"crc_32_sum"`
}

func (m *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		FHeader *FixedHeader `json:"f_header"`
		VHeader string       `json:"v_header"`
		Payload string       `json:"payload"`
	}{
		FHeader: m.FHeader,
		VHeader: string(m.VHeader),
		Payload: string(m.Payload),
	})
}
