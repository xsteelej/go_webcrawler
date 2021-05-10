package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"
	"xsteelej/webcrawler/internal/crawler"
)

func main() {
	startPageURL, host, err := parseCommandLine()
	if err != nil {
		fmt.Println("Error parsing command line: " + err.Error())
		os.Exit(1)
	}
	log.Printf("Crawling %s...", startPageURL)
	links, err := startCrawler(startPageURL, host)
	if err != nil {
		fmt.Println("Error returned from crawler: " + err.Error())
		os.Exit(1)
	}
	displayLinks(links)
}

func displayLinks(links crawler.PageLinks) {
	keys := make([]string, 0)

	for key := range links {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, pageURL := range keys {
		log.Printf("Page: %s No of Links: %d", pageURL, len(links[pageURL]))
		sort.Strings(links[pageURL])
		for _, link := range links[pageURL] {
			log.Printf("\tLink: %s", link)
		}
		log.Printf("-----")
	}
}

func parseCommandLine() (string, string, error) {
	startPageURL := flag.String("u", "", "The initial page URL to initialise the webcrawler.")
	flag.Parse()
	parsedURL, err := url.Parse(*startPageURL)
	if err != nil {
		return "", "", err
	}
	return *startPageURL, parsedURL.Host, nil
}

func startCrawler(pageURL string, host string) (crawler.PageLinks, error) {
	c := crawler.New(pageURL, host)
	_, err := c.Start(context.Background())
	if err != nil {
		return nil, err
	}
	return c.Links, nil
}
