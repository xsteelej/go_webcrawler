package crawler

import (
	"context"
	"log"
	"sync"
	"xsteelej/webcrawler/internal/request"
)

// PageLinks map of page URLS to an array of link URLS found on that page
type PageLinks = map[string][]string
type pages = map[string]*Page

// Crawler crawls the StartURL and finds all the links in StartURL host
type Crawler struct {
	Requester     request.Requester
	StartURL      string
	Links         PageLinks
	host          string
	maxGoRoutines int
}

// New creates a new Crawler, pass the startURL
func New(startURL string, host string) *Crawler {
	return &Crawler{
		Requester:     request.HTTPRequest{},
		StartURL:      startURL,
		Links:         make(PageLinks),
		host:          host,
		maxGoRoutines: 20,
	}
}

// Start finds all the links, starting from StartURL that are in the same domain
func (c *Crawler) Start(ctx context.Context) (PageLinks, error) {
	pagesList := pages{c.StartURL: NewPage(c.StartURL, c.host)}
	for len(pagesList) > 0 {
		c.findLinks(ctx, pagesList)
		pagesList = c.nextPages(pagesList)
	}

	return c.Links, nil
}

func (c *Crawler) nextPages(pageList pages) pages {
	newPageLinks := make(pages)
	for pageURL, page := range pageList {
		if _, ok := c.Links[pageURL]; ok { // Don't process if page already listed
			continue
		}
		c.Links[pageURL] = make([]string, 0)
		for pageLinkURL := range page.Links {
			c.Links[pageURL] = append(c.Links[pageURL], pageLinkURL)
			newPageLinks[pageLinkURL] = NewPage(pageLinkURL, c.host)
		}
	}
	return newPageLinks
}

func (c *Crawler) findLinks(ctx context.Context, pageList pages) {
	var wg sync.WaitGroup
	semaphore := make(chan bool, c.maxGoRoutines)

	for _, page := range pageList {
		if _, ok := c.Links[page.PageURL]; ok {
			continue
		}

		semaphore <- true
		wg.Add(1)
		go c.getPage(ctx, semaphore, &wg, page)
	}
	wg.Wait()
}

func (c *Crawler) getPage(ctx context.Context, semaphore chan bool, wg *sync.WaitGroup, page *Page) {
	defer func() {
		wg.Done()
		<-semaphore
	}()

	err := c.Requester.Get(ctx, page.PageURL, page)
	if err != nil {
		log.Printf("Error reading page: %s", page.PageURL)
		return
	}
}
