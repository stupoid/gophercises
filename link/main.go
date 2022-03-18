package link

import (
	"io"
	"strings"

	"github.com/gammazero/deque"
	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

// getText uses a DFS to find Text elements
// from HTML
func getText(n *html.Node) string {
	var sb strings.Builder
	var nodesToVisit deque.Deque
	nodesToVisit.PushFront(n)
	for nodesToVisit.Len() > 0 {
		i := nodesToVisit.PopFront()
		c := i.(*html.Node)
		if c.Type == html.TextNode {
			sb.WriteString(c.Data)
		}
		if c.FirstChild != nil {
			nodesToVisit.PushFront(c.FirstChild)
		}
		if c.NextSibling != nil {
			nodesToVisit.PushBack(c.NextSibling)
		}
	}
	return strings.TrimSpace(sb.String())
}

// ParseLinks parses HTML using x/net/html and
// uses a DFS to find Anchor elements.
func ParseLinks(r io.Reader) ([]Link, error) {
	links := []Link{}
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	var nodesToVisit deque.Deque
	nodesToVisit.PushFront(doc)
	for nodesToVisit.Len() > 0 {
		i := nodesToVisit.PopFront()
		c := i.(*html.Node)

		if c.Type == html.ElementNode && c.Data == "a" {
			var href, text string
			for _, a := range c.Attr {
				if a.Key == "href" {
					href = a.Val
					break
				}
			}

			// extract text from Anchor element
			// and skip going deeper into this node
			if c.FirstChild != nil {
				text = getText(c.FirstChild)
			}

			links = append(links,
				Link{
					Href: href,
					Text: text,
				},
			)

			if c.NextSibling != nil {
				nodesToVisit.PushBack(c.NextSibling)
			}
		} else {
			if c.FirstChild != nil {
				nodesToVisit.PushFront(c.FirstChild)
			}
			if c.NextSibling != nil {
				nodesToVisit.PushBack(c.NextSibling)
			}
		}
	}

	return links, nil
}
