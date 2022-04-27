package runner

import (
	"github.com/Buzz2d0/xssfinder/chrome/browser"
	"github.com/Buzz2d0/xssfinder/logger"
	"github.com/Buzz2d0/xssfinder/proxy"
	"github.com/sirupsen/logrus"
)

type Options struct {
	Debug       bool
	VeryVerbose bool
	Mitm        proxy.Config
	Log         logger.Config
	Browser     browser.Config
}

func (o *Options) LogLevel() logrus.Level {
	l := logrus.InfoLevel
	if o.Debug {
		l = logrus.DebugLevel
	}
	if o.VeryVerbose {
		l = logrus.TraceLevel
	}
	return l
}

func NewOptions() *Options {
	return &Options{
		Mitm: proxy.Config{
			Addr: ":8080",
		},
		Log: logger.Config{
			Level: logrus.InfoLevel,
		},
	}
}

func (o *Options) Set(opts ...Option) {
	for i := range opts {
		opts[i](o)
	}
}

type Option func(*Options)

func WithDebug(debug bool) Option {
	return func(opt *Options) {
		opt.Debug = debug
	}
}

func WithVeryVerbose(vverbose bool) Option {
	return func(opt *Options) {
		opt.VeryVerbose = vverbose
	}
}

// log
func WithLogOutJson(o bool) Option {
	return func(opt *Options) {
		opt.Log.OutJson = o
	}
}

func WithLogNoColor(n bool) Option {
	return func(opt *Options) {
		opt.Log.NoColor = n
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
