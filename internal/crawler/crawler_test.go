package crawler_test

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"
	"xsteelej/webcrawler/internal/crawler"
	"xsteelej/webcrawler/internal/request"
)

const crawlHtmlTest = `<!DOCTYPE html>` +
	`<html>` +
	`<body>` +
	`<h2>Absolute URLs</h2>` +
	`<p><a href="https://www.w3.org/">W3C</a></p>` +
	`<p><a href="https://www.google.com/">Google</a></p>` +
	`<p><a href="https://www.google.com/">Google</a></p>` +
	`<h2>Relative URLs</h2>` +
	`<p><a href="html_images.asp">HTML Images</a></p>` +
	`<p><a href="/css/default.asp">CSS Tutorial</a></p>` +
	`</body>` +
	`</html>`

func TestSinglePageCrawler(t *testing.T) {
	crawl, r := NewCrawlerAndMockRequest("")
	r.AddCannedResponse(request.MockResponse{URL: "", ResponseBody: []byte(crawlHtmlTest), StatusCode: http.StatusOK})
	links, err := crawl.Start(context.Background())
	if err != nil {
		log.Fatalf("Error returned when starting crawler: %s", err.Error())
	}
	expectResponse := crawler.LinksMap{
		"":                        true,
		"https://www.w3.org/":     true,
		"https://www.google.com/": true,
		"html_images.asp":         true,
		"/css/default.asp":        true,
	}

	if len(expectResponse) != len(links) {
		log.Fatalf("Expecting %d but got %d", len(expectResponse), len(links))
	}

	for link := range links {
		if _, ok := expectResponse[link]; !ok {
			log.Fatalf("Url: %s not expected", link)
		}
	}
}

func TestCrawlingMockWebPage(t *testing.T) {
	crawl, r := NewCrawlerAndMockRequest("https://www.google.com")

	initialPageResponse, err := ioutil.ReadFile("../../testdata/initial_page_html.txt")
	if err != nil {
		log.Fatal("Error reading ../../testdata/initial_page_html.txt")
	}
	contactsPageResponse, err := ioutil.ReadFile("../../testdata/google_contacts.txt")
	if err != nil {
		log.Fatal("Error reading ../../testdata/google_contacts.txt")
	}

	helpPageResponse, err := ioutil.ReadFile("../../testdata/google_help.txt")
	if err != nil {
		log.Fatal("Error reading ../../testdata/google_help.txt")
	}

	searchPageResponse, err := ioutil.ReadFile("../../testdata/google_search.txt")
	if err != nil {
		log.Fatal("Error reading ../../testdata/google_search.txt")
	}

	testPageResponse, err := ioutil.ReadFile("../../testdata/google_test.txt")
	if err != nil {
		log.Fatal("Error reading ../../testdata/google_test.txt")
	}

	r.AddCannedResponse(request.MockResponse{URL: "https://www.google.com", ResponseBody: []byte(initialPageResponse), StatusCode: http.StatusOK})
	r.AddCannedResponse(request.MockResponse{URL: "https://www.google.com/contacts", ResponseBody: []byte(contactsPageResponse), StatusCode: http.StatusOK})
	r.AddCannedResponse(request.MockResponse{URL: "https://www.google.com/help", ResponseBody: []byte(helpPageResponse), StatusCode: http.StatusOK})
	r.AddCannedResponse(request.MockResponse{URL: "https://www.google.com/search", ResponseBody: []byte(searchPageResponse), StatusCode: http.StatusOK})
	r.AddCannedResponse(request.MockResponse{URL: "https://www.google.com/test", ResponseBody: []byte(testPageResponse), StatusCode: http.StatusOK})

	links, err := crawl.Start(context.Background())
	if err != nil {
		log.Fatalf("Error returned when starting crawler: %s", err.Error())
		return
	}

	expectResponse := crawler.PageLinks{
		"https://www.google.com":          []string{"https://www.google.com/contacts", "https://www.google.com/help", "https://www.google.com/search", "https://www.google.com"},
		"https://www.google.com/contacts": []string{"https://www.google.com/test"},
		"https://www.google.com/help":     []string{"https://www.google.com/search"},
		"https://www.google.com/search":   []string{"https://www.google.com/test", "https://www.google.com/help"},
		"https://www.google.com/test":     []string{"https://www.google.com"},
	}

	if len(expectResponse) != len(links) {
		log.Fatalf("TestCrawlingAMockWebPage - Expecting %d but got %d", len(expectResponse), len(links))
		return
	}

	for page, pageLinks := range links {
		if _, ok := expectResponse[page]; !ok {
			log.Fatalf("Url: %s not expected", page)
			return
		}

		if len(expectResponse[page]) != len(pageLinks) {
			log.Fatalf("Expecting %d links, but got %d for page: %s", len(expectResponse[page]), len(pageLinks), page)
			return
		}

		// Check the expected response is actually in the list of pages
		for _, expectedLink := range expectResponse[page] {
			found := false
			for _, pageLink := range pageLinks {
				if pageLink == expectedLink {
					found = true
					break
				}
			}

			if !found {
				log.Fatalf("Could not find link %s for page %s", expectedLink, page)
			}
		}

	}
}

func NewCrawlerAndMockRequest(initialUrl string) (*crawler.Crawler, *request.MockRequest) {
	u, _ := url.Parse(initialUrl)
	c := crawler.New(initialUrl, u.Host)
	r := &request.MockRequest{}
	c.Requester = r
	return c, r
}
