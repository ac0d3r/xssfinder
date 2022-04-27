package proxy

import (
	"net/url"
	"path/filepath"
	"testing"
)

func TestURL(t *testing.T) {
	URL, _ := url.Parse("http://foo.com/static/ab.css")
	t.Log(URL.Path, filepath.Ext(URL.Path), ignoreRequest(URL.Path))
	URL, _ = url.Parse("http://foo.com/static/ab.js")
	t.Log(URL.Path, filepath.Ext(URL.Path), ignoreRequest(URL.Path))
	URL, _ = url.Parse("http://foo.com/static/ab.png")
	t.Log(URL.Path, filepath.Ext(URL.Path), ignoreRequest(URL.Path))
}
