package crawler

import (
	"bytes"
	"io"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

const (
	anchor  = "a"
	linkRef = "href"
)

// LinksMap is a unique list of web links parsed from the page
type LinksMap = map[string]bool

// Page is the interface that retrieves a web page, parses it, and returns the unique list of urls
type Page struct {
	PageURL string
	Host    string
	Links   LinksMap
}

// NewPage creates a new page for finding links, each separate web request needs a new page
func NewPage(pageURL string, host string) *Page {
	return &Page{
		PageURL: pageURL,
		Host:    host,
		Links:   make(LinksMap),
	}
}

func (pg *Page) Write(p []byte) (n int, err error) {
	err = pg.parser(bytes.NewReader([]byte(p)))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// Links parses the HTML and returns the found links
func (pg *Page) parser(r io.Reader) error {
	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		token := z.Token()
		if tt == html.ErrorToken {
			break
		}
		if token.Data == anchor {
			pg.addLink(token.Attr)
		}
	}
	return nil
}

func (pg *Page) addLink(attrs []html.Attribute) {
	for _, a := range attrs {
		if a.Key != linkRef {
			continue
		}
		if !pg.linkInDomain(a.Val) {
			continue
		}
		pg.Links[pg.addBaseURL(strings.TrimSpace(a.Val))] = true
	}
}

func (pg *Page) linkInDomain(linkURL string) bool {
	parsedURL, err := url.Parse(linkURL)
	if err != nil {
		return false
	}

	// Ignore links in the same page
	if strings.HasPrefix(linkURL, "#") {
		return false
	}

	// Ignore links that refer to link on the previous page backwards
	if strings.HasPrefix(linkURL, "..") {
		return false
	}

	// If the scheme is blank then it must be a simple link such as "/contact"
	extension := filepath.Ext(linkURL)
	if parsedURL.Scheme == "" && (extension == "" || extension == ".html") {
		return true
	}

	// If hosts match then we've found a link in this domain
	if pg.Host == "" || parsedURL.Host == pg.Host {
		return true
	}
	return false
}

func (pg *Page) addBaseURL(linkURL string) string {
	parsedURL, err := url.Parse(linkURL)
	if err != nil {
		log.Printf("Error parsing url: %s", linkURL)
		return linkURL
	}
	if parsedURL.Host != "" {
		return linkURL
	}

	return pg.baseURL(linkURL) + linkURL
}

func (pg *Page) baseURL(linkURL string) string {
	if pg.PageURL == "" {
		return ""
	}
	parsedPageURL, err := url.Parse(pg.PageURL)
	if err != nil {
		log.Printf("Error parsing url: %s", pg.PageURL)
		return pg.PageURL
	}

	pathSeparator := ""
	if !strings.HasPrefix(linkURL, "/") && !strings.HasSuffix(parsedPageURL.Host, "/") {
		pathSeparator = "/"
	}

	return parsedPageURL.Scheme + "://" + parsedPageURL.Host + pathSeparator
}
