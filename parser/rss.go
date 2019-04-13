// Package parser parses
//
// TODO
//
//   - Implement custom unmarshaling to trim space
package parser

import (
	"encoding/xml"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	TagStartContentEncoded = `<content:encoded>`
	TagEndContentEncoded   = `</content:encoded>`
)

type Feed struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"Default link"`
	Image       string `xml:"image>url"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title           string `xml:"title"`
	Link            string `xml:"Default link"`
	PublicationDate string `xml:"pubDate"`
}

func LoadURL(url string) (Feed, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Feed{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Feed{}, errors.Errorf("server responded with status %d", resp.StatusCode)
	}

	feed, err := Load(resp.Body)
	if err != nil {
		return Feed{}, err
	}

	return feed, nil
}

func LoadFile(file string) (Feed, error) {
	f, err := os.Open(file)
	if err != nil {
		return Feed{}, err
	}
	defer f.Close()
	return Load(f)
}

func Load(r io.Reader) (Feed, error) {
	var feed Feed
	dec := xml.NewDecoder(r)
	dec.DefaultSpace = "Default"
	if err := dec.Decode(&feed); err != nil {
		return Feed{}, err
	}

	return feed, nil
}

func Quote(input, startTag, endTag string) string {
	return ProcessElementText(input, startTag, endTag, strconv.Quote)
}

type StringModifierFunc func(input string) string

func ProcessElementText(input, startTag, endTag string, f StringModifierFunc) string {
	var output strings.Builder

	currentIndex := 0
	for {
		// 1. Find index of first start tag
		startTagIndex := strings.Index(input[currentIndex:], startTag)
		if startTagIndex == -1 {
			// 9. Append the rest
			output.WriteString(input[currentIndex:])
			// 10. Return the result
			return output.String()
		}
		startTagIndex += currentIndex

		// 2. Advance the index by the length of the start tag
		startTagIndex += len(startTag)

		// 3. Append everything up to end of the start tag
		output.WriteString(input[currentIndex:startTagIndex])

		// 4. Find the index of the first end tag
		endTagIndex := strings.Index(input[startTagIndex:], endTag)
		if endTagIndex == -1 {
			// This is a malformed string with no matching end tag
			return ""
		}
		endTagIndex += startTagIndex

		// 5. Append the modified version of the interstitial text
		output.WriteString(f(input[startTagIndex:endTagIndex]))

		// 6. Append the end tag
		output.WriteString(endTag)

		// 7. Advance the index by the length of the end tag
		currentIndex = endTagIndex + len(endTag)

		// 8. Repeat from step 1 until start tag not present
	}
}
