package lib

import (
	"bufio"
	"errors"
	"hash/crc32"

	"github.com/vimcoders/go-driver"
)

const (
	HEADER_LENGTH = 10
)

type Message struct {
	version  uint8
	protocol uint16
	payload  []byte
}

func (m *Message) Header() []byte {
	var header [HEADER_LENGTH]byte
	length := len(m.payload) + HEADER_LENGTH
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

func (m *Message) Length() int32 {
	return int32(len(m.payload)) + HEADER_LENGTH
}

func NewMessage(reader *bufio.Reader) (driver.Message, error) {
	header, err := reader.Peek(HEADER_LENGTH)

	if err != nil {
		return nil, err
	}

	length := uint32(uint32(header[1])<<16 | uint32(header[2])<<8 | uint32(header[3]))

	buf, err := reader.Peek(int(length))

	if len(buf) < HEADER_LENGTH {
		return nil, nil
	}

	version := buf[0]
	protocol := uint16(buf[4])<<8 | uint16(buf[5])
	code := uint32(buf[6])<<24 | uint32(buf[7])<<16 | uint32(buf[8])<<8 | uint32(buf[9])

	if code != crc32.ChecksumIEEE(buf[HEADER_LENGTH:]) {
		return nil, errors.New("code !")
	}

	return &Message{version, protocol, buf[HEADER_LENGTH:]}, nil
}
