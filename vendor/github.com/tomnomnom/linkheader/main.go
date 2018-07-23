// Package linkheader provides functions for parsing HTTP Link headers
package linkheader

import (
	"fmt"
	"strings"
)

// A Link is a single URL and related parameters
type Link struct {
	URL    string
	Rel    string
	Params map[string]string
}

// HasParam returns if a Link has a particular parameter or not
func (l Link) HasParam(key string) bool {
	for p := range l.Params {
		if p == key {
			return true
		}
	}
	return false
}

// Param returns the value of a parameter, or an error on failure
func (l Link) Param(key string) (string, error) {
	for k, v := range l.Params {
		if key == k {
			return v, nil
		}
	}
	return "", fmt.Errorf("Could not find param '%s'", key)
}

// Links is a slice of Link structs
type Links []Link

// FilterByRel filters a group of Links by the provided Rel attribute
func (l Links) FilterByRel(r string) Links {
	links := make(Links, 0)
	for _, link := range l {
		if link.Rel == r {
			links = append(links, link)
		}
	}
	return links
}

// Parse parses a raw Link header in the form:
//   <url>; rel="foo", <url>; rel="bar"; wat="dis"
// returning a slice of Link structs
func Parse(raw string) Links {
	links := make(Links, 0)

	// One chunk: <url>; rel="foo"
	for _, chunk := range strings.Split(raw, ",") {

		link := Link{URL: "", Rel: "", Params: make(map[string]string)}

		// Figure out what each piece of the chunk is
		for _, piece := range strings.Split(chunk, ";") {

			piece = strings.Trim(piece, " ")
			if piece == "" {
				continue
			}

			// URL
			if piece[0] == '<' && piece[len(piece)-1] == '>' {
				link.URL = strings.Trim(piece, "<>")
				continue
			}

			// Params
			key, val := parseParam(piece)
			if key == "" {
				continue
			}

			// Special case for rel
			if strings.ToLower(key) == "rel" {
				link.Rel = val
			}

			link.Params[key] = val

		}

		links = append(links, link)
	}

	return links
}

// ParseMultiple is like Parse, but accepts a slice of headers
// rather than just one header string
func ParseMultiple(headers []string) Links {
	links := make(Links, 0)
	for _, header := range headers {
		links = append(links, Parse(header)...)
	}
	return links
}

// parseParam takes a raw param in the form key="val" and
// returns the key and value as seperate strings
func parseParam(raw string) (key, val string) {

	parts := strings.SplitN(raw, "=", 2)
	if len(parts) != 2 {
		return "", ""
	}

	key = parts[0]
	val = strings.Trim(parts[1], "\"")

	return key, val

}