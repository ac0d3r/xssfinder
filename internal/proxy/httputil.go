package proxy

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
)

type Request struct {
	Method        string
	URL           string
	Header        http.Header
	Host          string
	Form          url.Values
	PostForm      url.Values
	MultipartForm *multipart.Form
	Response      Response
}

type Response struct {
	Status int
	Header http.Header
	Body   []byte
}

func MakeRequest(req *http.Request) Request {
	r2 := Request{
		Method: req.Method,
		Host:   req.Host,
		URL:    req.URL.String(),
	}
	r2.Header = req.Header.Clone()
	r2.Form = cloneURLValues(req.Form)
	r2.PostForm = cloneURLValues(req.PostForm)
	r2.MultipartForm = cloneMultipartForm(req.MultipartForm)
	return r2
}

func MakeResponse(resp *http.Response) Response {
	r2 := Response{
		Status: resp.StatusCode,
		Header: resp.Header.Clone(),
	}
	data, err := io.ReadAll(resp.Body)
	if err == nil {
		r2.Body = data
	}
	return r2
}

func cloneURLValues(v url.Values) url.Values {
	if v == nil {
		return nil
	}
	// http.Header and url.Values have the same representation, so temporarily
	// treat it like http.Header, which does have a clone:
	return url.Values(http.Header(v).Clone())
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
