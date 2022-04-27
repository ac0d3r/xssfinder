package checker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

func CheckGetTypePoc(url, kerword string, res *bool, timeout time.Duration, before ...chromedp.Action) chromedp.Tasks {
	return chromedp.Tasks{
		network.Enable(),
		chromedp.Tasks(before),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var (
				fired    = make(chan struct{})
				captured = make(chan struct{})
			)
			lctx, cancel := context.WithCancel(ctx)
			defer cancel()

			chromedp.ListenTarget(lctx, func(ev interface{}) {
				switch e := ev.(type) {
				case *page.EventJavascriptDialogOpening:
					logrus.Traceln("[CheckPoc]", e)
					if strings.Contains(e.Message, kerword) {
						*res = true
						captured <- struct{}{}
					}
				case *page.EventLoadEventFired:
					fired <- struct{}{}
				}
			})

			var err error
			go func(e *error) {
				_, _, errorText, err := page.Navigate(url).Do(ctx)
				if err != nil {
					*e = err
				}
				if errorText != "" {
					*e = fmt.Errorf("page load error %s", errorText)
				}
			}(&err)

			select {
			case <-captured:
				// 如果成功触发对话框，页面处于 loadding 状态
				if err := page.StopLoading().Do(ctx); err != nil {
					return err
				}
				if err := page.HandleJavaScriptDialog(true).Do(ctx); err != nil {
					return err
				}
				logrus.Traceln("[CheckPoc] captured")
			case <-time.After(timeout):
				if err := page.StopLoading().Do(ctx); err != nil {
					return err
				}
				logrus.Traceln("[CheckPoc] timeout")
			case <-fired:
				logrus.Traceln("[CheckPoc] fired")
			}
			return err
		}),
		page.Close(),
	}
}
