package options

import (
	"github.com/Buzz2d0/xssfinder/internal/logger"
	"github.com/Buzz2d0/xssfinder/pkg/chrome/browser"
	"github.com/Buzz2d0/xssfinder/pkg/proxy"
	"github.com/sirupsen/logrus"
)

type Options struct {
	Debug        bool
	Verbose      bool
	NotifierYaml string
	Mitm         proxy.Config
	Log          logger.Config
	Browser      browser.Config
}

func (o *Options) LogLevel() logrus.Level {
	l := logrus.InfoLevel
	if o.Debug {
		l = logrus.DebugLevel
	}
	if o.Verbose {
		l = logrus.TraceLevel
	}
	return l
}

func New() *Options {
	return &Options{
		Mitm: proxy.Config{
			Addr: "127.0.0.1:8222",
		},
		Log: logger.Config{
			Level: logrus.InfoLevel,
		},
	}
}

func (o *Options) Set(opts ...Option) *Options {
	for i := range opts {
		opts[i](o)
	}
	return o
}

type Option func(*Options)

func WithDebug(debug bool) Option {
	return func(opt *Options) {
		opt.Debug = debug
	}
}

func WithVeryVerbose(verbose bool) Option {
	return func(opt *Options) {
		opt.Verbose = verbose
	}
}

func WithNotifierYaml(f string) Option {
	return func(opt *Options) {
		opt.NotifierYaml = f
	}
}

// log
func WithLogOutJson(o bool) Option {
	return func(opt *Options) {
		opt.Log.OutJson = o
	}
}

// mitm
func WithMitmAddr(addr string) Option {
	return func(opt *Options) {
		opt.Mitm.Addr = addr
	}
}

func WithMitmVerbose(verbose bool) Option {
	return func(opt *Options) {
		opt.Mitm.Verbose = verbose
	}
}

func WithMitmTargetHosts(hosts ...string) Option {
	return func(opt *Options) {
		if opt.Mitm.TargetHosts == nil {
			opt.Mitm.TargetHosts = make([]string, 0)
		}
		opt.Mitm.TargetHosts = append(opt.Mitm.TargetHosts, hosts...)
	}
}

func WithMitmParentProxy(p string) Option {
	return func(opt *Options) {
		opt.Mitm.ParentProxy = p
	}
}

// browser
func WithBrowserExecPath(p string) Option {
	return func(opt *Options) {
		opt.Browser.ExecPath = p
	}
}

func WithBrowserNoHeadless(n bool) Option {
	return func(opt *Options) {
		opt.Browser.NoHeadless = n
	}
}

func WithBrowserIncognito(i bool) Option {
	return func(opt *Options) {
		opt.Browser.Incognito = i
	}
}

func WithBrowserProxy(p string) Option {
	return func(opt *Options) {
		opt.Browser.Proxy = p
	}
}
