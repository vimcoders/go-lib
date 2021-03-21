package lib

const DEFAULT_BUFFER_SIZE = 128

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
	return NewBufferSize(DEFAULT_BUFFER_SIZE)
}

func NewBufferSize(n int) *Buffer {
	return &Buffer{
		buf: make([]byte, n),
	}
}
