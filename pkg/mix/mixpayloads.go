// Package mix used to mix payload with different rules and scopes
package mix

import (
	"net/url"
	"path"
	"strings"
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
	// var urls []url.URL
	index := 0
	for _, payload := range payloads {
		for key := range baseQuery {
			query := CloneValues(baseQuery)
			for _, rule := range rules {
				newQuery := generateQueryWithRule(query, key, payload, rule)
				u.RawQuery = newQuery.Encode()
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
				paths[index] = generateStrWithRule(paths[index], payload, rule)
				newPath := CloneSlices(paths)
				u.Path = path.Join(newPath...)
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

func CloneValues(values url.Values) url.Values {
	if values == nil {
		return nil
	}
	// find total number of values.
	nv := 0
	for _, vv := range values {
		nv += len(vv)
	}
	// shared backing array for headers' values
	sv := make([]string, nv)
	v2 := make(url.Values, len(values))
	for k, vv := range values {
		n := copy(sv, vv)
		v2[k] = sv[:n:n]
		sv = sv[n:]
	}
	return v2
}

// CloneSlices returns a copy of the given slice.
// based on golang 1.18+ generics
func CloneSlices[T any](s []T) []T {
	if s == nil {
		return nil
	}
	return append([]T{}, s...)
}
