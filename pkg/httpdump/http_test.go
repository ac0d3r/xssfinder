package httpdump

import (
	"net/http"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	req := Request{
		Method: http.MethodGet,
		URL:    "https://www.baidu.com",
	}

	t.Log(Do(req, time.Second))

	t.Log(req)
}
