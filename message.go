package lib

import (
	"hash/crc32"

	"github.com/vimcoders/go-driver"
)

const (
	headerLength = 10
)

type Message struct {
	version  uint8
	protocol uint16
	payload  []byte
}

func (m *Message) Header() []byte {
	var header [headerLength]byte
	length := len(m.payload) + headerLength
	code := crc32.ChecksumIEEE(m.payload)

	header[0] = m.version
	header[1] = uint8(length >> 16)
	header[2] = uint8(length >> 8)
	header[3] = uint8(length)
	header[4] = uint8(m.protocol >> 8)
	header[5] = uint8(m.protocol)
	header[6] = uint8(code >> 24)
	header[7] = uint8(code >> 16)
	header[8] = uint8(code >> 8)
	header[9] = uint8(code)

	return header[:]
}

func (m *Message) Payload() []byte {
	return m.payload
}

func (m *Message) Protocol() uint16 {
	return m.protocol
}

func (m *Message) Version() uint8 {
	return m.version
}

func NewMessage(version uint8, protocol uint16, payload []byte) driver.Message {
	return &Message{version, protocol, payload}
}
