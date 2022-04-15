package bypass

import (
	"context"
	_ "embed"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

//go:embed stealth.min.js
var stealthJS string

func BypassHeadlessDetect() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		_, err := page.AddScriptToEvaluateOnNewDocument(stealthJS).Do(ctx)
		return err
	})
}
