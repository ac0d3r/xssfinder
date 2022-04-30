package reflex

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func Reflex(request http.Request) (chromedp.Tasks, error) {
	return chromedp.Tasks{}, nil
}

func reflexGet() {

}

func GetParamsWithURL(rawUrl string) (chromedp.Tasks, error) {
	_, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	var (
		getRequestIDFlag bool
		requestID        network.RequestID
		respBody         []byte
	)
	return chromedp.Tasks{
		network.Enable(),
		chromedp.ActionFunc(func(ctx context.Context) error {
			ch := make(chan struct{})
			once := sync.Once{}
			lCtx, cancel := context.WithCancel(ctx)
			defer cancel()

			chromedp.ListenTarget(lCtx, func(ev interface{}) {
				switch e := ev.(type) {
				case *network.EventResponseReceived:
					if !getRequestIDFlag && e.Type == "Document" {
						getRequestIDFlag = true
						requestID = e.RequestID
					}
				case *page.EventLoadEventFired:
					once.Do(func() {
						close(ch)
					})
				}
			})

			{ // chromedp.Navigate()
				_, _, errorText, err := page.Navigate(rawUrl).Do(ctx)
				if err != nil {
					return err
				}
				if errorText != "" {
					return fmt.Errorf("page load error %s", errorText)
				}
			}

			// notice: should cancel the context and release the bounding
			// Or it may cause resource leak
			select {
			case <-time.After(time.Second * 30):
				if err := page.StopLoading().Do(ctx); err != nil {
					return err
				}
			case <-ch:
			}

			// Get response body
			if b, err := network.GetResponseBody(requestID).Do(ctx); err == nil {
				respBody = b
			}

			fmt.Println(string(respBody), rawUrl)
			return nil
		}),
		page.Close(),
	}, nil
}
