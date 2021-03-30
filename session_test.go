package lib

import (
	"context"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/websocket"
)

func TestMain(m *testing.M) {
	http.Handle("/", websocket.Handler(func(c *websocket.Conn) {
		sess := NewSession(context.Background(), c)

		for {
			sess.Write(NewMessage([]byte("response")))
		}
	}))

	go func() {
		if err := http.ListenAndServe(":8888", nil); err != nil {
			log.Fatal("ListenAndServe:", err)
		}
	}()

	time.Sleep(time.Second)

	m.Run()
}

func TestWebsocket(t *testing.T) {
	var waitGroup sync.WaitGroup

	for i := 0; i < 1024; i++ {
		waitGroup.Add(1)

		go func() {
			ws, err := websocket.Dial("ws://localhost:8888", "", "http://localhost:8888/")

			if err != nil {
				t.Error(err)
				return
			}

			//defer ws.Close() //关闭连接

			sess := NewSession(context.Background(), ws)

			for {
				sess.Write(NewMessage([]byte("hello")))
			}
		}()
	}

	waitGroup.Wait()
}
