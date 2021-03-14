package lib

import (
	"context"
	"net"

	driver "github.com/vimcoders/go-driver"
)

type Session struct {
	id int64
	driver.WriteCloser
	v map[interface{}]interface{}
}

func (s *Session) SessionID() int64 {
	return s.id
}

func (s *Session) Set(key, value interface{}) error {
	s.v[key] = value
	return nil
}

func (s *Session) Get(key interface{}) interface{} {
	return s.v[key]
}

func (s *Session) Delete(key interface{}) error {
	delete(s.v, key)
	return nil
}

func (s *Session) Close() error {
	return nil
}

func NewSession(ctx context.Context, c net.Conn) driver.Session {
	var s Session

	conn := &Conn{
		Conn: c,
		OnMessage: func(pkg driver.Message) (err error) {
			return nil
		},
		OnClose: func(e interface{}) {
			s.Close()
		},
		PushMessageQuene: make(chan driver.Message),
	}

	go conn.Pull(ctx)
	go conn.Push(ctx)

	s.WriteCloser = conn

	return &s
}
