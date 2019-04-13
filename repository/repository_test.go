package repository_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/haleyrc/rss"
	"github.com/haleyrc/rss/parser"
	"github.com/haleyrc/rss/repository"
)

func TestFeedItems(t *testing.T) {
	db := sqlx.MustConnect("postgres", "host=localhost user=postgres password=test port=5433 dbname=rss sslmode=disable")
	client := repository.New(db)

	feed, err := rss.NewFeed("test feed", "this is a test", "http://example.com", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := client.CreateFeed(feed); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	item1, err := rss.NewItem(feed.ID, "Initial title", "http://example.com/item/1", time.Now().Add(-24*time.Hour))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	item2, err := rss.NewItem(feed.ID, "New title", "http://example.com/item/1", time.Now().Add(24*time.Hour))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := client.CreateItem(item1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := client.CreateItem(item2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if item1.ID != item2.ID {
		t.Errorf("expected ids to be equal, %d != %d", item1.ID, item2.ID)
	}

	got, err := client.GetItem(item1.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Title != item2.Title {
		t.Errorf("expected title %q, got %q", item2.Title, got.Title)
	}
	if got.Link != item2.Link {
		t.Errorf("expected link %q, got %q", item2.Link, got.Link)
	}
	wantDate := item2.PublicationDate.Round(time.Millisecond)
	gotDate := got.PublicationDate.Round(time.Millisecond)
	if !gotDate.Equal(wantDate) {
		t.Errorf("expected pub date %s, got %s", wantDate, gotDate)
	}
}

func TestRepository(t *testing.T) {
	db := sqlx.MustConnect("postgres", "host=localhost user=postgres password=test port=5433 dbname=rss sslmode=disable")
	client := repository.New(db)

	for _, filename := range []string{"hackernews.xml", "checkly.xml"} {
		xmlFeed, err := parser.LoadFile(filepath.Join("..", "testdata", filename))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		feed, err := rss.NewFromChannel(xmlFeed.Channel)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if err := client.CreateFeed(feed, feed.Items...); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if feed.ID == 0 {
			t.Fatalf("expected id to be set, but got %d", feed.ID)
		}
	}
}
