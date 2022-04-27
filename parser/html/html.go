package html

import (
	"io"

	"github.com/Buzz2d0/xssfinder/parser/javascript"
	"github.com/gokitx/pkgs/slicex"
	"golang.org/x/net/html"
)

const (
	inputTag  = "input"
	scriptTag = "script"
)

func GetParams(r io.Reader) ([]string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	params := make([]string, 0)
	var f func(*html.Node)
	f = func(n *html.Node) {
		switch n.Type {
		case html.ElementNode:
			switch n.Data {
			case inputTag:
				for _, a := range n.Attr {
					if a.Key == "name" {
						params = append(params, a.Val)
						break
					}
				}
			case scriptTag:
				if n.FirstChild != nil {
					if vars, err := javascript.GetAllVariable(n.FirstChild.Data); err == nil {
						params = append(params, vars...)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return slicex.RevemoRepByMap(params), nil
}
