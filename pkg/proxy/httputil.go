package proxy

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"path/filepath"
	"sync"

	"github.com/gokitx/pkgs/slicex"
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
	Method        string
	URL           string
	Header        http.Header
	Host          string
	Form          url.Values
	PostForm      url.Values
	MultipartForm *multipart.Form
	Cookies       []*http.Cookie
	Response      Response
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
	r2.Form = urlx.CloneUrlValues(req.Form)
	r2.PostForm = urlx.CloneUrlValues(req.PostForm)
	r2.MultipartForm = cloneMultipartForm(req.MultipartForm)
	return r2
}

func MakeResponse(req Request, resp *http.Response) Response {
	r2 := Response{
		Status: resp.StatusCode,
		Header: resp.Header.Clone(),
		Body:   nil,
	}
	cookies := resp.Cookies()
	for i := range cookies {
		if cookies[i].Domain == "" {
			cookies[i].Domain = resp.Request.Host
		}
	}
	req.Cookies = cookies

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

var (
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
	exts = []string{
		".aac",
		".abw",
		".arc",
		".avif",
		".avi",
		".azw",
		".bin",
		".bmp",
		".bz",
		".bz2",
		".cda",
		".csh",
		".css",
		".csv",
		".doc",
		".docx",
		".eot",
		".epub",
		".gz",
		".gif",
		".ico",
		".ics",
		".jpeg",
		".jpg",
		".js",
		".json",
		".jsonld",
		".mid",
		".midi",
		".mjs",
		".mp3",
		".mp4",
		".mpeg",
		".mpkg",
		".odp",
		".ods",
		".odt",
		".oga",
		".ogv",
		".ogx",
		".opus",
		".otf",
		".png",
		".pdf",
		".ppt",
		".pptx",
		".rar",
		".rtf",
		".sh",
		".svg",
		".swf",
		".tar",
		".tif ",
		".tiff",
		".ts",
		".ttf",
		".txt",
		".vsd",
		".wav",
		".weba",
		".webm",
		".webp",
		".woff",
		".woff2",
		".xls",
		".xlsx",
		".xml",
		".xul",
		".zip",
		".3gp",
		".3g2",
		".7z",
	}
)

func ignoreRequestWithPath(path string) bool {
	return slicex.ContainsIn(exts,
		filepath.Ext(path))
}
