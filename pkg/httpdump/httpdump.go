package httpdump

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"sync"
	"time"

	"github.com/gokitx/pkgs/urlx"
	"github.com/sirupsen/logrus"
)

var (
	bufferPool = sync.Pool{
		New: func() any {
			return bytes.NewBuffer(make([]byte, 1024))
		},
	}
)

type Request struct {
	Method   string
	URL      string
	Header   http.Header
	Host     string
	PostForm url.Values
	Cookies  []Cookie
	Response Response
}

type Cookie struct {
	Name     string
	Value    string
	Domain   string
	Path     string
	HttpOnly bool
	Secure   bool
	Expires  time.Time
}

type Response struct {
	Status int
	Header http.Header
	Body   []byte
}

func MakeRequest(req *http.Request) Request {
	if req.URL.Fragment == "" {
		req.URL.Fragment = "zznq"
		defer func() { req.URL.Fragment = "" }()
	}

	r2 := Request{
		Method: req.Method,
		Host:   req.Host,
		URL:    req.URL.String(),
	}
	r2.Header = req.Header.Clone()
	req.ParseForm()
	r2.PostForm = urlx.CloneUrlValues(req.PostForm)
	return r2
}

func MakeResponse(req Request, resp *http.Response) Response {
	r2 := Response{
		Status: resp.StatusCode,
		Header: resp.Header.Clone(),
		Body:   nil,
	}
	// reponse setcookies
	cookies := resp.Cookies()
	for i := range cookies {
		if cookies[i].Domain == "" {
			cookies[i].Domain = resp.Request.Host
		}
	}
	req.Cookies = append(req.Cookies, dumpCookies(cookies)...)

	if resp.Body != nil {
		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		defer func() {
			if buf != nil {
				buf.Reset()
				bufferPool.Put(buf)
			}
			buf = nil
		}()

		if _, err := io.Copy(buf, resp.Body); err != nil {
			resp.Body.Close()
			logrus.Errorln("[httputil] copy resp.body error:", err)
		} else {
			resp.Body.Close()
			data := buf.Bytes()
			r2.Body = data
			resp.Body = io.NopCloser(bytes.NewReader(data))
		}
	}
	return r2
}

func cloneURL(u *url.URL) *url.URL {
	if u == nil {
		return nil
	}
	u2 := new(url.URL)
	*u2 = *u
	if u.User != nil {
		u2.User = new(url.Userinfo)
		*u2.User = *u.User
	}
	return u2
}

func cloneMultipartForm(f *multipart.Form) *multipart.Form {
	if f == nil {
		return nil
	}
	f2 := &multipart.Form{
		Value: (map[string][]string)(http.Header(f.Value).Clone()),
	}
	if f.File != nil {
		m := make(map[string][]*multipart.FileHeader)
		for k, vv := range f.File {
			vv2 := make([]*multipart.FileHeader, len(vv))
			for i, v := range vv {
				vv2[i] = cloneMultipartFileHeader(v)
			}
			m[k] = vv2
		}
		f2.File = m
	}
	return f2
}

func cloneMultipartFileHeader(fh *multipart.FileHeader) *multipart.FileHeader {
	if fh == nil {
		return nil
	}
	fh2 := new(multipart.FileHeader)
	*fh2 = *fh
	fh2.Header = textproto.MIMEHeader(http.Header(fh.Header).Clone())
	return fh2
}

func dumpCookies(c []*http.Cookie) []Cookie {
	r := make([]Cookie, len(c))
	for i := range c {
		r[i] = Cookie{
			Name:     c[i].Name,
			Value:    c[i].Value,
			Domain:   c[i].Domain,
			Path:     c[i].Path,
			HttpOnly: c[i].HttpOnly,
			Secure:   c[i].Secure,
			Expires:  c[i].Expires,
		}
	}
	return r
}
