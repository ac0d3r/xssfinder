package reflex

import (
	"context"
	"testing"

	"github.com/chromedp/chromedp"
)

func TestRefex(t *testing.T) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		// chromedp.ProxyServer("http://127.0.0.1:7890"),
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
			t.Fatal(err)
		}
	}()

	tasks, _ := GetParamsWithURL("https://example.com")

	if err := chromedp.Run(ctx, tasks); err != nil {
		t.Fatal(err)
	}
}
