package rss

import (
	"encoding/xml"
	"io"
)

type Feed struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Image       string `xml:"image>url"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title           string `xml:"title"`
	Description     string `xml:"description"`
	Link            string `xml:"link"`
	PublicationDate string `xml:"pubDate"`
}

func Load(r io.Reader) (Feed, error) {
	var feed Feed
	if err := xml.NewDecoder(r).Decode(&feed); err != nil {
		return Feed{}, err
	}
	return feed, nil
}
