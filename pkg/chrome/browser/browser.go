package browser

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

const (
	userAgent = "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.0 Safari/537.36"
)

type Config struct {
	ExecPath   string
	NoHeadless bool
	Incognito  bool
	Proxy      string
}

type Browser struct {
	once sync.Once
	mux  sync.RWMutex

	opts             []chromedp.ExecAllocatorOption
	ctx              *context.Context
	cancelA, cancelC context.CancelFunc
	running          int64
}

func NewBrowser(c Config) *Browser {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", !c.NoHeadless),         // 无头模式
		chromedp.Flag("disable-gpu", true),               // 禁用GPU，不显示GUI
		chromedp.Flag("incognito", c.Incognito),          // 隐身模式启动
		chromedp.Flag("no-sandbox", true),                // 取消沙盒模式
		chromedp.Flag("ignore-certificate-errors", true), // 忽略证书错误
		chromedp.Flag("disable-images", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-xss-auditor", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("allow-running-insecure-content", true),
		chromedp.Flag("disable-webgl", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.UserAgent(userAgent),
		chromedp.WindowSize(1920, 1080),
	)
	if c.Proxy != "" { // 设置浏览器代理
		if _, err := url.Parse(c.Proxy); err == nil {
			opts = append(opts, chromedp.ProxyServer(c.Proxy))
		}
	}
	if c.ExecPath != "" { // 设置浏览器执行路径
		opts = append(opts, chromedp.ExecPath(c.ExecPath))
	}

	browser := &Browser{
		opts: opts,
	}
	return browser
}

func (b *Browser) Start(ctx context.Context) error {
	var err error
	b.once.Do(func() {
		err = b.init(ctx)
	})
	return err
}

func (b *Browser) init(ctx context.Context) error {
	b.mux.Lock()
	defer b.mux.Unlock()

	ctx, cancelA := chromedp.NewExecAllocator(ctx, b.opts...)
	bctx, cancelC := chromedp.NewContext(ctx)

	if err := chromedp.Run(bctx); err != nil {
		return err
	}
	b.ctx = &bctx
	b.cancelA = cancelA
	b.cancelC = cancelC
	return nil
}

func (b *Browser) Scan(tasks chromedp.Tasks) error {
	b.mux.RLock()
	defer b.mux.RUnlock()

	if b.ctx == nil || b.cancelA == nil || b.cancelC == nil {
		return fmt.Errorf("browser engine not started")
	}

	atomic.AddInt64(&b.running, 1)
	defer atomic.AddInt64(&b.running, -1)

	ctx, cancel := chromedp.NewContext(*b.ctx)
	defer cancel()

	return chromedp.Run(ctx, tasks)
}

func (b *Browser) Running() int64 {
	return b.running
}

func (b *Browser) Close() {
	logrus.Debugln("browser close")

	b.mux.Lock()
	defer b.mux.Unlock()
	if b.ctx == nil || b.cancelA == nil || b.cancelC == nil {
		return
	}

	// close main browser
	b.cancelC()
	b.cancelA()
	if err := browser.Close().Do(*b.ctx); err != nil {
		logrus.Errorln("browser page close error:", err)
	}
}
