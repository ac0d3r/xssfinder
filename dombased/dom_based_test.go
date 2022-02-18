package dombased

import (
	"context"
	"testing"

	"github.com/chromedp/chromedp"
)

func TestDom(t *testing.T) {
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

	if err := chromedp.Run(ctx, DomBased("https://example.com/#123")); err != nil {
		t.Fatal(err)
	}
}

func TestRex(t *testing.T) {
	ss := scriptContentRex.FindAllStringSubmatch(`<script>
	document.write("Hello2 " + b + "!");
	</script>
	<script>
	console.log("123");
	</script>`, -1)
	// fmt.Printf("%#v \n", ss)
	for i := range ss {
		for j := range ss[i] {
			t.Log(ss[i][j])
		}
	}

}
