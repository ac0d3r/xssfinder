package cookies

import (
	"context"

	"github.com/Buzz2d0/xssfinder/pkg/httpdump"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func SetWithHttpCookie(c []httpdump.Cookie) chromedp.Action {
	cookies := make([]*network.CookieParam, len(c))
	for i := range c {

		expr := cdp.TimeSinceEpoch(c[i].Expires)
		cookies[i] = &network.CookieParam{
			Name:     c[i].Name,
			Value:    c[i].Value,
			Domain:   c[i].Domain,
			Path:     c[i].Path,
			HTTPOnly: c[i].HttpOnly,
			Secure:   c[i].Secure,
			Expires:  &expr,
		}
	}
	return chromedp.ActionFunc(func(ctx context.Context) error {
		return network.SetCookies(cookies).Do(ctx)
	})
}
