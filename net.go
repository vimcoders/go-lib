package lib

import (
	"bufio"
	"context"
	"errors"
	"net"
	"time"

	driver "github.com/vimcoders/go-driver"
	"google.golang.org/protobuf/proto"
)

const (
	Version           = 1
	HeaderLength      = 4
	DefaultBufferSize = 128
)

type Buffer struct {
	buf []byte
}

func (b *Buffer) Take(n int) []byte {
	if n < len(b.buf) {
		return b.buf[:n]
	}

	return make([]byte, n)
}

func NewBuffer() *Buffer {
	return NewBufferSize(DefaultBufferSize)
}

func NewBufferSize(n int) *Buffer {
	return &Buffer{
		buf: make([]byte, n),
	}
}

type Encoder struct {
	proto.Message
}

func (e *Encoder) ToBytes() (b []byte, errr error) {
	return proto.Marshal(e.Message)
}

func NewEncoder(msg proto.Message) driver.Message {
	return &Encoder{msg}
}

type Decoder struct {
	b []byte
}

func (d *Decoder) ToBytes() (b []byte, err error) {
	return d.b, nil
}

func NewDecoder(b []byte) driver.Message {
	return &Decoder{b}
}

type Conn struct {
	net.Conn
	OnMessage        func(pkg driver.Message) (err error)
	OnClose          func(e interface{})
	PushMessageQuene chan driver.Message
}

func (c *Conn) Write(pkg driver.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			c.OnClose(err)
		}

		c.OnClose(err)
	}()

	c.PushMessageQuene <- pkg

	return nil
}

func (c *Conn) Push(ctx context.Context) (err error) {
	defer func() {
		if err := recover(); err != nil {
			c.OnClose(err)
			close(c.PushMessageQuene)
			return
		}

		c.OnClose(err)
		close(c.PushMessageQuene)
	}()

	buffer := NewBuffer()

	for {
		select {
		case <-ctx.Done():
		default:
		}

		pkg, ok := <-c.PushMessageQuene

		if !ok {
			return errors.New("shutdown")
		}

		b, err := pkg.ToBytes()

		if err != nil {
			return err
		}

		header := make([]byte, HeaderLength)
		length := len(b)

		header[0] = Version
		header[1] = uint8(length >> 16)
		header[2] = uint8(length >> 8)
		header[3] = uint8(length)

		buf := buffer.Take(len(header) + len(b))
		copy(buf, header)
		copy(buf[len(header):], b)

		if err := c.SetWriteDeadline(time.Now().Add(time.Second * 5)); err != nil {
			return err
		}

		if _, err := c.Conn.Write(buf); err != nil {
			return err
		}
	}
}

func (c *Conn) Pull(ctx context.Context) (err error) {
	defer func() {
		if err := recover(); err != nil {
			c.OnClose(err)
		}
	}()

	reader := bufio.NewReaderSize(c.Conn, DefaultBufferSize)

	for {
		select {
		case <-ctx.Done():
		default:
		}

		header, err := reader.Peek(HeaderLength)

		if err != nil {
			return err
		}

		if header[0] != Version {
			return errors.New("Version is Unknown")
		}

		length := int(uint32(uint32(header[1])<<16 | uint32(header[2])<<8 | uint32(header[3])))

		buf, err := reader.Peek(length + len(header))

		if err != nil {
			return err
		}

		if err := c.OnMessage(NewDecoder(buf[len(header):])); err != nil {
			return err
		}

		if _, err := reader.Discard(len(buf)); err != nil {
			return err
		}
	}
}
