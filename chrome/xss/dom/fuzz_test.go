package dom

import (
	"fmt"
	"net/url"
	"testing"
)

func TestURL(t *testing.T) {
	u, _ := url.Parse("http://localhost:8080/dom_test.html#12323")
	t.Logf("%#v", u)

	u.Fragment = "abc" + "12"
	t.Log(u.String())

	u, _ = url.Parse("https://hyuga.icu?a=1&b=2")
	t.Logf("%#v %#v", u, u.Query())
	q2 := u.Query()
	q2.Set("a", "2333")
	u.RawQuery = q2.Encode()

	t.Log(u.String())
}

func TestGenPocUrls(t *testing.T) {
	urls, _ := GenPocUrls(VulPoint{
		Url: "https://hyuga.icu?a=1&b=2#cc",
		// Url: "http://localhost:8080/dom_test.html#12323",
	})
	logurls(urls)
}

func logurls(urls []FuzzUrl) {
	for _, f := range urls {
		fmt.Printf("%s - %s\n\n", f.Url, f.Rand)
	}
}
