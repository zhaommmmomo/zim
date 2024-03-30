package gateway

import (
	"encoding/binary"
	"fmt"
	"github.com/zhaommmmomo/zim/common/domain"
	"io"
)

func encoder(m *domain.Message) []byte {
	buf := make([]byte, 14+m.FHeader.VarHLen+m.FHeader.PayloadLen)
	buf[0] = m.FHeader.V
	buf[1] = m.FHeader.Cmd
	binary.BigEndian.PutUint32(buf[2:6], m.FHeader.VarHLen)
	binary.BigEndian.PutUint32(buf[6:10], m.FHeader.PayloadLen)
	binary.BigEndian.PutUint32(buf[10:14], m.FHeader.Crc32sum)
	copy(buf[14:], m.VHeader)
	copy(buf[14+m.FHeader.VarHLen:], m.Payload)
	return buf
}

func decoder(r io.Reader) (*domain.Message, error) {
	// 处理 fixed header
	fHeaderBuf := make([]byte, 14)
	if _, err := r.Read(fHeaderBuf); err != nil {
		return nil, fmt.Errorf("read fixed header fail: %w", err)
	}
	// 处理 var header
	vHeaderLen := binary.BigEndian.Uint32(fHeaderBuf[2:6])
	var vHeader []byte
	if vHeaderLen > 0 {
		vHeader = make([]byte, vHeaderLen)
		if _, err := r.Read(vHeader); err != nil {
			return nil, fmt.Errorf("read var header fail: %w", err)
		}
	}
	// 处理 payload
	payloadLen := binary.BigEndian.Uint32(fHeaderBuf[6:10])
	var payload []byte
	if payloadLen > 0 {
		payload = make([]byte, payloadLen)
		if _, err := r.Read(payload); err != nil {
			return nil, fmt.Errorf("read payload fail: %w", err)
		}
	}
	return &domain.Message{
		FHeader: &domain.FixedHeader{
			V:          fHeaderBuf[0],
			Cmd:        fHeaderBuf[1],
			VarHLen:    vHeaderLen,
			PayloadLen: payloadLen,
			Crc32sum:   binary.BigEndian.Uint32(fHeaderBuf[10:14]),
		},
		VHeader: vHeader,
		Payload: payload,
	}, nil
}
