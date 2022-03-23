package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/http"
	nurl "net/url"
	"os"
	"strings"
	"sync"

	"github.com/stupoid/gophercises/link"
)

type URL struct {
	XMLName    xml.Name `xml:"url"`
	Loc        string   `xml:"loc"`
	Lastmod    string   `xml:"lastmod,omitempty"`
	Changefreq string   `xml:"changefreq,omitempty"`
	Priority   float64  `xml:"priority,omitempty"`
}

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	Urls    []*URL   `xml:"url"`
}

func main() {
	urlFlag := flag.String("url", "https://www.sitemaps.org", "url to generate sitemap")
	formatFlag := flag.String("format", "xml", "output file format (xml/txt).")
	fileFlag := flag.String("file", "sitemap.xml", "output file")
	maxDepth := flag.Int("max-depth", 3, "max depth to search")

	flag.Parse()
	rootUrl, err := nurl.Parse(*urlFlag)
	if err != nil {
		log.Fatal(err)
	}

	urls := map[string]struct{}{
		rootUrl.String(): {},
	}
	urlsVisited := make(map[string]struct{})
	urlset := []string{}

	for depth := 0; depth <= *maxDepth || len(urls) > len(urlsVisited); depth++ {
		wg := sync.WaitGroup{}
		ch := make(chan map[string][]string)

		for url := range urls {
			if _, exists := urlsVisited[url]; !exists {
				wg.Add(1)
				go getHrefs(url, urlset, ch, &wg)
				urlsVisited[url] = struct{}{}
			}
		}

		go func() {
			wg.Wait()
			close(ch)
		}()

		for hrefMap := range ch {
			// resUrl may be different from the one used in the GET request
			// e.g. index.php -> index.html
			for resUrl, hrefs := range hrefMap {
				urlset = append(urlset, resUrl)
				for _, href := range hrefs {
					if url, ok := getRelativeUrl(href, *urlFlag); ok {
						urls[url] = struct{}{}
					}
				}
			}
		}
	}
	var out []byte
	switch *formatFlag {
	case "xml":
		XMLURLSet := URLSet{
			XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		}
		for _, url := range urlset {
			XMLURLSet.Urls = append(XMLURLSet.Urls, &URL{Loc: url})
		}
		out, _ = xml.MarshalIndent(XMLURLSet, " ", "  ")
		out = append([]byte(xml.Header), out...)
	case "txt":
		var sb strings.Builder
		for _, url := range urlset {
			sb.WriteString(url + "\n")
		}
		out = []byte(sb.String())
	}
	err = os.WriteFile(*fileFlag, out, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s sitemap saved to %s\n", *urlFlag, *fileFlag)
}

func getHrefs(url string, urlset []string, ch chan<- map[string][]string, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	hrefs, err := link.ParseHref(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
	ch <- map[string][]string{
		resp.Request.URL.String(): hrefs,
	}
}

func getRelativeUrl(rawUrl, rawBaseUrl string) (string, bool) {
	if strings.HasPrefix(rawUrl, "#") {
		return "", false
	}

	url, err := nurl.Parse(rawUrl)
	if err != nil {
		return "", false
	}

	baseUrl, err := nurl.Parse(rawBaseUrl)
	if err != nil {
		return "", false
	}

	if url.Host == "" {
		url.Scheme = baseUrl.Scheme
		url.Host = baseUrl.Host
	}
	if url.Scheme == baseUrl.Scheme && url.Host == baseUrl.Host {
		url.Fragment = ""
		return url.String(), true
	}

	return "", false
}
