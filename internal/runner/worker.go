package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Buzz2d0/xssfinder/pkg/chrome/browser"
	"github.com/Buzz2d0/xssfinder/pkg/chrome/bypass"
	"github.com/Buzz2d0/xssfinder/pkg/chrome/cookies"
	"github.com/Buzz2d0/xssfinder/pkg/chrome/xss/checker"
	"github.com/Buzz2d0/xssfinder/pkg/chrome/xss/dom"
	"github.com/Buzz2d0/xssfinder/pkg/httpdump"
	"github.com/Buzz2d0/xssfinder/pkg/notify"
	"github.com/chromedp/chromedp"
	"github.com/gokitx/pkgs/limiter"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	*limiter.Limiter

	notifier   notify.Notifier
	browser    *browser.Browser
	preActions []chromedp.Action
}

func NewWorker(limitNum int64,
	notifier notify.Notifier,
	browser *browser.Browser) *Worker {
	return &Worker{
		Limiter:    limiter.New(limitNum),
		notifier:   notifier,
		browser:    browser,
		preActions: []chromedp.Action{bypass.BypassHeadlessDetect()},
	}
}

func (w *Worker) Start(ctx context.Context, C <-chan httpdump.Request) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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
			logrus.Infoln("[worker] received task:", task.URL, task.Response.Status)
			go func(ctx context.Context, req httpdump.Request) {
				defer w.Done()
				if err := w.scan(ctx, req); err != nil {
					logrus.Errorln("[worker] scan task error:", err)
				}
			}(ctx, task)
		}
	}
}

func (w *Worker) scan(ctx context.Context, req httpdump.Request) error {
	var preTasks chromedp.Tasks
	preTasks = w.preActions[:]
	if len(req.Cookies) != 0 {
		logrus.Debugf("[worker] cookies: %#v \n", req.Cookies)
		preTasks = append(preTasks, cookies.SetWithHttpCookie(req.Cookies))
	}

	switch req.Method {
	case http.MethodGet:
		// TODO only detecting dom-based XSS
		vuls := make([]dom.VulPoint, 0)
		if err := w.browser.Scan(dom.GenTasks(req.URL, &vuls, time.Second*8, preTasks)); err != nil {
			return err
		}
		if len(vuls) != 0 {
			logrus.Info("[worker] dom-based scan vuls:")
			for _, vulPoint := range vuls {
				logrus.Infof("\t url: %s source: %s sink: %s \n", vulPoint.Url, vulPoint.Source.Label, vulPoint.Sink.Label)
			}
			// TODO 根据 source 构造 payload
			// w.checkDomPoc(ctx, vuls, preTasks)
		}
	}
	return nil
}

func (w *Worker) checkDomPoc(ctx context.Context, points []dom.VulPoint, preActions chromedp.Tasks) {
	var res bool
	for _, point := range points {
		pocUrls, err := dom.GenPocUrls(point)
		if err != nil {
			logrus.Errorln("[worker] gen dom poc urls error:", err)
			continue
		}

		for _, poc := range pocUrls {
			res = false
			if err := w.browser.Scan(checker.GenTasks(poc.Url,
				poc.Rand,
				&res,
				time.Second*5,
				preActions)); err != nil {
				logrus.Errorln("[worker] check dom poc url error:", err)
			}

			if res {
				w.reportDom(poc.Url, point)
			}
		}
	}
}

func (w *Worker) reportDom(url string, point dom.VulPoint) {
	// TODO report
	p, _ := json.MarshalIndent(point, "", "  ")
	logrus.Infof("[report] url: %s\n\ttype: dom-based\n\tdesc: %s\n", url, string(p))
	if w.notifier != nil {
		w.notifier.Notify(url,
			fmt.Sprintf("- type dom-based\n- desc \n\n ```\n%s\n```", string(p)),
		)
	}
}
