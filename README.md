# Webcrawler
Author: John Steele

A simple Web crawler that finds all the links in the same domain as the URL

## Running tests

```go test -v -coverprofile cover.out ./...```

Example output:

```
PASS
coverage: 93.8% of statements
ok  	xsteelej/webcrawler/internal/crawler	0.088s	coverage: 93.8% of statements
```

## Linting

```golangci-lint run```

## Compile and run
```
cd cmd/webcrawler
go build -o webcrawler main.go
./webcrawler -u https://www.monzo.com
```

## Description of Design

![alt text](https://github.com/xsteelej/go_webcrawler/blob/main/docs/diagram.png?raw=true)

### Package: ```crawler```

#### Type: ```Crawler```
Is responsible for orchestrating the main logic for generating the list of pages and links. It is initially seeded with the starting URL, it gets each page and generates a list of new pages to visit.

A new Go routine is created to retrieve and parse each page, the number of Go routines are limited to 20 by default.

#### Type: ```Page```

Implements the io.Writer interface and is responsible for parsing the page and finding the links in the domain being parsed.

### Package: ```request```

#### Interface: ```Requester```
Defines the single method: ```Get(ctx context.Context, uri string, w io.Writer) error``` which implementations use to do the actual work of getting the raw page. The ```w io.Writer``` parameter is used to write the raw bytes to an object that implements that interface, such as the ```Page``` type.

#### Type: ```HTTPRequest```
Implements the Requester interface and performs http network calls to get raw page data.

#### Type: ```MockResponse```
Implements the Requester interface and is used for testing the functionality of ```Crawler``` and ```Page```. Tests configure instances of this object with canned response data for unit test scenarios. For testing, this is the only Mock object required.
