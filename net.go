package lib

import (
	"crypto/rand"
	"crypto/rsa"
	"net"
	"time"

	driver "github.com/vimcoders/go-driver"
)

const (
	DefaultBufferSize = 512
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

func (b *Buffer) Buffer() []byte {
	return b.buf
}

func NewBuffer() *Buffer {
	return NewBufferSize(DefaultBufferSize)
}

func NewBufferSize(n int) *Buffer {
	return &Buffer{
		buf: make([]byte, n),
	}
}

type Decoder struct {
	key     *rsa.PrivateKey
	message []byte
}

func (d *Decoder) ToBytes() (b []byte, err error) {
	return rsa.DecryptPKCS1v15(rand.Reader, d.key, d.message)
}

func NewDecoder(msg []byte, k *rsa.PrivateKey) driver.Message {
	return &Decoder{k, msg}
}

type Encoder struct {
	key     *rsa.PublicKey
	message []byte
}

func (e *Encoder) ToBytes() (b []byte, err error) {
	return rsa.EncryptPKCS1v15(rand.Reader, e.key, e.message)
}

func NewEncoder(msg []byte, k *rsa.PublicKey) driver.Message {
	return &Encoder{k, msg}
}

type Message struct {
	message []byte
}

func (m *Message) ToBytes() (b []byte, err error) {
	return m.message, nil
}

func NewMessage(msg []byte) driver.Message {
	return &Message{msg}
}

type Writer struct {
	net.Conn
	b       *Buffer
	timeout time.Duration
}

func (w *Writer) Write(pkg driver.Message) (err error) {
	b, err := pkg.ToBytes()

	if err != nil {
		return err
	}

	const header = 2
	length := len(b)

	buf := w.b.Take(length + header)

	copy(buf[header:], b)

	buf[0] = uint8(length >> 8)
	buf[1] = uint8(length)

	if err := w.SetDeadline(time.Now().Add(w.timeout)); err != nil {
		return err
	}

	if _, err := w.Conn.Write(buf); err != nil {
		return err
	}

	return nil
}

func NewWriter(c net.Conn, b *Buffer, t time.Duration) driver.Writer {
	return &Writer{c, b, t}
}

type Reader struct {
	net.Conn
	buffer  *Buffer
	timeout time.Duration
}

func (r *Reader) Read() (pkg driver.Message, err error) {
	if err := r.SetDeadline(time.Now().Add(r.timeout)); err != nil {
		return nil, err
	}

	buffer := r.buffer.Take(DefaultBufferSize)

	if _, err := r.Conn.Read(buffer); err != nil {
		return nil, err
	}

	length := int(uint32(buffer[0])<<8 | uint32(buffer[1]))

	const header = 2

	body := buffer[header : header+length]

	return &Message{body}, nil
}

func NewReader(c net.Conn, b *Buffer, t time.Duration) driver.Reader {
	return &Reader{c, b, t}
}
