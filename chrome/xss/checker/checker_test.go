package checker

import (
	"context"
	"testing"
	"time"

	"github.com/Buzz2d0/xssfinder/logger"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	logger.Init("xssfinder-debug", logger.Config{
		Level:   logrus.DebugLevel,
		NoColor: true,
	})
	m.Run()
}

func TestChecker(t *testing.T) {
	url := "http://localhost:8080/dom_test.html#12323%27%22%3E%3Cimg%20src=x%20onerror=alert(42689)%3E"

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
	)
	var (
		cancelA context.CancelFunc
		ctx     context.Context
	)
	ctx, cancelA = chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelA()

	var cancelC context.CancelFunc
	ctx, cancelC = chromedp.NewContext(ctx)
	defer cancelC()

	defer func() {
		if err := chromedp.Cancel(ctx); err != nil {
			t.Log(err)
		}
	}()

	var res bool
	if err := chromedp.Run(ctx, CheckGetTypePoc(url, "42689", &res, time.Second*12)); err != nil {
		t.Log(err)
	}

	t.Log(res)
}
