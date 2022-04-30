package dingbot

import (
	"os"
	"testing"
)

func TestDingbot(t *testing.T) {
	d := New(os.Getenv("dingbot_token"), os.Getenv("dingbot_secret"))
	d.Notify("http://localhost:8080/", `- type dom`)
}
