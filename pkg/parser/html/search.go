package html

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type LocationType uint8

const (
	InTag LocationType = iota + 1
	InComment
	InHtml
	Inscript
	InStyle
	InAttr
)

type ReflexLocation struct {
	Type    LocationType
	TagName string // in-html
	Content string
}

func MarkReflexLocation(param string, r io.Reader) ([]ReflexLocation, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	var (
		f  func(*html.Node)
		ls []ReflexLocation = make([]ReflexLocation, 0)
	)
	f = func(n *html.Node) {
		switch n.Type {
		case html.ElementNode:
			if strings.Contains(n.Data, param) {
				ls = append(ls, ReflexLocation{
					Type:    InTag,
					TagName: n.Data,
					Content: n.Data,
				})
			}
			for _, a := range n.Attr {
				if strings.Contains(a.Key, param) {
					ls = append(ls, ReflexLocation{
						Type:    InAttr,
						Content: a.Key,
					})
				} else if strings.Contains(a.Val, param) {
					ls = append(ls, ReflexLocation{
						Type:    InAttr,
						Content: a.Val,
					})
				}
			}
			switch n.Data {
			case scriptTag:
				if n.FirstChild != nil &&
					n.FirstChild.Type == html.TextNode &&
					strings.Contains(n.FirstChild.Data, param) {
					ls = append(ls, ReflexLocation{
						Type:    Inscript,
						Content: strings.TrimSpace(n.FirstChild.Data),
					})
				}
			case styleTag:
				if n.FirstChild != nil &&
					n.FirstChild.Type == html.TextNode &&
					strings.Contains(n.FirstChild.Data, param) {
					ls = append(ls, ReflexLocation{
						Type:    InStyle,
						Content: strings.TrimSpace(n.FirstChild.Data),
					})
				}
			default:
				if n.FirstChild != nil &&
					n.FirstChild.Type == html.TextNode &&
					strings.Contains(n.FirstChild.Data, param) {
					ls = append(ls, ReflexLocation{
						Type:    InHtml,
						TagName: n.Data,
						Content: strings.TrimSpace(n.FirstChild.Data),
					})
				}
			}
		case html.CommentNode:
			if strings.Contains(n.Data, param) {
				ls = append(ls, ReflexLocation{
					Type:    InComment,
					Content: strings.TrimSpace(n.Data),
				})
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return ls, nil
}
