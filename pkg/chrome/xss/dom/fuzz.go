package dom

import (
	"net/url"
	"strings"

	"github.com/gokitx/pkgs/random"
	"github.com/gokitx/pkgs/urlx"
)

type VulPoint struct {
	Url    string     `json:"url"`
	Source TrackChain `json:"source"`
	Sink   TrackChain `json:"sink"`
}

type TrackChain struct {
	Label      string `json:"label"`
	Stacktrace []struct {
		Url    string `json:"url"`
		Line   string `json:"line"`
		Column string `json:"column"`
	} `json:"stacktrace"`
}

var (
	fuzzPrefixes = []string{
		`javascript://alert({{RAND}})//`,
	}
	fuzzSuffixes = []string{
		`'-alert({{RAND}})-'`,
		`"-alert({{RAND}})-"`,
		`-alert({{RAND}})-`,
		`'"><img src=x onerror=alert({{RAND}})>`,
	}
)

func genRand(s string) (string, string) {
	r := random.RandomDigitString(5)
	return strings.ReplaceAll(s, "{{RAND}}", r), r
}

type FuzzUrl struct {
	Url  string `json:"url"`
	Rand string `json:"rand"`
}

func GenPocUrls(point VulPoint) ([]FuzzUrl, error) {
	payloads := make([]FuzzUrl, 0)

	u, err := url.Parse(point.Url)
	if err != nil {
		return nil, err
	}

	oriQuery := u.Query()
	if len(oriQuery) != 0 {
		for k := range oriQuery {
			for _, pre := range fuzzPrefixes {
				pre, rand := genRand(pre)
				q2 := urlx.CloneUrlValues(oriQuery)
				q2.Set(k, pre+q2.Get(k))
				u.RawQuery = q2.Encode()
				payloads = append(payloads, FuzzUrl{
					Url:  u.String(),
					Rand: rand,
				})
			}
			for _, suf := range fuzzSuffixes {
				suf, rand := genRand(suf)
				q2 := urlx.CloneUrlValues(oriQuery)
				q2.Set(k, q2.Get(k)+suf)
				u.RawQuery = q2.Encode()
				payloads = append(payloads, FuzzUrl{
					Url:  u.String(),
					Rand: rand,
				})
			}
		}
		// 还原query
		u.RawQuery = oriQuery.Encode()
	}

	hash := u.Fragment
	if hash != "" {
		for _, pre := range fuzzPrefixes {
			pre, rand := genRand(pre)
			u.Fragment = pre + hash
			payloads = append(payloads, FuzzUrl{
				Url:  u.String(),
				Rand: rand,
			})
		}
		for _, suf := range fuzzSuffixes {
			suf, rand := genRand(suf)
			u.Fragment = hash + suf
			payloads = append(payloads, FuzzUrl{
				Url:  u.String(),
				Rand: rand,
			})
		}
	}

	// TODO referrer
	// if strings.Contains(point.Source.Label, "referrer") {
	// 	for _, suf := range fuzzSuffixes {
	// 		u, rand := genRand(fmt.Sprintf("%s&%s", point.Url, suf))
	// 		payloads = append(payloads, FuzzUrl{
	// 			Url:  u,
	// 			Rand: rand,
	// 		})
	// 	}
	// }

	return payloads, nil
}
