package proxy

import (
	"testing"
)

func TestMitmProxy(t *testing.T) {
	mitm := NewMitmServer(":8080", false, "hyuga.icu")
	mitm.Run()
}
