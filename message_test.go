package lib

import "testing"

func TestMessageHeader(t *testing.T) {
	message := NewMessage(1, 1, []byte{1, 2, 3})

	t.Log(message.Header())
}
