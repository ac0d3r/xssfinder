package bypass

import (
	"context"
	"io/ioutil"
	"log"
	"testing"

	"github.com/chromedp/chromedp"
)

func TestByass(t *testing.T) {
	buf, err := screenshot("https://bot.sannysoft.com/")
	if err != nil {
		t.Fatal(err)
	}

	err = saveToFile(buf, "BypassHeadlessDetect.png")
	if err != nil {
		t.Fatal(err)
	}
}

func screenshot(url string) ([]byte, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.0 Safari/537.36"),
	)
	var (
		ctx              context.Context
		cancelA, cancelC context.CancelFunc
	)
	ctx, cancelA = chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelA()

	ctx, cancelC = chromedp.NewContext(ctx)
	defer cancelC()

	defer func() {
		if err := chromedp.Cancel(ctx); err != nil {
			log.Printf("[Screenshot] Cancel chrome error: %v", err)
		}
	}()

	var buf []byte
	if err := chromedp.Run(ctx, chromedp.Tasks{
		BypassHeadlessDetect(),
		chromedp.Navigate(url),
		chromedp.CaptureScreenshot(&buf),
	}); err != nil {
		return nil, err
	}
	return buf, nil
}

func saveToFile(buf []byte, file string) error {
	if err := ioutil.WriteFile(file, buf, 0644); err != nil {
		return err
	}
	return nil
}
