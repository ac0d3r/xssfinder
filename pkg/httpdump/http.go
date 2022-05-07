package httpdump

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"time"

	"github.com/gokitx/pkgs/urlx"
	"github.com/sirupsen/logrus"
)

var (
	defaultClient *http.Client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
)

func Do(req Request, timeout time.Duration) (Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp := Response{}
	hreq, err := http.NewRequest(req.Method, req.URL, nil)
	if err != nil {
		return resp, err
	}
	hreq.WithContext(ctx)
	hreq.Header = req.Header.Clone()
	if req.Method != http.MethodGet {
		hreq.PostForm = urlx.CloneUrlValues(req.PostForm)
	}
	hresp, err := defaultClient.Do(hreq)
	if err != nil {
		return resp, err
	}

	resp.Status = hresp.StatusCode
	resp.Header = hresp.Header.Clone()
	if resp.Body != nil {
		defer hresp.Body.Close()

		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		defer func() {
			if buf != nil {
				buf.Reset()
				bufferPool.Put(buf)
			}
			buf = nil
		}()

		if _, err := io.Copy(buf, hresp.Body); err != nil {
			logrus.Errorln("[httputil] copy resp.body error:", err)
		} else {
			resp.Body = buf.Bytes()
		}
	}
	return resp, nil
}
