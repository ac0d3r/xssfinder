package runner

type Options struct {
	LogVerbose bool
	Mitm       struct {
		Addr        string
		Verbose     bool     // porxy 代理日志
		TargetHosts []string // 指定目标主机，将忽略其他主机；默认为所有
		ParentProxy string   // 请求的代理
	}
	Chrome struct {
		HeadLess  bool
		UserAgent string
		Proxy     string
	}
}

type Option func(Options)

func NewOptions() Options {
	return Options{}
}

func WithLogVerbose(verbose bool) Option {
	return func(opt Options) {
		opt.LogVerbose = verbose
	}
}

func WithMitmAddr(addr string) Option {
	return func(opt Options) {
		opt.Mitm.Addr = addr
	}
}

func WithMitmParentProxy(p string) Option {
	return func(opt Options) {
		opt.Mitm.ParentProxy = p
	}
}

func WithMitmTargetHosts(hosts ...string) Option {
	return func(opt Options) {
		if opt.Mitm.TargetHosts == nil {
			opt.Mitm.TargetHosts = make([]string, 0)
		}
		opt.Mitm.TargetHosts = append(opt.Mitm.TargetHosts, hosts...)
	}
}

func WithChromeHeadLess(headless bool) Option {
	return func(opt Options) {
		opt.Chrome.HeadLess = headless
	}
}

func WithChromeProxy(proxy string) Option {
	return func(opt Options) {
		opt.Chrome.Proxy = proxy
	}
}

func WithChromeUserAgent(useragent string) Option {
	return func(opt Options) {
		opt.Chrome.UserAgent = useragent
	}
}
