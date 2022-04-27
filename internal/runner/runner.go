package runner

import (
	"context"
	"net/http"
	"time"

	"github.com/Buzz2d0/xssfinder/chrome/browser"
	"github.com/Buzz2d0/xssfinder/chrome/bypass"
	"github.com/Buzz2d0/xssfinder/chrome/cookie"
	"github.com/Buzz2d0/xssfinder/chrome/xss/checker"
	"github.com/Buzz2d0/xssfinder/chrome/xss/dom"
	"github.com/Buzz2d0/xssfinder/proxy"
	"github.com/chromedp/chromedp"
	"github.com/gokitx/pkgs/limiter"
	"github.com/sirupsen/logrus"
)

type Runner struct {
	worker     *ScanWorker
	mitmServer *proxy.MitmServer
}

func NewRunner(opt *Options) *Runner {
	browser := browser.NewBrowser(opt.Browser)
	mitm := proxy.NewMitmServer(opt.Mitm)
	worker := NewScanWorker(5e1, browser)

	r := &Runner{
		mitmServer: mitm,
		worker:     worker,
	}
	return r
}

func (r *Runner) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := r.worker.Run(ctx, r.mitmServer.C); err != nil {
			logrus.Error(err)
		}
	}()
	return r.mitmServer.ListenAndServe()
}

func (r *Runner) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := r.mitmServer.Shutdown(ctx)

	logrus.Debug(r.worker.Count())
	r.worker.Wait()
	return err
}

type ScanWorker struct {
	*limiter.Limiter

	browser    *browser.Browser
	preActions []chromedp.Action
}

func NewScanWorker(limitNum int64, browser *browser.Browser) *ScanWorker {
	return &ScanWorker{
		Limiter:    limiter.New(limitNum),
		browser:    browser,
		preActions: []chromedp.Action{bypass.BypassHeadlessDetect()},
	}
}

func (w *ScanWorker) Run(ctx context.Context, C <-chan proxy.Request) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// first start browser engine
	if err := w.browser.Start(ctx); err != nil {
		return err
	}
	defer w.browser.Close()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case task := <-C:
			w.Allow()
			go func(ctx context.Context, req proxy.Request) {
				defer w.Done()
				logrus.Debugln("[ScanWorker] task:", req.URL, req.Response.Status)
				if err := w.scan(ctx, req); err != nil {
					logrus.Errorln("[ScanWorker] scan", err)
				}
			}(ctx, task)
		}
	}
}

func (w *ScanWorker) scan(ctx context.Context, req proxy.Request) error {
	var preTasks chromedp.Tasks
	preTasks = w.preActions[:]

	if len(req.Cookies) != 0 {
		preTasks = append(preTasks, cookie.SetWithHttpCookie(req.Cookies))
	}

	// TODO only detecting dom-based XSS
	if req.Method == http.MethodGet {
		vuls := make([]dom.VulPoint, 0)
		if err := w.browser.Scan(dom.DomBased(req.URL, &vuls, time.Second*5, preTasks)); err != nil {
			return err
		}
		if len(vuls) != 0 {
			w.checkPoc(ctx, vuls, preTasks)
		}
	}

	return nil
}

func (w *ScanWorker) checkPoc(ctx context.Context, points []dom.VulPoint, preActions chromedp.Tasks) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var res bool
	for _, point := range points {
		pocUrls, err := dom.GenPocUrls(point)
		if err != nil {
			logrus.Errorln("[ScanWorker] GenPocUrls", err)
			continue
		}

		for _, poc := range pocUrls {
			res = false
			if err := w.browser.Scan(checker.CheckGetTypePoc(poc.Url, poc.Rand, &res, time.Second*5, preActions)); err != nil {
				logrus.Errorln("[ScanWorker] Scan", err)
			}

			if res {
				w.report(point, poc.Url)
			}
		}
	}

}

func (w *ScanWorker) report(point dom.VulPoint, poc string) {
	// TODO
	logrus.Infoln(point, poc)
}
