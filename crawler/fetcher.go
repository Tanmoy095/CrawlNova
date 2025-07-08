package crawler

import (
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func ExtractLinks(body io.Reader, base string) ([]string, error) {
	//initialize the valid link slice.....
	var links []string

	//baseUrl will return a object
	baseUrl, err := url.Parse(base)
	if err != nil {
		return nil, err

	}
	//tokenizer breaks the html body into small pieces..one token at a time..

	tokenizer := html.NewTokenizer(body)
	//now loop through html tokens............
	for {
		tt := tokenizer.Next() // it returns token type, like html.errorToken, tagtoke, etc......

		switch tt {
		case html.ErrorToken:
			return links, nil
		case html.StartTagToken, html.SelfClosingTagToken:
			//get the token and check for <a> tag
			t := tokenizer.Token() //returns a tag name t.Data and tag attribute t.Attt
			if t.Data == "a" {
				for _, attr := range t.Attr {
					if attr.Key == "href" {
						href := strings.TrimSpace(attr.Val)

						//skip empty fragment or js loink
						if href == "" || strings.HasPrefix(href, "#") || strings.HasPrefix(href, "javascript") {
							continue
						}

						u, err := url.Parse(href)
						if err != nil {
							continue

						}
						absoluteURL := baseUrl.ResolveReference(u)
						links = append(links, absoluteURL.String())
					}
				}

			}

		}

	}

}
