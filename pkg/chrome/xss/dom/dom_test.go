package dom

import (
	"context"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	logrus.SetLevel(logrus.FatalLevel)
	m.Run()
}

func TestDom(t *testing.T) {
	url := "http://localhost:8080/dom_test.html#123232"
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
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
			t.Log(err)
		}
	}()

	vuls := make([]VulPoint, 0)
	if err := chromedp.Run(ctx, GenTasks(url, &vuls, time.Second*8)); err != nil {
		t.Log(err)
	}

	t.Log(vuls)
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
