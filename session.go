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
	OnMessage        func(pkg driver.Message) (err error)
	pushMessageQuene chan driver.Message
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

func (s *Session) Write(pkg driver.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			//TODO:log
		}
	}()

	s.pushMessageQuene <- pkg

	return nil
}

func (s *Session) Push(ctx context.Context) (err error) {
	defer func() {
		if err := recover(); err != nil {
		}

		close(s.pushMessageQuene)
	}()

	buffer := NewBuffer()

	for {
		select {
		case <-ctx.Done():
		default:
		}

		pkg, ok := <-s.pushMessageQuene

		if !ok {
			return errors.New("shutdown")
		}

		if err := s.SetWriteDeadline(time.Now().Add(time.Second)); err != nil {
			return err
		}

		header, payload := pkg.Header(), pkg.Payload()

		buf := buffer.Take(len(header) + len(payload))
		copy(buf, header)
		copy(buf[HEADER_LENGTH:], payload)

		if _, err := s.Conn.Write(buf); err != nil {
			return err
		}
	}
}

func (s *Session) Pull(ctx context.Context) (err error) {
	defer func() {
		if err := recover(); err != nil {
			//TODO::
		}
	}()

	reader := bufio.NewReader(s.Conn)

	for {
		select {
		case <-ctx.Done():
		default:
		}

		pkg, err := NewMessage(reader)

		if err != nil {
			return err
		}

		if err := s.OnMessage(pkg); err != nil {
			return err
		}

		if _, err := reader.Discard(int(pkg.Length())); err != nil {
			return err
		}
	}
}

func NewSession(ctx context.Context, c net.Conn) (session driver.Session) {
	s := &Session{Conn: c}

	s.OnMessage = func(message driver.Message) (err error) {
		//TODO::config
		var methodName string

		t, _ := reflect.TypeOf(s), reflect.ValueOf(s)
		//t, v := reflect.TypeOf(s), reflect.ValueOf(s)

		for i := 0; i < t.NumMethod(); i++ {
			if strings.ToLower(t.Method(i).Name) != methodName {
				continue
			}

			//TODO::
			//v.Method(i).Call([]reflect.Value{arg1, arg2})

			return nil
		}

		return errors.New("unknow")
	}

	go s.Pull(ctx)
	go s.Push(ctx)

	return s
}
