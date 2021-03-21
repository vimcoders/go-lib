package lib

import (
	"bufio"
	"errors"
	"fmt"
	"hash/crc32"
	"io"

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

func NewMessage(version uint8, protocol uint16, payload []byte) driver.Message {
	return &Message{version, protocol, payload}
}

type Reader struct {
	reader *bufio.Reader
}

func (r *Reader) Read() (driver.Message, error) {
	header, err := r.reader.Peek(HEADER_LENGTH)

	if err != nil {
		return nil, err
	}

	version := header[0]
	length := uint32(uint32(header[1])<<16 | uint32(header[2])<<8 | uint32(header[3]))
	protocol := uint16(header[4])<<8 | uint16(header[5])
	code := uint32(header[6])<<24 | uint32(header[7])<<16 | uint32(header[8])<<8 | uint32(header[9])

	buf, err := r.reader.Peek(int(length))

	if err != nil {
		return nil, err
	}

	if len(buf) < HEADER_LENGTH {
		return nil, err
	}

	if code != crc32.ChecksumIEEE(buf[HEADER_LENGTH:]) {
		return nil, errors.New(fmt.Sprintf("protocol %v code uncomplete", protocol))
	}

	return NewMessage(version, protocol, buf[HEADER_LENGTH:]), nil
}

func (r *Reader) Discard(msg driver.Message) error {
	if _, err := r.reader.Discard(len(msg.Payload()) + HEADER_LENGTH); err != nil {
		return err
	}

	return nil
}

func NewReader(r io.Reader) *Reader {
	return &Reader{bufio.NewReader(r)}
}
