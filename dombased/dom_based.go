package dombased

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/base64"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/Buzz2d0/xssfinder/dombased/hookconv"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

var (
	//go:embed js/preload.js
	preloadJS string
	//go:embed js/bridge.js
	bridgeJS string
)

func DomBased(url string) chromedp.Tasks {
	return chromedp.Tasks{
		runtime.Enable(),
		network.Enable(),
		// 开启响应拦截
		fetch.Enable().WithPatterns([]*fetch.RequestPattern{
			{
				URLPattern:   "*",
				ResourceType: network.ResourceTypeDocument,
				RequestStage: fetch.RequestStageResponse,
			},
			{
				URLPattern:   "*",
				ResourceType: network.ResourceTypeScript,
				RequestStage: fetch.RequestStageResponse,
			},
		}),
		runtime.AddBinding("PushDomVul"),
		chromedp.ActionFunc( // load javascript
			func(ctx context.Context) error {
				var err error
				_, err = page.AddScriptToEvaluateOnNewDocument(preloadJS).Do(ctx)
				if err != nil {
					return err
				}
				_, err = page.AddScriptToEvaluateOnNewDocument(bridgeJS).Do(ctx)
				return err
			}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var (
				ch   = make(chan struct{})
				once = sync.Once{}
				wg   sync.WaitGroup
			)
			lctx, cancel := context.WithCancel(ctx)
			defer cancel()

			chromedp.ListenTarget(lctx, func(ev interface{}) {
				switch e := ev.(type) {
				case *fetch.EventRequestPaused:
					fmt.Printf("EventRequestPaused: %s %d\n", e.Request.URL, e.ResponseStatusCode)

					wg.Add(1)
					go func() { // convert javascript
						defer wg.Done()
						if err := convertResJs(lctx, e); err != nil {
							fmt.Printf("[EventRequestPaused] %s error: %s\n", e.Request.URL, err)
						}
					}()
				case *runtime.EventBindingCalled:
					switch e.Name {
					case "PushDomVul":
						fmt.Println(e.Payload)
					}
				case *page.EventLoadEventFired:
					once.Do(func() {
						close(ch)
					})
				}
			})

			{ // chromedp.Navigate()
				_, _, errorText, err := page.Navigate(url).Do(ctx)
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
			case <-time.After(time.Second * 60):
				if err := page.StopLoading().Do(ctx); err != nil {
					return err
				}
			case <-ch:
			}
			wg.Wait()
			return nil
		}),
		page.Close(),
	}
}

var (
	scriptContentRex = regexp.MustCompile(`<script[^/>]*?>(?:\s*<!--)?\s*(\S[\s\S]+?\S)\s*(?:-->\s*)?<\/script>`)
)

func convertResJs(ctx context.Context, e *fetch.EventRequestPaused) error {
	resBody, err := fetch.GetResponseBody(e.RequestID).Do(ctx)
	if err != nil {
		return err
	}

	switch e.ResourceType {
	case network.ResourceTypeDocument:
		ss := scriptContentRex.FindAllSubmatch(resBody, -1)
		for i := range ss {
			convedBody, err := hookconv.HookConv(string(ss[i][1]))
			if err != nil {
				fmt.Printf("[HookConv] %s error: %s\n", e.Request.URL, err)
				continue
			}
			resBody = bytes.Replace(resBody, ss[i][1], []byte(convedBody), 1)
		}
		return fetch.FulfillRequest(e.RequestID, e.ResponseStatusCode).WithBody(base64.StdEncoding.EncodeToString(resBody)).Do(ctx)
	case network.ResourceTypeScript:
		convertedResBody, err := hookconv.HookConv(string(resBody))
		if err != nil {
			return err
		}
		return fetch.FulfillRequest(e.RequestID, e.ResponseStatusCode).WithBody(base64.StdEncoding.EncodeToString([]byte(convertedResBody))).Do(ctx)
	}
	return nil
}
