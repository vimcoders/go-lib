package lib

import (
	"context"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	driver "github.com/vimcoders/go-driver"
	"golang.org/x/net/websocket"
)

func TestWebsocket(t *testing.T) {
	go func() {
		http.Handle("/", websocket.Handler(func(c *websocket.Conn) {
			s := &Session{Conn: c, pushMessageQuene: make(chan driver.Message)}

			s.OnMessage = func(message driver.Message) (err error) {
				s.Write(message)
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

				return nil
			}

			go s.Pull(context.Background())
			s.Push(context.Background())
		}))

		if err := http.ListenAndServe(":8888", nil); err != nil {
			log.Fatal("ListenAndServe:", err)
		}
	}()

	var waitGroup sync.WaitGroup

	for i := 0; i < 10240; i++ {
		waitGroup.Add(1)

		go func() {
			ws, err := websocket.Dial("ws://localhost:8888", "", "http://localhost:8888/")

			if err != nil {
				t.Error(err)
				return
			}

			sess := NewSession(context.Background(), ws)

			for {
				sess.Write(NewMessage([]byte("hello")))

				time.Sleep(time.Second)
			}
		}()
	}

	waitGroup.Wait()
}

func TestTcp(t *testing.T) {
	var waitGroup sync.WaitGroup

	go func() {
		listener, err := net.Listen("tcp", ":8888")

		if err != nil {
			t.Error(err)
			return
		}

		for {
			conn, err := listener.Accept()

			if err != nil {
				t.Error(err)
				return
			}

			NewSession(context.Background(), conn)
		}
	}()

	time.Sleep(time.Second)

	for i := 0; i < 10240; i++ {
		waitGroup.Add(1)

		go func() {
			conn, err := net.Dial("tcp", ":8888")

			if err != nil {
				t.Error(err)
				return
			}

			sess := NewSession(context.Background(), conn)

			for {
				sess.Write(NewMessage([]byte("hello")))

				time.Sleep(time.Second)
			}
		}()
	}

	waitGroup.Wait()
}
