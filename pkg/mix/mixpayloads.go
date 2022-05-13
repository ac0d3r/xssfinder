// Package mix used to mix payload with different rules and scopes
package mix

import (
	"net/url"
	"path"
	"strings"

	"github.com/gokitx/pkgs/urlx"
)

func Payloads(u url.URL, payloads []string, rules []Rule, scopes []Scope) []url.URL {
	var urls []url.URL
	for _, scope := range scopes {
		switch scope {
		case ScopeQuery:
			urls = append(urls, mixQuery(u, payloads, rules)...)
		case ScopePath:
			urls = append(urls, mixPath(u, payloads, rules)...)
		case ScopeFragment:
			urls = append(urls, mixFragment(u, payloads, rules)...)
		}
	}
	return urls
}

func mixQuery(u url.URL, payloads []string, rules []Rule) []url.URL {
	if len(u.Query()) == 0 {
		return nil
	}
	baseQuery := u.Query()
	urls := make([]url.URL, len(payloads)*len(rules)*len(baseQuery))

	var index int
	for _, payload := range payloads {
		for key := range baseQuery {
			for _, rule := range rules {
				u.RawQuery = generateQueryWithRule(
					urlx.CloneUrlValues(baseQuery),
					key, payload, rule,
				).Encode()
				urls[index] = u
				index++
			}
		}
	}
	return urls
}

func generateQueryWithRule(query url.Values, key, payload string, rule Rule) url.Values {
	switch rule {
	case RuleAppendPrefix:
		query.Set(key, payload+query.Get(key))
		return query
	case RuleAppendSuffix:
		query.Set(key, query.Get(key)+payload)
		return query
	case RuleReplace:
		query.Set(key, payload)
		return query
	}
	return query
}

func mixPath(u url.URL, payloads []string, rules []Rule) []url.URL {
	var urls []url.URL
	paths := strings.Split(u.Path, "/")
	if len(paths) <= 1 {
		return nil
	}
	for _, payload := range payloads {
		for index := range paths {
			for _, rule := range rules {
				if paths[index] == "" {
					continue
				}
				brefore := paths[index]
				paths[index] = generateStrWithRule(brefore, payload, rule)
				u.Path = path.Join(paths...)
				paths[index] = brefore
				urls = append(urls, u)
			}
		}
	}
	return urls
}

func generateStrWithRule(old, payload string, rule Rule) string {
	switch rule {
	case RuleAppendPrefix:
		return payload + old
	case RuleAppendSuffix:
		return old + payload
	case RuleReplace:
		return payload
	}
	return old
}

func mixFragment(u url.URL, payloads []string, rules []Rule) []url.URL {
	var urls []url.URL
	if u.Fragment == "" {
		return nil
	}
	fragment := u.Fragment
	for _, payload := range payloads {
		for _, rule := range rules {
			u.Fragment = generateStrWithRule(fragment, payload, rule)
			urls = append(urls, u)
		}
	}
	return urls
}

// mixFilter decides whether a payload or url is allowed to be mixed or not.
// TODO: filter numeric or string types
type mixFilter int

const (
	FilterNumeric mixFilter = iota + 1
	FilterString
	FilterIP
	FilterDomain
	FilterURI
	FilterEmail
)

type Rule int

const (
	// RuleAppendPrefix appends the payload to the beginning of the string
	RuleAppendPrefix Rule = iota + 1
	// RuleAppendSuffix appends the payload to the end of the string
	RuleAppendSuffix
	// RuleReplace replaces the string with the payload
	RuleReplace
)

var DefaultMixRules = []Rule{RuleAppendPrefix, RuleAppendSuffix, RuleReplace}

type Scope int

const (
	ScopeQuery Scope = iota + 1
	ScopePath
	ScopeFragment
)

var DefaultScopes = []Scope{ScopeQuery, ScopeFragment}
