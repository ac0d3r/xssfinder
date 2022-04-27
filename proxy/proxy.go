package proxy

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/elazarl/goproxy"
	"github.com/gokitx/pkgs/slicex"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := SetCA(caCert, caKey); err != nil {
		panic(err)
	}
}

type LogrusLogger struct {
}

func (l LogrusLogger) Printf(format string, v ...interface{}) {
	logrus.Tracef(format, v...)
}

type Config struct {
	Addr        string
	Verbose     bool     // porxy 代理日志
	TargetHosts []string // 指定目标主机，将忽略其他主机；默认为所有
	ParentProxy string   // 请求的代理
}

type MitmServer struct {
	addr    string
	reqs    sync.Map
	goProxy *goproxy.ProxyHttpServer
	srv     *http.Server

	C <-chan Request
}

func NewMitmServer(conf Config) *MitmServer {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = conf.Verbose
	proxy.Logger = LogrusLogger{}

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest(goproxy.DstHostIs("xssfinder.ca")).
		DoFunc(DownloadCaHandlerFunc)

	c := make(chan Request, 5e1)
	mitm := &MitmServer{
		addr:    conf.Addr,
		goProxy: proxy,
		C:       c,
	}

	if len(conf.ParentProxy) != 0 {
		mitm.SetParentProxy(conf.ParentProxy)
	}

	var targetFilte goproxy.ReqConditionFunc
	if len(conf.TargetHosts) != 0 {
		targetFilte = goproxy.ReqConditionFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
			return slicex.ContainsIn(conf.TargetHosts, req.Host,
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
	m.goProxy.Tr.Proxy = func(req *http.Request) (*url.URL, error) {
		return url.Parse(parentProxy)
	}
}

func (m *MitmServer) OnRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	logrus.Debugln("[mitm]OnRequest ", ctx.Session, req.URL.String())

	// ignore requests for static resources
	if ignoreRequest(req.URL.Path) {
		return req, nil
	}
	m.reqs.Store(ctx.Session, MakeRequest(req))
	return req, nil
}

func (m *MitmServer) MakeOnResponse(c chan Request) func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	return func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		logrus.Debugln("[mitm]OnResponse ", ctx.Session, resp)

		contentType := resp.Header.Get("Content-Type")
		if strings.Contains(contentType, "text/html") ||
			strings.Contains(contentType, "text/htm") {
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

func (m *MitmServer) ListenAndServe() error {
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
