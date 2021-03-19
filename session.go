package lib

import (
	"bufio"
	"context"
	"errors"
	"net"
	"reflect"
	"strings"
	"time"

	driver "github.com/vimcoders/go-driver"
)

type Session struct {
	UserID int64

	net.Conn

	OnMessage func(message driver.Message) (err error)

	pushMessageQune chan driver.Message
}

func (s *Session) SessionID() int64 {
	return 0
}

func (s *Session) Set(key, v interface{}) error {
	return nil
}

func (s *Session) Delete(key interface{}) error {
	return nil
}

func (s *Session) Get(key interface{}) interface{} {
	return nil
}

func (s *Session) Send(msg driver.Message) (err error) {
	return nil
}

func (s *Session) PullMessage(ctx context.Context) (err error) {
	reader := bufio.NewReader(s.Conn)

	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		default:
		}

		pkg, err := ReadMessage(reader)

		if err != nil {
			return err
		}

		if s.OnMessage == nil {
			continue
		}

		if err := s.OnMessage(pkg); err != nil {
			return err
		}

		if _, err := reader.Discard(int(pkg.Length())); err != nil {
			return err
		}
	}
}

func (s *Session) PushMessage(ctx context.Context) (err error) {
	defaultBufferSize := 1024

	buffer := NewBuffer(defaultBufferSize)

	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		default:
			pkg, ok := <-s.pushMessageQune

			if !ok || pkg == nil {
				return errors.New("shut down!")
			}

			header, payload := pkg.Header(), pkg.Payload()

			buf := buffer.Take(len(header) + len(payload))

			copy(buf, header)
			copy(buf[len(header):], payload)

			if err := s.Conn.SetWriteDeadline(time.Now().Add(time.Second)); err != nil {
				return err
			}

			if _, err := s.Conn.Write(buf); err != nil {
				return err
			}
		}
	}
}

func NewSession(ctx context.Context, c net.Conn) (session driver.Session) {
	s := &Session{
		Conn:            c,
		pushMessageQune: make(chan driver.Message),
	}

	s.OnMessage = func(message driver.Message) (err error) {
		var methodName string

		t, v := reflect.TypeOf(s), reflect.ValueOf(s)

		for i := 0; i < t.NumMethod(); i++ {
			if strings.ToLower(t.Method(i).Name) != methodName {
				continue
			}

			arg1 := reflect.ValueOf(context.Background())
			arg2 := reflect.ValueOf(message.Payload())

			v.Method(i).Call([]reflect.Value{arg1, arg2})

			return nil
		}

		return errors.New("unknow")
	}

	go s.PullMessage(ctx)
	go s.PushMessage(ctx)

	return s
}

type buffer struct {
	buf []byte
}

func (b *buffer) Take(n int) []byte {
	if n < len(b.buf) {
		return b.buf[:n]
	}

	return make([]byte, n)
}

func NewBuffer(n int) *buffer {
	return &buffer{
		buf: make([]byte, n),
	}
}
