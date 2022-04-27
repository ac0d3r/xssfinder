package cookie

import (
	"context"
	"net/http"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func SetOne(name, value,
	domain, path string,
	httpOnly, secure bool,
	expires int64) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		expr := cdp.TimeSinceEpoch(time.Unix(expires, 0))
		return network.SetCookie(name, value).
			WithExpires(&expr).
			WithDomain(domain).
			WithPath(path).
			WithHTTPOnly(httpOnly).
			WithSecure(secure).
			Do(ctx)
	})
}

func SetWithHttpCookie(c []*http.Cookie) chromedp.Action {
	cookies := make([]*network.CookieParam, len(c))
	for i := range c {
		cookies[i].Name = c[i].Name
		cookies[i].Value = c[i].Value
		cookies[i].Domain = c[i].Domain
		cookies[i].Path = c[i].Path
		cookies[i].HTTPOnly = c[i].HttpOnly
		cookies[i].Secure = c[i].Secure
		expr := cdp.TimeSinceEpoch(c[i].Expires)
		cookies[i].Expires = &expr
	}
	return chromedp.ActionFunc(func(ctx context.Context) error {
		return network.SetCookies(cookies).Do(ctx)
	})
}
