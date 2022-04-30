package runner

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/Buzz2d0/xssfinder/internal/options"
	"github.com/Buzz2d0/xssfinder/pkg/chrome/browser"
	"github.com/Buzz2d0/xssfinder/pkg/notify"
	"github.com/Buzz2d0/xssfinder/pkg/proxy"
	"github.com/sirupsen/logrus"
)

type Runner struct {
	opt        *options.Options
	worker     *Worker
	mitmServer *proxy.MitmServer
}

func NewRunner(opt *options.Options) (*Runner, error) {
	var (
		notifier notify.Notifier
		err      error
	)
	if opt.NotifierYaml != "" {
		notifier, err = notify.NewNotifierWithYaml(opt.NotifierYaml)
		if err != nil {
			return nil, err
		}
	}
	browser := browser.NewBrowser(opt.Browser)
	mitm := proxy.NewMitmServer(opt.Mitm)
	worker := NewWorker(5e1, notifier, browser)

	r := &Runner{
		opt:        opt,
		mitmServer: mitm,
		worker:     worker,
	}
	return r, nil
}

func (r *Runner) Start() {
	var (
		err  error
		once sync.Once
	)
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	e := make(chan os.Signal, 1)
	signal.Notify(e, os.Interrupt)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := r.mitmServer.ListenAndServe(); err != nil {
			logrus.Error(err)
			once.Do(func() { close(e) })
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := r.worker.Start(ctx, r.mitmServer.C); err != nil {
			logrus.Error(err)
			once.Do(func() { close(e) })
		}
	}()

	<-e
	cancel()
	if err = r.mitmServer.Shutdown(ctx); err != nil {
		logrus.Error(err)
	}
	r.worker.Wait()
	wg.Wait()
}
