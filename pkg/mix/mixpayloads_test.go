package mix

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

var urlTestCases = []struct {
	urlOBJ url.URL
}{
	{
		urlOBJ: url.URL{
			Scheme:   "https",
			Host:     "www.hackerone.com",
			Path:     "/pathA/pathB/pathC",
			RawQuery: "queryA=111&queryB=222",
			Fragment: "fragment",
		},
	},
}

var payloads = []string{
	"moond4rk",
	"buzz2d0",
}

func TestMixPayloads(t *testing.T) {
	t.Parallel()
	var urls []url.URL
	for _, testCase := range urlTestCases {
		urls = Payloads(testCase.urlOBJ, payloads, DefaultMixRules, DefaultScopes)
	}
	for _, v := range urls {
		fmt.Println(v.String())
	}
}

func TestMixQuery(t *testing.T) {
	t.Parallel()
	var urls []url.URL
	var queryIndex int
	for _, tc := range urlTestCases {
		queryIndex += len(tc.urlOBJ.Query())
		urls = append(urls, mixQuery(tc.urlOBJ, payloads, DefaultMixRules)...)
	}
	if len(urls) != len(urlTestCases)*len(payloads)*len(DefaultMixRules)*queryIndex {
		t.Errorf("Expected %d urls, got %d", len(urlTestCases)*len(payloads)*len(DefaultMixRules)*queryIndex, len(urls))
	}
}

func TestMixPath(t *testing.T) {
	t.Parallel()
	var urls []url.URL
	pathIndex := 0
	for _, tc := range urlTestCases {
		for _, p := range strings.Split(tc.urlOBJ.Path, "/") {
			if p != "" {
				pathIndex++
			}
		}
		urls = append(urls, mixPath(tc.urlOBJ, payloads, DefaultMixRules)...)
	}
	if len(urls) != len(urlTestCases)*len(payloads)*len(DefaultMixRules)*pathIndex {
		t.Errorf("Expected %d urls, got %d", len(urlTestCases)*len(payloads)*len(DefaultMixRules)*pathIndex, len(urls))
	}
}

func TestMixFragment(t *testing.T) {
	t.Parallel()
	var urls []url.URL
	var index int
	for _, tc := range urlTestCases {
		if tc.urlOBJ.Fragment != "" {
			index++
		}
		urls = append(urls, mixFragment(tc.urlOBJ, payloads, DefaultMixRules)...)
	}
	if len(urls) != len(urlTestCases)*len(payloads)*len(DefaultMixRules)*index {
		t.Errorf("Expected %d urls, got %d", len(urlTestCases)*len(payloads)*len(DefaultMixRules), len(urls))
	}
}
