package crawler_test

import (
	"bytes"
	"io"
	"testing"
	"xsteelej/webcrawler/internal/crawler"
)

const htmlTest = `<!DOCTYPE html>` +
	`<html>` +
	`<body>` +
	`<h2>Absolute URLs</h2>` +
	`<p><a href="https://www.w3.org/">W3C</a></p>` +
	`<p><a href="https://www.google.com/">Google</a></p>` +
	`<h2>Relative URLs</h2>` +
	`<p><a href="html_images.asp">HTML Images</a></p>` +
	`<p><a href="/css/default.asp">CSS Tutorial</a></p>` +
	`</body>` +
	`</html>`

func TestThatParsingEmptyHtmlWorksWithoutError(t *testing.T) {
	const simpleHtml = ``
	pg := crawler.NewPage("", "")
	_, err := io.Copy(pg, bytes.NewReader([]byte(simpleHtml)))
	if err != nil {
		t.Fatalf("Error returned for blank parser")
	}
	if len(pg.Links) > 0 {
		t.Fatal("Links array is not zero")
	}
}

func TestThatNoErrorIsReturnedFromASimplePage(t *testing.T) {
	const simpleHtml = `<!DOCTYPE html><html><body></body></html>`
	pg := crawler.NewPage("", "")
	_, err := io.Copy(pg, bytes.NewReader([]byte(simpleHtml)))
	if err != nil {
		t.Fatalf("Error returned for simple html document: %s", err.Error())
	}
	if len(pg.Links) > 0 {
		t.Fatal("Links array is not zero")
	}
}

func TestThatASingleLinkIsReturned(t *testing.T) {
	const simpleHtml = `<!DOCTYPE html><html><body><a href="https://example.com">Website</a></body></html>`
	pg := crawler.NewPage("", "")
	_, err := io.Copy(pg, bytes.NewReader([]byte(simpleHtml)))
	if err != nil {
		t.Fatalf("Error returned: %s", err.Error())
	}
	if len(pg.Links) != 1 {
		t.Fatal("Links array is not 1")
	}
	if _, ok := pg.Links["https://example.com"]; !ok {
		t.Fatal("Did not find the expected link")
	}
}

func TestThatInvalidHrefDoesNotReturnALink(t *testing.T) {
	const simpleHtml = `<!DOCTYPE html><html><body><a hf="https://example.com">Website</a></body></html>`
	pg := crawler.NewPage("", "")
	_, err := io.Copy(pg, bytes.NewReader([]byte(simpleHtml)))
	if err != nil {
		t.Fatalf("Error returned: %s", err.Error())
	}
	if len(pg.Links) == 1 {
		t.Fatal("Links array is set to 1, it should be 0")
	}
}

func TestThatABlankURLReturnsAnError(t *testing.T) {
	const simpleHtml = `<!DOCTYPE html><html><body><a href="http$://"></a></body></html>`
	pg := crawler.NewPage("", "")
	_, err := io.Copy(pg, bytes.NewReader([]byte(simpleHtml)))
	if err != nil {
		t.Fatalf("Error returned: %s", err.Error())
	}
	if len(pg.Links) != 0 {
		t.Fatal("Links array is set to 1, it should be 0")
	}
}

func TestThatMultipleUrlsAreParsed(t *testing.T) {
	pg := crawler.NewPage("", "")
	_, err := io.Copy(pg, bytes.NewReader([]byte(htmlTest)))

	if err != nil {
		t.Fatalf("Error returned: %s", err.Error())
	}
	expectedLinks := crawler.LinksMap{
		"https://www.w3.org/":     true,
		"https://www.google.com/": true,
		"html_images.asp":         true,
		"/css/default.asp":        true,
	}
	if len(expectedLinks) != len(pg.Links) {
		t.Fatalf("Expecting %d links but got %d", len(expectedLinks), len(pg.Links))
	}
	if len(expectedLinks) != len(pg.Links) {
		t.Fatalf("Expected number links %d does not equal the actual numnber of %d", len(expectedLinks), len(pg.Links))
	}
	for key := range expectedLinks {
		if _, ok := pg.Links[key]; !ok {
			t.Fatalf("Expecting %s but not found", key)
		}
	}
}

func TestThatMultipleUrlsAreFiltered(t *testing.T) {
	pg := crawler.NewPage("https://www.google.com/", "www.google.com")
	_, err := io.Copy(pg, bytes.NewReader([]byte(htmlTest)))
	if err != nil {
		t.Fatalf("Error returned: %s", err.Error())
	}
	expectedLinks := crawler.LinksMap{
		"https://www.google.com/": true,
	}
	if len(expectedLinks) != len(pg.Links) {
		t.Fatalf("Expecting %d links but got %d", len(expectedLinks), len(pg.Links))
	}
	for key := range expectedLinks {
		if _, ok := pg.Links[key]; !ok {
			t.Fatalf("Expecting %s but not found", key)
		}
	}
}
