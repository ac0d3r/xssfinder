package browser

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Buzz2d0/xssfinder/logger"
	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	logger.Init("xssfinder-debug", logger.Config{
		Level:   logrus.TraceLevel,
		NoColor: true,
	})
	m.Run()
}

func TestBrowser(t *testing.T) {
	browser := NewBrowser(Config{
		NoHeadless: true,
	})

	t.Log(browser.Start(context.Background()))

	for i := 0; i < 3; i++ {
		t.Log(browser.Scan(chromedp.Tasks{chromedp.Navigate("https://www.baidu.com")}))
		time.Sleep(5 * time.Second)
	}
	browser.Close()
}

func TestChromedpTabs(t *testing.T) {
	// https://github.com/chromedp/chromedp/issues/824
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	bctx, cancel := chromedp.NewContext(ctx) // the next line does not have any effect because there are not browsers created by this context.

	// call Run on browser context to start a browser first
	chromedp.Run(bctx)
	defer cancel()

	for i := 0; i < 3; i++ {
		tctx, cancel := chromedp.NewContext(bctx)
		// since the parent context has browser allocated,
		// the Run will create new tab in the browser.
		chromedp.Run(tctx, chromedp.Navigate("https://www.baidu.com"))
		if err := browser.Close().Do(tctx); err != nil {
			log.Println(err)
		}
		cancel()
		time.Sleep(2 * time.Second)
	}

	if err := browser.Close().Do(bctx); err != nil {
		log.Println(err)
	}
}

func TestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-ctx.Done()
		fmt.Println("123")
	}()

	time.Sleep(time.Second)
	fmt.Println("xxx")
}
