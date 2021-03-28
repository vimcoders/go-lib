package lib

import (
	"github.com/vimcoders/go-driver"
)

const (
	HeaderLength = 4
)

type Message struct {
	payload []byte
}

func (m *Message) Header() []byte {
	header := make([]byte, HeaderLength)
	length := len(m.payload) + len(header)

	header[1] = uint8(length >> 16)
	header[2] = uint8(length >> 8)
	header[3] = uint8(length)

	return header[:]
}

func (m *Message) Payload() []byte {
	return m.payload
}

func NewMessage(payload []byte) driver.Message {
	return &Message{payload}
}
