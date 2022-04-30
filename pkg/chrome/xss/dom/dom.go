package dom

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/gokitx/pkgs/bytesconv"
	"github.com/sirupsen/logrus"
)

const (
	eventPushVul = "xssfinderPushDomVul"
)

var (
	//go:embed preload.js
	preloadJS string
)

func GenTasks(url string, vuls *[]VulPoint, timeout time.Duration, before ...chromedp.Action) chromedp.Tasks {
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
		runtime.AddBinding(eventPushVul),
		chromedp.Tasks(before),
		loadDomPreloadJavaScript(),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var (
				once  sync.Once
				mux   sync.Mutex
				wg    sync.WaitGroup
				fired = make(chan struct{})
			)
			lctx, cancel := context.WithCancel(ctx)
			defer cancel()

			chromedp.ListenTarget(lctx, func(ev interface{}) {
				switch e := ev.(type) {
				case *fetch.EventRequestPaused:
					logrus.Debugf("[dom-based] EventRequestPaused: %s %d\n", e.Request.URL, e.ResponseStatusCode)
					if e.ResponseStatusCode == 0 {
						return
					}
					wg.Add(1)
					go func() { // convert javascript
						defer wg.Done()
						if err := parseDomHooktResponseJs(lctx, e); err != nil {
							logrus.Errorf("[dom-based] hook %s error: %s\n", e.Request.URL, err)
						}
					}()
				case *runtime.EventBindingCalled:
					switch e.Name {
					case eventPushVul:
						logrus.Traceln("[dom-based] EventBindingCalled", e.Payload)

						points := make([]VulPoint, 0)
						if err := json.Unmarshal([]byte(e.Payload), &points); err != nil {
							logrus.Errorln("[dom-based] json.Unmarshal error:", err)
							return
						}
						mux.Lock()
						*vuls = append(*vuls, points...)
						mux.Unlock()
					}
				case *page.EventLoadEventFired:
					logrus.Traceln("[dom-based] EventLoadEventFired")
					once.Do(func() { close(fired) })
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
			case <-time.After(timeout):
				if err := page.StopLoading().Do(ctx); err != nil {
					return err
				}
			case <-fired:
			}
			wg.Wait()
			return err
		}),
		page.Close(),
	}
}

func loadDomPreloadJavaScript() chromedp.Action {
	return chromedp.ActionFunc(
		func(ctx context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument(preloadJS).Do(ctx)
			return err
		})
}

var (
	scriptContentRex = regexp.MustCompile(`<script[^/>]*?>(?:\s*<!--)?\s*(\S[\s\S]+?\S)\s*(?:-->\s*)?<\/script>`)
)

func parseDomHooktResponseJs(ctx context.Context, event *fetch.EventRequestPaused) error {
	resBody, err := fetch.GetResponseBody(event.RequestID).Do(ctx)
	if err != nil {
		return err
	}

	switch event.ResourceType {
	case network.ResourceTypeDocument:
		ss := scriptContentRex.FindAllSubmatch(resBody, -1)
		for i := range ss {
			convedBody, err := HookParse(bytesconv.BytesToString(ss[i][1]))
			if err != nil {
				logrus.Errorf("[dom-based] hookconv %s error: %s\n", event.Request.URL, err)
				continue
			}
			resBody = bytes.Replace(resBody, ss[i][1], bytesconv.StringToBytes(convedBody), 1)
		}
		return fetch.FulfillRequest(event.RequestID, event.ResponseStatusCode).WithBody(base64.StdEncoding.EncodeToString(resBody)).Do(ctx)
	case network.ResourceTypeScript:
		convertedResBody, err := HookParse(bytesconv.BytesToString(resBody))
		if err != nil {
			return err
		}
		return fetch.FulfillRequest(event.RequestID, event.ResponseStatusCode).WithBody(base64.StdEncoding.EncodeToString(
			bytesconv.StringToBytes(convertedResBody))).Do(ctx)
	}
	return nil
}
