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
	verbose bool
}

func (l LogrusLogger) Printf(format string, v ...interface{}) {
	if l.verbose {
		// TODO 先忽略此日志
		// logrus.Tracef(format, v...)
	}
}

type Config struct {
	Addr        string
	Verbose     bool     // porxy 代理日志
	TargetHosts []string // 指定目标主机，将忽略其他主机；默认为所有
	ParentProxy string   // 请求的代理
	CaHost      string
}

type MitmServer struct {
	addr    string
	cahost  string
	reqs    sync.Map
	goProxy *goproxy.ProxyHttpServer
	srv     *http.Server

	C <-chan Request
}

func NewMitmServer(conf Config) *MitmServer {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = conf.Verbose
	proxy.Logger = LogrusLogger{verbose: conf.Verbose}

	if conf.CaHost == "" {
		conf.CaHost = "xssfinder.ca"
	}
	proxy.OnRequest(goproxy.DstHostIs(conf.CaHost)).
		DoFunc(DownloadCaHandlerFunc)
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	c := make(chan Request, 5e1)
	mitm := &MitmServer{
		addr:    conf.Addr,
		cahost:  conf.CaHost,
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
	// ignore requests for static resources
	if ignoreRequestWithPath(req.URL.Path) {
		return req, nil
	}
	m.reqs.Store(ctx.Session, MakeRequest(req))
	return req, nil
}

func (m *MitmServer) MakeOnResponse(c chan Request) func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	return func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if resp == nil {
			m.reqs.Delete(ctx.Session)
			return resp
		}

		contentType := resp.Header.Get("Content-Type")
		if strings.Contains(contentType, "text/html") ||
			strings.Contains(contentType, "text/htm") {
			if req, ok := m.reqs.LoadAndDelete(ctx.Session); ok {
				if request, ok := req.(Request); ok {
					logrus.Debugln("[mitm] received:", request.URL)
					request.Response = MakeResponse(request, resp)
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
	logrus.Infoln("[mitm] listen at: ", m.addr)
	logrus.Infoln("[mitm] ca:  ", "http://"+m.cahost)
	return serv.ListenAndServe()
}

func (m *MitmServer) Shutdown(ctx context.Context) error {
	return m.srv.Shutdown(ctx)
}
