package rss

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/haleyrc/rss/parser"
)

var (
	ErrTitleRequired         = errors.New("title is required")
	ErrDescriptionIsRequired = errors.New("description is required")
	ErrLinkRequired          = errors.New("link is required")
	ErrFeedRequired          = errors.New("feed is required")
)

func NewFeed(title, description, link, image string, items ...*Item) (*Feed, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrTitleRequired
	}
	description = strings.TrimSpace(description)
	if description == "" {
		return nil, ErrDescriptionIsRequired
	}
	link = strings.TrimSpace(link)
	if link == "" {
		return nil, ErrLinkRequired
	}
	image = strings.TrimSpace(image)
	return &Feed{
		Title:       title,
		Description: description,
		Link:        link,
		Image:       image,
		Items:       items,
	}, nil
}

func NewFromChannel(c parser.Channel) (*Feed, error) {
	var items []*Item
	for _, item := range c.Items {
		pubDate, err := time.Parse(time.RFC1123, item.PublicationDate)
		if err != nil {
			log.Printf("error parsing publication date: %s: %v: skipping\n", item.PublicationDate, err)
			continue
		}
		newItem, err := NewItem(-1, item.Title, item.Link, pubDate)
		if err != nil {
			log.Printf("invalid item: %v: skipping\n", err)
			continue
		}
		items = append(items, newItem)
	}
	feed, err := NewFeed(c.Title, c.Description, c.Link, c.Image, items...)
	if err != nil {
		return nil, err
	}
	return feed, nil
}

type Feed struct {
	ID          int64
	Title       string
	Description string
	Link        string
	Image       string
	Items       []*Item
}

func NewItem(feed int64, title, link string, pub time.Time) (*Item, error) {
	if feed == 0 {
		return nil, ErrFeedRequired
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrTitleRequired
	}
	link = strings.TrimSpace(link)
	if link == "" {
		return nil, ErrLinkRequired
	}
	if pub.IsZero() {
		pub = time.Now()
	}
	return &Item{
		FeedID:          feed,
		Title:           title,
		Link:            link,
		PublicationDate: pub,
	}, nil
}

type Item struct {
	ID              int64     `db:"id"`
	FeedID          int64     `db:"feed_id"`
	Title           string    `db:"title"`
	Link            string    `db:"link"`
	PublicationDate time.Time `db:"publication_date"`
	Read            bool      `db:"read"`
	Starred         bool      `db:"starred"`
	Ignored         bool      `db:"ignored"`
}
