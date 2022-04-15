package proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/elazarl/goproxy"
	"github.com/gokitx/pkgs/slices"
)

func init() {
	if err := SetCA(caCert, caKey); err != nil {
		panic(err)
	}
}

type MitmServer struct {
	addr    string
	reqs    sync.Map
	goProxy *goproxy.ProxyHttpServer
	srv     *http.Server
	C       <-chan Request
}

func NewMitmServer(addr string, verbose bool, targets ...string) *MitmServer {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest(goproxy.DstHostIs("xssfinder.ca")).
		DoFunc(DownloadCaHandlerFunc)

	c := make(chan Request, 1e4)
	mitm := &MitmServer{
		addr:    addr,
		goProxy: proxy,
		C:       c,
	}

	var targetFilte goproxy.ReqConditionFunc
	if len(targets) != 0 {
		targetFilte = goproxy.ReqConditionFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
			return slices.ContainsIn(targets, req.Host,
				func(v, sub string) bool {
					return strings.Contains(sub, v) || strings.Contains(v, sub)
				})
		})
		proxy.OnRequest(targetFilte).DoFunc(mitm.OnRequest)
	} else {
		proxy.OnRequest().DoFunc(mitm.OnRequest)
	}

	proxy.OnResponse().DoFunc(mitm.MakeOnResponse(c))
	return mitm
}

func (m *MitmServer) SetParentProxy(parentProxy string) {
	if parentProxy != "" {
		m.goProxy.Tr.Proxy = func(req *http.Request) (*url.URL, error) {
			return url.Parse(parentProxy)
		}
	}
	m.goProxy.Tr.Proxy = http.ProxyFromEnvironment
}

func (m *MitmServer) OnRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	fmt.Println(ctx.Session, req.URL.String())

	m.reqs.Store(ctx.Session, MakeRequest(req))
	return req, nil
}

func (m *MitmServer) MakeOnResponse(c chan Request) func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	return func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		fmt.Println(ctx.Session, resp.Status)

		if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
			if req, ok := m.reqs.LoadAndDelete(ctx.Session); ok {
				if request, ok := req.(Request); ok {
					request.Response = MakeResponse(resp)
					c <- request
				}
			}
		} else {
			m.reqs.Delete(ctx.Session)
		}

		return resp
	}
}

func (m *MitmServer) Run() error {
	serv := &http.Server{
		Addr:    m.addr,
		Handler: m.goProxy,
	}
	m.srv = serv
	return serv.ListenAndServe()
}

func (m *MitmServer) Shutdown(ctx context.Context) error {
	return m.srv.Shutdown(ctx)
}
